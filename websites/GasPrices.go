package websites

/*
This file contains methods to fetch Petrol, Diesel and LPG Gas Price
from the official NOC website.
*/

import (
	"log"
	"strconv"
	"strings"

	"github.com/gocolly/colly"
)

// Prices of petroleum products
type PetroleumPrices struct {
	EffectiveDate       string
	Petrol, Diesel, LPG float64
}

type NOCTable struct {
}

func FetchPrices() PetroleumPrices {

	// Official website of Nepal Oil Corporation (NOC)
	const url = "http://noc.org.np/retailprice"

	prices := PetroleumPrices{}
	/*
	 table ID="DataTables_Table_0"
	 Fetch first row only as it is the latest prices
	 effective Date |	effective Time | 	petrol |	diesel	|kerosene	| LPG |	ATF (DP)|	ATF (DF)
	 2022.09.01(2079.05.16) | 	24:00 hrs | 	181.00| 	178.00| 	178.00| 	1800.00| 	190.00| 	1645.00
	*/

	collector := colly.NewCollector(
	// Visit only allowed domains
	//colly.AllowedDomains("nepalnews.com", "https://www.nepalnews.com"),
	)

	// Setting Custom User Agent
	collector.UserAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/42.0.2311.135 Safari/537.36 Edge/12.246"

	collector.OnRequest(func(r *colly.Request) {
		log.Println("Visiting ", r.URL)

	})

	// On finding a tbody, run this function
	collector.OnHTML("tbody", func(h *colly.HTMLElement) {

		// for each "tr" element
		h.ForEach("tr", func(x int, el *colly.HTMLElement) {

			// since we only need the first item as it is latest
			if x == 0 {
				el.ForEach("td", func(i int, el *colly.HTMLElement) {

					switch i {
					case 0:
						// first item is date:
						prices.EffectiveDate = el.Text

					case 2:
						// Third item is petrol:

						temp, _ := strconv.ParseFloat(trimAway(el.Text), 64)
						prices.Petrol = temp
					case 3:
						// Fourth item is diesel:
						temp, _ := strconv.ParseFloat(trimAway(el.Text), 64)

						prices.Diesel = temp
					case 5:
						// Sixth item is lpg price:
						temp, _ := strconv.ParseFloat(trimAway(el.Text), 64)

						prices.LPG = temp

					}
				})
			}
		})
	})

	collector.Visit(url)

	collector.Wait()

	log.Println(url, ": Data scraping completion!!")
	return prices

}

// Removes unnecessary \n in text data extracted from table
func trimAway(s string) string {
	return strings.TrimSpace(strings.ReplaceAll(s, "\n", ""))
}
