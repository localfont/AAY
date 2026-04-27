package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
)

type Package struct {
	Name        string `json:"name"`
	Version     string `json:"version"`
	Description string `json:"description"`
	Maintainer  string `json:"maintainer"`
	Votes       int    `json:"votes"`
	URL         string `json:"url"`
}

func getAPIURL() string {
	api := os.Getenv("AAR_API")
	if api == "" {
		api = "http://localhost:8000"
	}
	return api
}

func searchPackage(args []string) {
	query := args[0]
	apiURL := fmt.Sprintf("%s/api/search?q=%s", getAPIURL(), url.QueryEscape(query))

	resp, err := http.Get(apiURL)
	if err != nil {
		fmt.Printf("\033[31mHata:\033[0m API'ye bağlanılamadı: %v\n", err)
		return
	}
	defer resp.Body.Close()

	var packages []Package
	if err := json.NewDecoder(resp.Body).Decode(&packages); err != nil {
		fmt.Printf("\033[31mHata:\033[0m Yanıt ayrıştırılamadı\n")
		return
	}

	if len(packages) == 0 {
		fmt.Printf("'%s' için sonuç bulunamadı\n", query)
		return
	}

	fmt.Printf("\n\033[34m::\033[0m %d paket bulundu:\n\n", len(packages))
	for _, p := range packages {
		fmt.Printf("\033[32maar/\033[0m\033[1m%s\033[0m \033[33m%s\033[0m (oy: %d)\n",
			p.Name, p.Version, p.Votes)
		fmt.Printf("    %s\n\n", p.Description)
	}
}

func packageInfo(name string) {
	apiURL := fmt.Sprintf("%s/api/info/%s", getAPIURL(), name)

	resp, err := http.Get(apiURL)
	if err != nil {
		fmt.Printf("\033[31mHata:\033[0m API'ye bağlanılamadı\n")
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode == 404 {
		fmt.Printf("'%s' paketi bulunamadı\n", name)
		return
	}

	var p Package
	json.NewDecoder(resp.Body).Decode(&p)

	fmt.Printf(`
\033[1mPaket Adı    :\033[0m %s
\033[1mVersiyon     :\033[0m %s
\033[1mAçıklama     :\033[0m %s
\033[1mYapımcı      :\033[0m %s
\033[1mURL          :\033[0m %s
\033[1mOy           :\033[0m %d
`, p.Name, p.Version, p.Description, p.Maintainer, p.URL, p.Votes)
}
