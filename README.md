# yelp
```
├── cmd/                # Dastur boshlanish nuqtasi.
│   └── main.go         # Dastur ishga tushiriladigan fayl.
├── internal/           # Ichki logika va biznes kodlar.
│   ├── entity/         # Biznes obyektlari va asosiy domen modellari.
│   ├── usecase/        # Asosiy biznes logika va interfeyslar.
│   ├── adapter/        # Tashqi interfeyslar uchun adapterlar.
│   │   ├── http/       # HTTP server, marshrutlash va controllerlar.
│   │   ├── db/         # Ma'lumotlar bazasi adapterlari.
│   │   ├── redis/      # Redis caching adapteri.
│   │   └── auth/       # JWT va Casbin adapterlari (rol boshqaruvi).
│   ├── repository/     # DB interfeyslarini amalga oshirish.
│   └── service/        # Repozitoriyga bog‘liq bo‘lmagan biznes logikalar.
├── config/             # Konfiguratsiya fayllari.
│   └── config.go       # Konfiguratsiyani o‘qish va yuklash.
├── migrations/         # Ma'lumotlar bazasi migratsiyalari.
├── api/                # API hujjatlari (Swagger).
│   └── swagger.yaml    # Swagger hujjati.
├── pkg/                # Umumiy kutubxonalar yoki yordamchi kodlar.
│   └── logger/         # Loglar bilan ishlash.
└── tools/              # Qo‘shimcha asboblar (masalan, migratsiya vositalari).
```