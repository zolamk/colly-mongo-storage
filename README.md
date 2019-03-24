# colly-mongo-storage
A MongoDB storage back end for the Colly web crawling/scraping framework https://go-colly.org

Example Usage:

```go
package main

import (
	"fmt"

	"github.com/gocolly/colly"
	"github.com/zolamk/colly-mongo-storage/colly/mongo"
)

func main() {

	c := colly.NewCollector()

	storage := &mongo.Storage{
		Database: "colly",
		URI:      "mongodb://127.0.0.1:27017",
	}

	if err := c.SetStorage(storage); err != nil {
		panic(err)
	}

	// Find and visit all links
	c.OnHTML("a[href]", func(e *colly.HTMLElement) {
		e.Request.Visit(e.Attr("href"))
	})

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL)
	})

	c.Visit("http://go-colly.org/")
}

```
