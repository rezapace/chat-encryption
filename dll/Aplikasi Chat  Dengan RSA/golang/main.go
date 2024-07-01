package main

import (
	"bufio"
	"crypto/rand"
	"crypto/rsa"
	"encoding/base64"
	"fmt"
	"math/big"
	"net"
	"os"
	"strconv"
	"strings"
	"time"
)

// Receive handles receiving of messages.
func receive(client net.Conn, bufferSize int, privateKey1 *rsa.PrivateKey) {
	name := os.Args[3]
	msgList := " Welcome! " + name
	msgList += " You are online!"
	for {
		msg, err := bufio.NewReader(client).ReadString('\n')
		if err != nil {
			break
		}
		decoded, _ := base64.StdEncoding.DecodeString(msg)
		decrypted, _ := rsa.DecryptPKCS1v15(rand.Reader, privateKey1, decoded)
		msgList += string(decrypted)
	}
}

// Send handles sending of messages.
func send(client net.Conn, myMsg string, publicKey2 *rsa.PublicKey) {
	name := os.Args[3]
	msg := name + ": " + myMsg
	encrypted, _ := rsa.EncryptPKCS1v15(rand.Reader, publicKey2, []byte(msg))
	encoded := base64.StdEncoding.EncodeToString(encrypted)
	client.Write([]byte(encoded + "\n"))
}

// OnClosing is to be called when the window is closed.
func onClosing() {
	msgList := ""
	msgList += "going offline..."
	time.Sleep(2 * time.Second)
	client.Close()
	os.Exit(0)
}

func main() {
	host := os.Args[1]
	port, _ := strconv.Atoi(os.Args[2])
	bufferSize := 1024
	address := host + ":" + strconv.Itoa(port)

	client, err := net.Dial("tcp", address)
	if err != nil {
		fmt.Println(err)
		return
	}

	publicKey1, privateKey1, _ := rsa.GenerateKey(rand.Reader, 2048)
	msg := strconv.Itoa(publicKey1.E) + "*" + publicKey1.N.String()
	client.Write([]byte(msg + "\n"))
	m, _ := bufio.NewReader(client).ReadString('\n')
	publicKey2Str := strings.Split(m, "*")
	e, _ := strconv.Atoi(publicKey2Str[0])
	n := new(big.Int)
	n.SetString(publicKey2Str[1], 10)
	publicKey2 := &rsa.PublicKey{N: n, E: e}

	go receive(client, bufferSize, privateKey1)

	myMsg := ""
	for {
		fmt.Scanln(&myMsg)
		send(client, myMsg, publicKey2)
	}

	onClosing()
}
