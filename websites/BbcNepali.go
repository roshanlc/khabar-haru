package websites

import (
	"log"

	"github.com/gocolly/colly"
)

func FetchBBCNepali() *[]News {

	// The url of The Himalayan Times news
	const url = "https://www.bbc.com/nepali"

	const prefix = "https://www.bbc.com"

	collector := colly.NewCollector()

	temp := make([]News, 0)

	// Setting Custom User Agent
	collector.UserAgent = "Mozilla/5.0 (compatible; Googlebot/2.1; +http://www.google.com/bot.html)"

	collector.OnRequest(func(r *colly.Request) {
		log.Println("Visiting ", r.URL)

	})

	collector.OnError(func(r *colly.Response, err error) {
		log.Println("Some error while scraping ", url, err.Error())
	})

	// On finding a tag, run this function

	collector.OnHTML("a.bbc-1fxtbkn", func(h *colly.HTMLElement) {

		link := prefix + h.Attr("href")
		title := h.ChildText("span")

		temp = append(temp, News{Title: title, Link: link})
	})

	collector.Visit(url)

	collector.Wait()

	log.Println(url, ": Data scraping completion!!")

	// Returns news from BBC
	return &temp
}
