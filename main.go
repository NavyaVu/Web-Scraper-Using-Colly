package main

import (
	"fmt"

	"github.com/gocolly/colly/v2"
	"github.com/samber/lo"
)

// initializing a data structure to keep the scraped data
// type PokemonProduct struct {
// 	url, image, name, price string
// }

type product struct {
	productUrl   string
	productImage string
	productName  string
}

func main() {
	c := colly.NewCollector(colly.MaxDepth(1), colly.Async(true))
	c.Limit(&colly.LimitRule{
		// limit the parallel requests to 4 request at a time
		Parallelism: 4,
	})

	// c.OnHTML("div h1", func(e *colly.HTMLElement) {
	// 	fmt.Println("text:", e.Text)
	// })
	// c.OnHTML("tbody tr.cmc-table-row", func(h *colly.HTMLElement) {
	// 	fmt.Println(h.ChildText(".cmc-table__column-name--name"))
	// 	// fmt.Println(h.ChildText(".cmc-table__cell--sort-by__symbol"))
	// 	//fmt.Println(h.ChildText(".cmc-table__cell--sort-by__percent-change-1-h"))
	// 	//cmc--change-negative
	// })

	var products []product
	c.OnHTML("ul.products li", func(h *colly.HTMLElement) {
		product := product{
			productUrl:   h.ChildAttr("a", "href"),
			productImage: h.ChildAttr("img", "src"),
			productName:  h.ChildText(".woocommerce-loop-product__title"),
		}
		products = append(products, product)
	})

	var next_pages, pages []string
	c.OnHTML("ul.page-numbers", func(h *colly.HTMLElement) {

		//fmt.Println("---------")
		var next_pages_internal []string = h.ChildAttrs("a", "href")

		//pages := lo.Uniq[string](next_pages_internal)
		pages = append(pages, next_pages_internal...)
		next_pages = lo.Uniq[string](pages)
		for _, pageToScrape := range next_pages {
			if pageToScrape != "https://scrapeme.live/shop/page/1/" {
				c.Visit(pageToScrape)
			}
		}
	})

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL)
	})

	c.Visit("https://scrapeme.live/shop/")
	c.Wait()

	fmt.Println(next_pages, len(next_pages), "Length of pro", len(products))

}
