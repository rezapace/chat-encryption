# Aplikasi Chat Menggunakan Algoritma RSA

## Deskripsi
Proyek ini adalah aplikasi chat sederhana yang menggunakan algoritma RSA untuk enkripsi dan dekripsi pesan. Aplikasi ini terdiri dari server dan dua klien yang dapat saling berkomunikasi dengan aman melalui enkripsi RSA.

## Kegunaan
Aplikasi ini berguna untuk:
- Mengamankan komunikasi antara dua pihak dengan menggunakan enkripsi RSA.
- Memahami implementasi dasar dari algoritma RSA dalam konteks aplikasi chat.
- Mempelajari cara kerja socket dan threading dalam Python untuk membangun aplikasi jaringan.

## Fungsi
Aplikasi ini terdiri dari beberapa modul dengan fungsi-fungsi sebagai berikut:

### RSA.py
- `gcd(a, b)`: Menghitung faktor persekutuan terbesar dari dua bilangan.
- `modinv(phi, m)`: Menghitung invers modulo.
- `coprimes(phi)`: Menghasilkan bilangan coprime dari phi.
- `key_generator()`: Menghasilkan pasangan kunci publik dan privat.
- `encrypt_block(m, e, n)`: Mengenkripsi blok pesan.
- `decrypt_block(c, d, n)`: Mendekripsi blok pesan.
- `encrypt_string(s, public_key)`: Mengenkripsi string pesan.
- `decrypt_string(s, private_key)`: Mendekripsi string pesan.

### Server.py
- `accept_incoming_connections()`: Menerima koneksi dari klien.
- `handle_client1(client_sock, client_addresses)`: Menangani komunikasi dengan klien pertama.
- `handle_client2(client_sock, client_addresses)`: Menangani komunikasi dengan klien kedua.

### client_1.py dan client_2.py
- `receive()`: Menerima pesan dari server dan mendekripsinya.
- `send(event=None)`: Mengirim pesan ke server setelah mengenkripsinya.
- `on_closing(event=None)`: Menutup koneksi dan keluar dari aplikasi.

### primes.py
- `choose_distinct_primes()`: Memilih dua bilangan prima yang berbeda dari daftar bilangan prima.

## Bagaimana Menjalankan
1. Jalankan `Server.py` untuk memulai server.
2. Jalankan `client_1.py` dan `client_2.py` untuk memulai dua klien.
3. Masukkan alamat host dan port yang sesuai pada klien.
4. Masukkan nama pengguna pada klien.
5. Mulai mengirim dan menerima pesan yang terenkripsi.

## Kesimpulan
Aplikasi ini menunjukkan bagaimana algoritma RSA dapat digunakan untuk mengamankan komunikasi dalam aplikasi chat. Dengan menggunakan enkripsi dan dekripsi, pesan yang dikirim antara dua klien dapat dijaga kerahasiaannya. Proyek ini juga memberikan pemahaman tentang penggunaan socket dan threading dalam Python untuk membangun aplikasi jaringan yang aman.