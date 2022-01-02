package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"sync"

	"github.com/gocolly/colly"
)

type News struct {
	Title string
	Link  string
}

type Page struct {
	PageTitle string
	AllNews   *[]News
}

func webServer(title string, links *[]News, wg *sync.WaitGroup) {

	tmpl := template.Must(template.ParseFiles("html/index.html"))

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {

		tmpl.Execute(w, Page{
			PageTitle: title,
			AllNews:   links,
		})

	})

	log.Printf("http server started at port 8080\n")
	log.Fatal(http.ListenAndServe(":8080", nil))

	wg.Done()
}

func scrapeEkantipur() []News {

	temp := make([]News, 0)
	collector := colly.NewCollector(
	// Visit only allowed domains
	//colly.AllowedDomains("nepalnews.com", "https://www.nepalnews.com"),
	)

	collector.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting ", r.URL)

	})

	collector.OnHTML("a[data-type=\"title\"]", func(h *colly.HTMLElement) {

		title := h.Text
		url := h.Attr("href")
		temp = append(temp, News{
			Title: title,
			Link:  url,
		})

	})

	collector.Visit("https://ekantipur.com")

	collector.Wait()

	fmt.Println("Data scraping completion!!")
	return temp
}

func main() {

	allLinks := scrapeEkantipur()

	var wg sync.WaitGroup

	wg.Add(1)

	go webServer("News Aggregatorr", &allLinks, &wg)

	// Keep the web server goroutine keep running
	wg.Wait()

}
