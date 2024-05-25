package engines

import (
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/PuerkitoBio/goquery"
)

func NewGoogleAPIEngine(options BingNewsAPIEngineOptions) (*BingNewsAPIEngine, error) {
	if options.BingApiKey == "" {
		return nil, errors.New("BingApiKey its mandatory")
	}
	if options.RowsProvider == nil {
		return nil, errors.New("RowsProvider its mandatory")
	}

	var bingNewsAPIEngine BingNewsAPIEngine
	bingNewsAPIEngine.RowsProvider = options.RowsProvider
	bingNewsAPIEngine.apiKey = options.BingApiKey
	bingNewsAPIEngine.apiURL = "https://api.bing.microsoft.com/v7.0/news/search"
	bingNewsAPIEngine.Engine = "BING_NEWS"
	return &bingNewsAPIEngine, nil
}

func Search(request models.IRMExtraSearchRequest, searchId int64) {

	// url := "https://www.google.com/search?q=pablo+petrecca&gl=us&hl=en"
	//	url = "https://www.google.com/search?oq=%D9%84%D9%8A%D9%88%D9%86%D9%84&q=%D9%84%D9%8A%D9%88%D9%86%D9%84#ip=1"
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatal(err)
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/101.0.4951.54 Safari/537.36")

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	c := 0
	doc.Find("div.g").Each(func(i int, result *goquery.Selection) {
		//	title := result.Find("h3").First().Text()
		link, _ := result.Find("a").First().Attr("href")
		//	snippet := result.Find(".VwiC3b").First().Text()

		//	fmt.Printf("Title: %s\n", title)
		fmt.Printf("Position: %d Link: %s \n ", c+1, link)
		//	fmt.Printf("Snippet: %s\n", snippet)
		//	fmt.Printf("Position: %d\n", c+1)
		//	fmt.Println()

		c++
	})
}

type GoogleAPIEngineOptions struct {
	RowsProvider common.RowsProvider
	BingApiKey   string
}

type GoogleAPIEngine struct {
	apiURL string
	apiKey string
	models.ExtraSearchEnginePlugin
}
