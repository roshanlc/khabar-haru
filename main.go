package main

import (
	"html/template"
	"log"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/roshanlc/news-aggregator/websites"
)

// Page related information
type Page struct {
	Ekantipur, BBCNepal *[]websites.News
	Prices              websites.PetroleumPrices
}

type pageWithLock struct {
	page Page
	rw   *sync.RWMutex
}

// Start a webserver
func webServer(content *pageWithLock, wg *sync.WaitGroup, port string) {

	tmpl := template.Must(template.ParseFiles("static/index.html"))

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {

		// lock for reading purposes
		content.rw.RLock()
		defer content.rw.RUnlock()

		tmpl.Execute(w, Page{
			Ekantipur: content.page.Ekantipur,
			BBCNepal:  content.page.BBCNepal,
			Prices:    content.page.Prices,
		})

	})

	log.Printf("http server started at port 8080\n")
	log.Fatal(http.ListenAndServe(":"+port, nil))

	wg.Done()
}

func main() {

	// Necesary for deploying on heroku
	port := os.Getenv("PORT")

	if port == "" {
		//For local run
		port = "8080"

	}

	// Content struct
	content := pageWithLock{page: Page{}, rw: &sync.RWMutex{}}

	var wg sync.WaitGroup

	wg.Add(2)

	go func() {
		var init bool = false

		// Run for the intial time
		if !init {
			log.Println("Running scrapping methods")
			scrapeWebsites(&content)
			init = true
		}
		// This function runs repeatedly to scrape sites
		ticker := time.NewTicker(1 * time.Minute)

		// loop over the ticks
		for range ticker.C {

			log.Println("Running scrapping methods")
			scrapeWebsites(&content)
		}

		//Add wg.Done() method
		wg.Done()
	}()

	go func() {
		//scrapeWebsites(&content)
		webServer(&content, &wg, port)
	}()

	// Keep the web server goroutine keep running
	wg.Wait()
}

// A routine function
func scrapeWebsites(content *pageWithLock) {
	// Lock for writing purpose
	content.rw.Lock()
	defer content.rw.Unlock()

	ek := websites.FetchEkantipur()
	bbc := websites.FetchBBCNepali()
	prices := websites.FetchPrices()

	content.page.BBCNepal = bbc
	content.page.Ekantipur = ek
	content.page.Prices = prices

}
