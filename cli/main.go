package main

import (
	"fmt"
	"os"
)

const VERSION = "0.1.0"

func main() {
	if len(os.Args) < 2 {
		printHelp()
		return
	}

	switch os.Args[1] {
	case "-S", "--sync":
		if len(os.Args) < 3 {
			fmt.Println("Hata: paket adı gerekli")
			os.Exit(1)
		}
		installPackage(os.Args[2:])
	case "-Ss", "--search":
		if len(os.Args) < 3 {
			fmt.Println("Hata: arama terimi gerekli")
			os.Exit(1)
		}
		searchPackage(os.Args[2:])
	case "-R", "--remove":
		if len(os.Args) < 3 {
			fmt.Println("Hata: paket adı gerekli")
			os.Exit(1)
		}
		removePackage(os.Args[2:])
	case "-Syu", "--upgrade":
		upgradeAll()
	case "-Qi", "--info":
		if len(os.Args) < 3 {
			fmt.Println("Hata: paket adı gerekli")
			os.Exit(1)
		}
		packageInfo(os.Args[2])
	case "-Q", "--list":
		listInstalled()
	case "-V", "--version":
		fmt.Printf("aay v%s - Alpine AUR Helper\n", VERSION)
	default:
		printHelp()
	}
}

func printHelp() {
	fmt.Printf(`aay v%s - Alpine AUR Helper
Kullanım: aay <işlem> [paket]

İşlemler:
  -S  <paket>    Paket kur
  -Ss <paket>    Paket ara
  -R  <paket>    Paket kaldır
  -Syu           Tüm AAR paketlerini güncelle
  -Qi <paket>    Paket bilgisi göster
  -Q             Kurulu AAR paketlerini listele
  -V             Versiyon bilgisi

Örnekler:
  aay -Ss htop
  aay -S htop
  aay -R htop
`, VERSION)
}
