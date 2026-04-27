#!/bin/sh
# aay kurulum scripti - Alpine Linux
# Kullanım: sh install.sh

set -e

REPO="alpineaar/aay"  # GitHub kullanıcı/repo adını değiştir
BIN="/usr/local/bin/aay"

# Mimariyi tespit et
ARCH=$(uname -m)
case "$ARCH" in
  x86_64)  TARGET="aay-linux-amd64" ;;
  aarch64) TARGET="aay-linux-arm64" ;;
  i386|i686) TARGET="aay-linux-386" ;;
  *)
    echo "Hata: Desteklenmeyen mimari: $ARCH"
    exit 1
    ;;
esac

echo "==> aay kuruluyor ($TARGET)..."

# Son sürümü bul
LATEST=$(wget -qO- "https://api.github.com/repos/${REPO}/releases/latest" \
  | grep '"tag_name"' \
  | sed 's/.*"tag_name": "\(.*\)".*/\1/')

if [ -z "$LATEST" ]; then
  echo "Hata: Son sürüm bulunamadı"
  exit 1
fi

echo "==> Sürüm: $LATEST"

# İndir
URL="https://github.com/${REPO}/releases/download/${LATEST}/${TARGET}"
wget -q --show-progress -O /tmp/aay "$URL" || {
  echo "Hata: İndirme başarısız: $URL"
  exit 1
}

chmod +x /tmp/aay
mv /tmp/aay "$BIN"

echo ""
echo "✓ aay $LATEST kuruldu!"
echo ""
echo "Kullanıma başlamak için AAR_API ortam değişkenini ayarlayın:"
echo "  export AAR_API=\"http://your-server:8000\""
echo ""
echo "Kullanım:"
echo "  aay -Ss <paket>   Ara"
echo "  aay -S  <paket>   Kur"
echo "  aay -Syu          Güncelle"
