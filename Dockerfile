# =============================================================================
# Build Stage
# Menggunakan Go image untuk compile source code menjadi binary
# =============================================================================
FROM golang:1.25-alpine AS builder

# Set working directory
WORKDIR /app

# Copy dependency files terlebih dahulu
# Dilakukan terpisah agar Docker cache layer ini
# sehingga tidak perlu download ulang jika kode berubah
COPY go.mod go.sum ./

# Download semua dependency
RUN go mod download

# Copy seluruh source code
COPY . .

# Build binary dengan optimasi untuk production
# CGO_ENABLED=0  — matikan CGO agar binary bisa jalan di alpine
# GOOS=linux     — target OS adalah Linux
# -o main        — nama output binary adalah "main"
RUN CGO_ENABLED=0 GOOS=linux go build -o main .

# =============================================================================
# Production Stage
# Menggunakan image minimal agar ukuran container sekecil mungkin
# =============================================================================
FROM alpine:latest

# Install dependencies yang dibutuhkan saat runtime
# ca-certificates — untuk koneksi HTTPS
# tzdata          — untuk timezone
RUN apk --no-cache add ca-certificates tzdata

# Set timezone ke Asia/Jakarta
ENV TZ=Asia/Jakarta

# Set working directory
WORKDIR /app

# Copy hanya binary hasil build dari stage sebelumnya
# Tidak perlu copy source code atau Go compiler ke production
COPY --from=builder /app/main .

# Copy file environment
COPY .env .

# Expose port aplikasi
EXPOSE 8080

# Jalankan aplikasi
CMD ["./main"]