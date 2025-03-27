# Booking Lapangan Mini Soccer!⚽

## Apaan sih ini? 🤔

Yok kenalan! Ini adalah sistem booking lapangan bola super kece yang dibangun pake teknologi canggih. Gak cuma sekedar booking, tapi serba digital abis!

## Teknologi Kece Bro! 🚀

- **Bahasa Coding**: Golang (Mantul!)
- **Frontend**: Next.js
- **Arsitektur**: Microservice (Biar keren)
- **Messaging**: Apache Kafka
- **Packaging**: Docker
- **Deploy**: Jenkins
- **Cloud**: Google Cloud Platform
- **Nyimpen File**: Google Cloud Storage
- **Database**: PostgreSQL
- **Service Discovery**: Consul KV

## Fitur Keren Banget! ✨

### Backend Jagoan
- Booking lapangan? Tinggal klik!
- Chat antar service real-time
- Atur lapangan sejagat raya
- Login super aman
- Transaksi nyaris tidak mungkin error

### Microservice Kece
- Service terpisah buat:
  - Ngatur User
  - Booking Lapangan
  - Urusan Duit
  - Kirim Notifikasi
- Komunikasi pake Kafka
- Bisa ngehandle jutaan user!

## Belajar Apa Aja Sih? 🎓

1. Microservice itu apaan?
2. Golang dari 0 sampai jago
3. Main-main sama cloud
4. Bikin backend kece abis!

## Ini Buat Siapa Aja? 👥

- Pemula yang pengen jadi jagoan backend
- Backend developer pengen upgrade skill
- Fullstack yang pengen belajar hal baru

## Mulai Ngoding Yuk! 💻

### Persiapan 
- Golang 1.20+
- Docker
- PostgreSQL
- Akun GCP
- Kafka

### Install Kilat!

```bash
# Copot repository
https://github.com/bintangnugrahaa/soccer-booking-field

# Masuk folder
cd soccer-booking-field

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

## Wajib Setting! 🔧

Salin `.env.example` terus diubah:
- Login GCP
- Koneksi Database
- Setting Kafka
- Pengamanan

**Selamat Ngoding, Semoga Jadi Developer Kece! 💻🔥**
