package handlers

import (
	"api/internal/database"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/go-rod/rod"
	"github.com/otiai10/gosseract/v2"
)

const ErrorImage = "must use jpeg, png, or tiff file types"

func Healthcheck(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Healthcheck!"})
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

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("test login")
	rows, err := database.DB().Query("SELECT * FROM users")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer rows.Close()

	type Users struct {
		name         string
		first        string
		last         string
		time_created string
	}
	fmt.Printf("%+v\n", rows)

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "login handler fired!"})
}

func scrape(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	// https://www.spiegel.de/international/world/escalating-violence-radical-settlers-on-the-west-bank-see-an-opportunity-a-9499f824-9b39-4739-b6db-36772bc2bb99
	// Decode the request body into the User struct
	url := r.URL.Query().Get("url")
	fmt.Println("url =>", url)

	// Launch a new browser with default options, and connect to it.

	browser := rod.New().MustConnect()

	// close it after main process ends.
	defer browser.MustClose()

	// Create a new page
	page := browser.MustPage(url).MustWaitStable()
	//DER SPIEGEL EXAMPLE
	//heading
	// fmt.Println(page.MustElement("#Inhalt > article > header > div > div").MustEval(`() => this.innerText`).String())
	heading := page.MustElement("#Inhalt > article > header > div > div").MustEval(`() => this.innerText`).String()
	// fmt.Printf(heading)

	// main := page.MustElement("main").MustEval(`() => {
	// 	return this.innerText
	// 	}`)
	// fmt.Printf("%+v\n", main)

	main := page.MustElements("header")
	fmt.Printf("%+v\n", main)

	// main := page.MustSearch("main").MustEval(`() => {
	// 	return this
	// 	}`)
	// fmt.Printf("____  TEST  ____ \n%+v\n %T", main, main)
	// 	var res
	// 	err := json.Unmarshal([]byte(str), &res)
	// 	fmt.Println(err)
	// 	fmt.Println(res)
	// 	for _, m := range main {

	//     // m is a map[string]interface.
	//     // loop over keys and values in the map.
	//     for k, v := range m {
	//         fmt.Println(k, "value is", v)
	//     }
	// }

	// header := page.MustElement("heading")
	// fmt.Printf("%+v\n", header)

	text := page.MustElements("p")
	fmt.Printf("%+v\n %T", text, text)
	//body
	// fmt.Println(page.MustElement("#Inhalt > article > div.relative > section.relative > div > div").MustEval(`() => this.innerText`).String())
	body := page.MustElement("#Inhalt > article > div.relative > section.relative > div > div").MustEval(`() => this.innerText`).String()
	// fmt.Printf(body)

	json.NewEncoder(w).Encode(map[string]string{"heading": heading, "body": body})
}
