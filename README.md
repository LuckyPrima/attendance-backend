# Attendance Backend

Attendance Backend dibangun menggunakan **Golang** (Go `v1.24.5`) dengan framework **Gin** dan ORM **GORM**.  
Database yang digunakan adalah **MySQL**, dengan konfigurasi koneksi dikelola melalui file `.env`.

---

## 🚀 Fitur Utama

- Manajemen **Employees**
- Manajemen **Departments**
- Pencatatan **Attendance (Clock In / Clock Out)**
- **Attendance Logs** dengan filter tanggal, departemen, dan pagination
- Migrasi database otomatis saat server dijalankan

---

## 📦 Persyaratan

- Go `>=1.24.5`
- MySQL `>=8`
- Git
- `go-migrate`

---

## ⚙️ Setup Project

### 1. Clone Repository

```bash
git clone https://github.com/LuckyPrima/attendance-backend.git
cd attendance-backend
```

### 2. Buat file .env

Isi .env sesuai konfigurasi database kamu:

```env
DB_USER=user
DB_PASS=yourpassword
DB_HOST=yourIPHost
DB_PORT=yourPortDB
DB_NAME=yourDBname
```

### 3. Install Dependency

```bash
go mod tidy
```

### 4. Jalankan Server

```bash
go run main.go
```

Server akan berjalan di

```
http://localhost:5000
```

---

## 🗄️ Struktur Proyek

```bash
attendance-backend/
│── config/        # Koneksi database (db.go)
│── controllers/   # Business logic (employee, department, attendance)
│── models/        # Struct model GORM
│── routes/        # Definisi routing API
│── main.go        # Entry point
│── .env           # Konfigurasi environment
```

---

## 📖 API Documentation

### 🔹 Employees

| Method | Endpoint         | Body                                                             | Deskripsi                          |
| ------ | ---------------- | ---------------------------------------------------------------- | ---------------------------------- |
| GET    | `/employees`     | -                                                                | Ambil semua employees (sorted A-Z) |
| POST   | `/employees`     | `{ "name": "Budi", "address": "Jakarta", "department_id": 1 }`   | Tambah employee                    |
| PUT    | `/employees/:id` | `{ "name": "Update", "address": "Bandung", "department_id": 2 }` | Update employee                    |
| DELETE | `/employees/:id` | -                                                                | Hapus employee                     |

### 🔹 Departments

| Method | Endpoint           | Body                                                                                      | Deskripsi              |
| ------ | ------------------ | ----------------------------------------------------------------------------------------- | ---------------------- |
| GET    | `/departments`     | -                                                                                         | Ambil semua departemen |
| POST   | `/departments`     | `{ "departement_name": "IT", "max_clock_in": "09:00", "max_clock_out": "17:00" }`         | Tambah departemen      |
| PUT    | `/departments/:id` | `{ "departement_name": "HR Updated", "max_clock_in": "08:30", "max_clock_out": "17:30" }` | Update departemen      |
| DELETE | `/departments/:id` | -                                                                                         | Hapus departemen       |

### 🔹 Attendance

| Method | Endpoint               | Body                                                              | Deskripsi                                      |
| ------ | ---------------------- | ----------------------------------------------------------------- | ---------------------------------------------- |
| POST   | `/attendance/clockin`  | `{ "employee_id": 1 }`                                            | Absen masuk                                    |
| PUT    | `/attendance/clockout` | `{ "employee_id": 1 }`                                            | Absen keluar                                   |
| GET    | `/attendance/logs`     | Query: `date=YYYY-MM-DD`, `department_id=1`, `page=1`, `limit=10` | Ambil log absensi (dengan filter + pagination) |

### Contoh Respons

```JSON
{
  "data": [
    {
      "history_id": 1,
      "employee_id": 3,
      "employee_name": "Budi",
      "department_id": 1,
      "department_name": "Human Resource",
      "date_attendance": "2025-09-18T05:01:19.325+07:00",
      "type": "Clock In",
      "description": "On Time",
      "accuracy": "On Time"
    }
  ],
  "page": 1,
  "limit": 10,
  "total": 20,
  "totalPages": 2
}

```

---

## 🛠️ Development Notes

1. Database migration otomatis saat menjalankan go run main.go.

2. Untuk production, gunakan migration manual agar lebih aman.

3. Pastikan .env di-setup dengan benar.
