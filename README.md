# Booking Lapangan Mini Soccer!⚽
## Goal? 
Sistem ini dirancang untuk memberikan kemudahan kepada para pecinta mini soccer dalam membooking lapangan mini soccer secara online dengan mudah. Dengan menggunakan sistem ini, semua proses mulai dari pemesanan hingga pembayaran dapat dilakukan melalui sistem ini.

## Teknologi 🚀

- **Programming Language**: Golang
- **Arsitektur**: Microservice
- **Framework**: Gin & Gorm
- **Database**: PostgreSQL
- **Message Broker**: Apache Kafka
- **Containerization**: Docker
- **Cloud**: Google Cloud Platform
- **Storage**: Google Cloud Storage

## Flow Aplikasi
-**1. Admin → Login**
-**2. Admin → Manage Field**
-**3. Admin → Manage Field Schedule**
-**4. Customer → Register**
-**5. Customer → Login**
-**6. Customer → Create Order**
-**7. Customer → Pay Order**
-**8. Done**


### Persiapan 
- Golang 1.20+
- Docker
- PostgreSQL
- Akun GCP
- Kafka

### Install Kilat!

```bash
# Copy repository
https://github.com/alghoziii/mini-soccer-booking

# Masuk folder
cd mini-soccer-booking

# Jalankan Docker
docker-compose up -d

# Migrate database
make migrate

# Nyalain service
make run
```

## Struktur Projek Kita 📂

```
booking-lapangan-bola/
│
├── service-user/           # Ngurusin User
├── service-booking/        # Booking Lapangan
├── service-pembayaran/     # Urus Duit
├── service-notifikasi/     # Kirim Pesan
├── api-gateway/            # Pintu Masuk
└── share/                  # Kode Umum
```

Documentation : https://wobbly-typhoon-983.notion.site/Konsep-Sistem-Booking-Lapangan-Mini-Soccer-Backend-Golang-Microservice-145e3a9651db8093a32fd571e0022670
https://wobbly-typhoon-983.notion.site/Konsep-Sistem-Booking-Lapangan-Mini-Soccer-Backend-Golang-Microservice-145e3a9651db8093a32fd571e0022670
