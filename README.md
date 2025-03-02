# Healcation Backend

Healcation Backend adalah aplikasi backend berbasis Golang yang menggunakan **Gorm** sebagai ORM dan **SQLite** sebagai database default. Backend ini menangani fitur user management, history perjalanan, dan akomodasi.

## 🚀 Features
- **User Management** (registrasi, autentikasi, hash password)
- **History & Travel Planning** (CRUD perjalanan, tempat, akomodasi)
- **Gorm ORM** (untuk database)
- **Postman API Testing** (terintegrasi dengan Postman)

---

## 📋 Requirements
- **Go 1.23**
- **SQLite Database** 
- **Postman** (untuk testing API)

---

## ⚙️ Installation
1. **Install Dependencies**
   ```sh
   go mod tidy
   ```

---

## 🛠 Configuration
Sebelum menjalankan aplikasi, pastikan telah mengatur environment variables.  
Buat file `.env` dan tambahkan:

```sh
PORT=3000
DB_PATH="healcation.db" 
```

---

## 📦 Running the Application
1. **Jalankan Server**
   ```sh
   go run main.go
   ```
2. **Cek API dengan Postman**
   - Gunakan Postman untuk melakukan request ke `http://localhost:3000`
   - Endpoint tersedia di `routes.go`

---

## 🔄 Database Migration & Seeding
Jika database belum ada, jalankan migrasi dan seeder secara otomatis:

```sh
go run main.go
```
Seeder akan menambahkan data awal (admin user & history perjalanan).

---

## 📜 License
MIT License © 2025 ForgeSherpa

---
