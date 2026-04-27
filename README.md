# aay - Alpine AUR Helper

Alpine Linux için AUR + yay benzeri topluluk paket yönetim sistemi.

## Bileşenler

```
aay/
├── cli/              # Go ile yazılmış CLI aracı
│   ├── main.go       # Giriş noktası, komut yönlendirici
│   ├── search.go     # Arama & bilgi
│   ├── install.go    # Kurulum, kaldırma, güncelleme
│   └── go.mod
├── api/              # Python FastAPI backend
│   ├── main.py       # REST API
│   ├── requirements.txt
│   └── Dockerfile
├── templates/
│   └── AARBUILD      # Paket tanım dosyası şablonu
├── scripts/
│   └── install.sh    # Tek satır kurulum scripti
├── .github/
│   └── workflows/
│       ├── build-cli.yml   # Go binary derleme
│       ├── build-api.yml   # Docker image build
│       └── release.yml     # Tam release pipeline
└── docker-compose.yml
```

---

## Kurulum

### CLI (Alpine Linux)
```sh
sh -c "$(wget -qO- https://raw.githubusercontent.com/KULLANICI/aay/main/scripts/install.sh)"
```

### API Sunucusu
```sh
git clone https://github.com/KULLANICI/aay
cd aay
docker-compose up -d
```

---

## Kullanım

```sh
# Ortam değişkeni (bir kez ayarla)
export AAR_API="http://localhost:8000"

# Paket ara
aay -Ss htop

# Paket kur
aay -S htop

# Paket kaldır
aay -R htop

# Tümünü güncelle
aay -Syu

# Kurulu paketleri listele
aay -Q

# Paket bilgisi
aay -Qi htop
```

---

## Paket Yayınlama (AARBUILD)

1. GitHub'da yeni repo oluştur: `github.com/KULLANICI/paketadim`
2. Repoya `AARBUILD` dosyası ekle (şablon: `templates/AARBUILD`)
3. AAR'a gönder:

```sh
curl -X POST http://localhost:8000/api/submit \
  -H "Content-Type: application/json" \
  -d '{
    "name": "paketadim",
    "version": "1.0.0",
    "description": "Paket açıklaması",
    "maintainer": "kullanici",
    "repo_url": "https://github.com/KULLANICI/paketadim.git",
    "depends": ["curl"]
  }'
```

---

## GitHub Actions

### Otomatik Build
- `main` branch'e push → tüm platformlar için binary derlenir
- `v*` tag'i → GitHub Release oluşturulur, binary'ler eklenir
- API için Docker image → `ghcr.io` üzerine push edilir

### Release
```sh
git tag v1.0.0
git push origin v1.0.0
# → GitHub Actions otomatik olarak release oluşturur
```

---

## API Endpointleri

| Metod | Endpoint | Açıklama |
|-------|----------|----------|
| GET | `/api/search?q=` | Paket ara |
| GET | `/api/info/:name` | Paket bilgisi |
| POST | `/api/submit` | Yeni paket gönder |
| PUT | `/api/update/:name` | Paket güncelle |
| DELETE | `/api/delete/:name` | Paket sil |
| POST | `/api/vote/:name` | Oy ver |
| GET | `/api/comments/:name` | Yorumları getir |
| POST | `/api/comments/:name` | Yorum ekle |
| GET | `/api/stats` | İstatistikler |

---

## Lisans

MIT
