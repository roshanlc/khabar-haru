package websites

import (
	"log"

	"github.com/gocolly/colly"
)

func FetchEkantipur() *[]News {

	// The url of The Himalayan Times news
	const url = "https://ekantipur.com"

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

	collector.OnHTML("a[data-type=\"title\"]", func(h *colly.HTMLElement) {

		title := h.Text
		url := h.Attr("href")
		temp = append(temp, News{
			Title: title,
			Link:  url,
		})

	})
	collector.Visit(url)

	collector.Wait()

	log.Println(url, ": Data scraping completion with ", len(temp), "items !!!")

	if len(temp) == 0 {
		return nil
	} else if len(temp) > 0 && len(temp) < 10 {
		return &temp
	}
	// Only the first ten items
	temp = temp[:10]

	// Returns news from Ekantipur
	return &temp
}
