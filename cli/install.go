package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

type BuildInfo struct {
	Name         string   `json:"name"`
	Version      string   `json:"version"`
	Description  string   `json:"description"`
	Depends      []string `json:"depends"`
	MakeDepends  []string `json:"makedepends"`
	RepoURL      string   `json:"repo_url"`
}

var dbPath = "/var/lib/aay/installed.json"

func installPackage(args []string) {
	for _, pkgname := range args {
		doInstall(pkgname)
	}
}

func doInstall(pkgname string) {
	fmt.Printf("\033[34m::\033[0m \033[1m%s\033[0m çözümleniyor...\n", pkgname)

	// API'den paket bilgisi al
	apiURL := fmt.Sprintf("%s/api/info/%s", getAPIURL(), pkgname)
	resp, err := http.Get(apiURL)
	if err != nil {
		fmt.Printf("\033[31mHata:\033[0m API bağlantısı başarısız\n")
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode == 404 {
		fmt.Printf("\033[31mHata:\033[0m '%s' paketi AAR'da bulunamadı\n", pkgname)
		return
	}

	var info BuildInfo
	json.NewDecoder(resp.Body).Decode(&info)

	// Bilgi göster
	fmt.Printf("\n  Paket  : %s %s\n", info.Name, info.Version)
	fmt.Printf("  Açıkl. : %s\n", info.Description)
	if len(info.Depends) > 0 {
		fmt.Printf("  Bağ.   : %s\n", strings.Join(info.Depends, ", "))
	}

	// Onay al
	fmt.Printf("\n\033[33m==> Devam edilsin mi? [E/h]\033[0m ")
	reader := bufio.NewReader(os.Stdin)
	confirm, _ := reader.ReadString('\n')
	confirm = strings.TrimSpace(strings.ToLower(confirm))
	if confirm == "h" || confirm == "hayır" {
		fmt.Println("İptal edildi.")
		return
	}

	// Build dizini
	buildDir := filepath.Join(os.TempDir(), "aay-build", pkgname)
	os.MkdirAll(buildDir, 0755)

	// Repoyu klonla
	fmt.Printf("\033[34m::\033[0m Kaynak kod indiriliyor...\n")
	if err := runCmd("git", "clone", info.RepoURL, buildDir); err != nil {
		fmt.Printf("\033[31mHata:\033[0m Git clone başarısız\n")
		return
	}

	// APK bağımlılıklarını kur
	if len(info.Depends) > 0 {
		fmt.Printf("\033[34m::\033[0m Bağımlılıklar kuruluyor...\n")
		for _, dep := range append(info.Depends, info.MakeDepends...) {
			runCmd("apk", "add", "--no-cache", dep)
		}
	}

	// AARBUILD çalıştır
	fmt.Printf("\033[34m::\033[0m Derleniyor...\n")
	buildScript := filepath.Join(buildDir, "AARBUILD")
	if _, err := os.Stat(buildScript); err == nil {
		if err := runCmdDir(buildDir, "sh", "AARBUILD"); err != nil {
			fmt.Printf("\033[31mHata:\033[0m Derleme başarısız\n")
			return
		}
	}

	// Kurulumu kaydet
	saveInstalled(info.Name, info.Version)

	fmt.Printf("\033[32m✓\033[0m \033[1m%s %s\033[0m başarıyla kuruldu!\n",
		info.Name, info.Version)

	// Temizlik
	os.RemoveAll(buildDir)
}

func removePackage(args []string) {
	for _, pkgname := range args {
		fmt.Printf("\033[34m::\033[0m \033[1m%s\033[0m kaldırılıyor...\n", pkgname)

		// apk del dene
		if err := runCmd("apk", "del", pkgname); err != nil {
			// Manuel kaldır
			runCmd("rm", "-f", "/usr/local/bin/"+pkgname)
		}

		removeInstalled(pkgname)
		fmt.Printf("\033[32m✓\033[0m \033[1m%s\033[0m kaldırıldı\n", pkgname)
	}
}

func upgradeAll() {
	installed := loadInstalled()
	if len(installed) == 0 {
		fmt.Println("Güncellenecek AAR paketi yok.")
		return
	}

	fmt.Printf("\033[34m::\033[0m %d AAR paketi kontrol ediliyor...\n", len(installed))
	for name := range installed {
		fmt.Printf("  güncelleniyor: %s\n", name)
		doInstall(name)
	}
}

func listInstalled() {
	installed := loadInstalled()
	if len(installed) == 0 {
		fmt.Println("Kurulu AAR paketi yok.")
		return
	}

	fmt.Printf("\n\033[1mKurulu AAR Paketleri:\033[0m\n")
	for name, ver := range installed {
		fmt.Printf("  \033[32m%s\033[0m %s\n", name, ver)
	}
}

// DB yardımcıları
func loadInstalled() map[string]string {
	os.MkdirAll(filepath.Dir(dbPath), 0755)
	data, err := os.ReadFile(dbPath)
	if err != nil {
		return map[string]string{}
	}
	var result map[string]string
	json.Unmarshal(data, &result)
	return result
}

func saveInstalled(name, version string) {
	db := loadInstalled()
	db[name] = version
	data, _ := json.MarshalIndent(db, "", "  ")
	os.WriteFile(dbPath, data, 0644)
}

func removeInstalled(name string) {
	db := loadInstalled()
	delete(db, name)
	data, _ := json.MarshalIndent(db, "", "  ")
	os.WriteFile(dbPath, data, 0644)
}

// Komut çalıştırıcılar
func runCmd(name string, args ...string) error {
	cmd := exec.Command(name, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func runCmdDir(dir, name string, args ...string) error {
	cmd := exec.Command(name, args...)
	cmd.Dir = dir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
