package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-rod/rod"
)

func Healthcheck(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Healthcheck!"})
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
