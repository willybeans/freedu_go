package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/go-rod/rod"
	_ "github.com/lib/pq"
	"github.com/otiai10/gosseract/v2"
)

const ErrorImage = "must use jpeg, png, or tiff file types"

func Healthcheck(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Healthcheck"})
}

func ImageHandler(w http.ResponseWriter, r *http.Request) {
	file, hdr, err := r.FormFile("image")
	if err != nil {
		// If there is an error that means form is empty. Return nil for err in order
		// to validate result as required.
		fmt.Println("error at reading form", err)
		return
	}
	defer file.Close()

	// Create physical file
	tempfile, err := os.CreateTemp("", "ocrserver"+"-")
	if err != nil {
		fmt.Println("error at temp")
		return
	}
	defer func() {
		tempfile.Close()
		os.Remove(tempfile.Name())
	}()

	switch hdr.Header.Get("Content-Type") {
	case "application/octet-stream":
		// add conversion later for fun?
		json.NewEncoder(w).Encode(map[string]string{"message": ErrorImage})
		break
	case "image/tiff":
		fallthrough
	case "image/png":
		fallthrough
	case "image/jpeg":
		// Make uploaded physical
		if _, err = io.Copy(tempfile, file); err != nil {
			fmt.Println("error at copy")
		}

		client := gosseract.NewClient()
		defer client.Close()

		client.SetImage(tempfile.Name())
		// client.SetLanguage("deu") // need to receive lang from api keys
		client.Languages = []string{"eng"}
		if langs := r.FormValue("languages"); langs != "" {
			client.Languages = strings.Split(langs, ",")
		}

		pageText, err := client.Text()
		if err != nil {
			fmt.Println("error at text out: ", err)
			return
		}
		json.NewEncoder(w).Encode(map[string]string{"message": pageText})

	default:
		json.NewEncoder(w).Encode(map[string]string{"message": ErrorImage})
	}
}

func ScrapeHandler(w http.ResponseWriter, r *http.Request) {
	// https://www.spiegel.de/international/world/escalating-violence-radical-settlers-on-the-west-bank-see-an-opportunity-a-9499f824-9b39-4739-b6db-36772bc2bb99
	// Decode the request body into the User struct
	url := r.URL.Query().Get("url")
	fmt.Println("scraping url =>", url)
	// Launch a new browser with default options, and connect to it.
	browser := rod.New().MustConnect()
	defer browser.MustClose()
	// Create a new page
	page := browser.MustPage(url).MustWaitStable()
	body := page.MustElement("main").MustEval(`() => this.innerText`).String()
	heading := page.MustElement("h1").MustEval(`() => this.innerText`).String()

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"heading": heading, "body": body})
}
