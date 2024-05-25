package engines

import (
	"fmt"
	"log"
	"net/http"
	"net/url"

	"github.com/JulioGastonPita/irm-scrap.git/models"

	"github.com/PuerkitoBio/goquery"
)

func NewGoogleAPIEngine(options GoogleAPIEngineOptions) (*GoogleAPIEngine, error) {

	if len(options.Headers) == 0 {
		options.Headers = append(options.Headers, Header{Key: "User-Agent", Value: "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/101.0.4951.54 Safari/537.36"}) // Handle the case when options.Headers has a length of zero
	}

	var googleAPIEngine GoogleAPIEngine
	googleAPIEngine.RowsProvider = options.RowsProvider
	googleAPIEngine.headers = options.Headers
	googleAPIEngine.apiURL = "https://www.google.com"
	googleAPIEngine.Engine = "GOOGLE"
	return &googleAPIEngine, nil
}

func (bnEngine GoogleAPIEngine) Search(request models.IRMExtraSearchRequest, searchId int64) {

	// armo la url de busqueda, y codifico para URL segura
	urlQuery := fmt.Sprintf("%s/search?q=%s", bnEngine.apiURL, url.QueryEscape(request.Query))

	// si vienen fechas especificadas aplico el filtro
	// tbs=cdr:1,cd_min:2/6/2024,cd_max:2/9/2024
	// cdr:1: Indica que se está utilizando un rango de fechas personalizado.
	if request.DateFrom != nil && request.DateTo != nil {
		urlQuery += "&" + fmt.Sprintf("tbs=cdr:1,cd_min:%s,cd_max:%s", url.QueryEscape(request.DateFrom.Format("01/02/2006")), url.QueryEscape(request.DateTo.Format("01/02/2006")))
	}

	req, err := http.NewRequest("GET", urlQuery, nil)
	if err != nil {
		log.Fatal(err)
	}

	// Agega los headers al request
	for _, header := range bnEngine.headers {
		req.Header.Set(header.Key, header.Value)
	}

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

	// recorro las urls de resultados
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
	RowsProvider any
	Headers      []Header
}

type GoogleAPIEngine struct {
	apiURL       string
	headers      []Header
	Engine       string
	RowsProvider any
	//models.ExtraSearchEnginePlugin
}

type Header struct {
	Key   string
	Value string
}

type GoogleAPIRequest struct {
	Query   string `json:"query"`
	Market  string `json:"market"`
	MaxURLs string `json:"maxUrls,omitempty"`
}
