# Todo API

REST API untuk manajemen Todo yang dibangun dengan Go, Gin, GORM, dan PostgreSQL.

## 🛠 Tech Stack

- **Language**: Go
- **Framework**: Gin
- **ORM**: GORM
- **Database**: PostgreSQL
- **Authentication**: JWT
- **Documentation**: Swagger
- **Containerization**: Docker

## 🚀 Fitur

- ✅ Authentication (Register & Login)
- ✅ JWT Authorization
- ✅ CRUD Todo per user
- ✅ Input Validation
- ✅ Error Handling yang konsisten
- ✅ Rate Limiting
- ✅ Logging dengan Zap
- ✅ Unit Testing (85.7% coverage)
- ✅ Docker & Docker Compose
- ✅ API Documentation (Swagger)

## 📁 Struktur Project
todo-api/
├── main.go
├── config/
│   └── config.go
├── database/
│   └── database.go
├── internal/
│   ├── domain/
│   ├── handler/
│   ├── middleware/
│   ├── repository/
│   └── service/
├── docs/
├── Dockerfile
└── docker-compose.yml

## ⚙️ Instalasi

**Tanpa Docker:**
```bash
# Clone repository
git clone https://github.com/ewinkcess/todo-api.git
cd todo-api

# Copy environment file
cp .env.example .env
# Edit .env sesuai konfigurasi Anda

# Install dependencies
go mod tidy

# Jalankan aplikasi
go run main.go
```

**Dengan Docker:**
```bash
git clone https://github.com/ewinkcess/todo-api.git
cd todo-api
docker-compose up --build
```

## 🔑 Environment Variables

Buat file `.env` dengan isi berikut:

```env
DB_HOST=localhost
DB_USER=postgres
DB_PASSWORD=yourpassword
DB_NAME=todo_db
DB_PORT=5432
SERVER_PORT=8080
JWT_SECRET=your_secret_key
```

## 📖 API Documentation

Setelah server berjalan, akses Swagger di:
http://localhost:8080/swagger/index.html
## 🔐 API Endpoints

### Auth
| Method | Endpoint | Deskripsi |
|--------|----------|-----------|
| POST | /auth/register | Register user baru |
| POST | /auth/login | Login & dapat token |

### Todo (Butuh Token)
| Method | Endpoint | Deskripsi |
|--------|----------|-----------|
| GET | /todos | Ambil semua todo |
| GET | /todos/:id | Ambil satu todo |
| POST | /todos | Buat todo baru |
| PUT | /todos/:id | Update todo |
| DELETE | /todos/:id | Hapus todo |

## 🧪 Testing

```bash
# Jalankan semua test
go test ./...

# Dengan coverage
go test ./... -cover
```

## 📝 Contoh Request

**Register:**
```json
POST /auth/register
{
    "name": "John Doe",
    "email": "john@gmail.com",
    "password": "Password123!"
}
```

**Login:**
```json
POST /auth/login
{
    "email": "john@gmail.com",
    "password": "Password123!"
}
```

**Buat Todo:**
```json
POST /todos
Authorization: Bearer <token>

{
    "title": "Belajar Golang",
    "description": "Belajar membuat REST API"
}
```

