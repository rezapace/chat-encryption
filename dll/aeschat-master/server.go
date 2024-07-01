package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"sync"
	"time"
)

// Config variables:
// Default listening interface (use 0.0.0.0 for all)
var CONFIG_HTTP_INTERFACE = "127.0.0.1"

// Default listening port
var CONFIG_HTTP_PORT = 8888

// Number of chat lines to keep in memory for each channel
var CONFIG_BUFFER_LINES = 32

var chatLines = make(map[string]map[int]string)
var chatLinesMutex = &sync.Mutex{}

func processJsonRequest(jsonString string) string {
	jsonMap := make(map[string]string)

	if err := json.Unmarshal([]byte(jsonString), &jsonMap); err != nil {
		return ""
	}

	function := jsonMap["function"]
	channel := jsonMap["chan"]
	if function == "" || channel == "" {
		return ""
	}

	if !validHash(channel) {
		return "{}"
	}

	switch function {
	case "get":
		t, err := strconv.Atoi(jsonMap["t"])
		if err != nil || t < 0 {
			return ""
		}
		chatLinesMutex.Lock()
		defer chatLinesMutex.Unlock()
		if _, ok := chatLines[channel]; !ok {
			return ""
		}
		log := make(map[string]string)
		lastLine := chatLines[channel][0]
		log["t"] = lastLine
		log["text"] = ""
		for i := t + 1; i <= strconv.Atoi(lastLine)+1; i++ {
			if _, ok := chatLines[channel][strconv.Itoa(i)]; ok {
				log["text"] += chatLines[channel]["time"+strconv.Itoa(i)] + "," + chatLines[channel][strconv.Itoa(i)] + "\n"
			}
		}
		jsonData, _ := json.Marshal(log)
		return string(jsonData)
	case "post":
		line := jsonMap["line"]
		if line == "" {
			return ""
		}
		if len(line) > 1024 {
			return ""
		}
		chatLinesMutex.Lock()
		defer chatLinesMutex.Unlock()
		lastLine := 1
		if _, ok := chatLines[channel]; ok {
			lastLineInt, _ := strconv.Atoi(chatLines[channel][0])
			lastLine = lastLineInt + 1
		} else {
			chatLines[channel] = make(map[int]string)
		}
		chatLines[channel][0] = strconv.Itoa(lastLine)
		chatLines[channel][lastLine] = line
		chatLines[channel]["time"+strconv.Itoa(lastLine)] = strconv.FormatInt(time.Now().Unix(), 10)

		log := make(map[string]string)
		log["t"] = strconv.Itoa(lastLine - 1)

		// Purge old chat lines from memory
		for key := range chatLines[channel] {
			if keyInt, err := strconv.Atoi(key); err == nil {
				if keyInt != 0 && keyInt <= (lastLine-CONFIG_BUFFER_LINES) {
					// Purge entry
					delete(chatLines[channel], keyInt)
					delete(chatLines[channel], "time"+strconv.Itoa(keyInt))
				}
			}
		}

		jsonData, _ := json.Marshal(log)
		return string(jsonData)

	}
	return ""
}

func validHash(hash string) bool {
	if match, _ := regexp.MatchString("^[a-f0-9]{64}$", hash); match {
		return true
	}
	return false
}

func main() {
	if len(os.Args) > 1 {
		port, err := strconv.Atoi(os.Args[1])
		if err == nil {
			CONFIG_HTTP_PORT = port
		}
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" {
			body, err := ioutil.ReadAll(r.Body)
			if err != nil {
				http.Error(w, "Error reading request body", http.StatusInternalServerError)
				return
			}

			ret := processJsonRequest(string(body))

			if ret == "" {
				ret = "{}"
			}

			w.Header().Set("Content-Type", "application/json")
			w.Write([]byte(ret))
			return
		}

		// Serve static files
		uri := r.URL.Path
		filename := filepath.Join("web", filepath.Clean(uri))

		if uri == "/" {
			filename = filepath.Join(filename, "index.html")
		}

		fileInfo, err := os.Stat(filename)
		if err != nil || fileInfo.IsDir() {
			http.Error(w, "404 Not Found", http.StatusNotFound)
			return
		}

		http.ServeFile(w, r, filename)
	})

	addr := fmt.Sprintf("%s:%d", CONFIG_HTTP_INTERFACE, CONFIG_HTTP_PORT)
	fmt.Printf("aeschat server running at\n  => http://%s/\nCTRL + C to shutdown\n", addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}
