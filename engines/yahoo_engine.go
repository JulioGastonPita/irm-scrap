package engines

import (
	"fmt"
	"net/http"
	"net/url"
	"sync"
	"time"

	"github.com/JulioGastonPita/irm-scrap.git/models"

	"github.com/PuerkitoBio/goquery"
)

func NewYahooAPIEngine(options YahooAPIEngineOptions) (*YahooAPIEngine, error) {

	if len(options.Headers) == 0 {
		options.Headers = append(options.Headers, Header{Key: "User-Agent", Value: "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/101.0.4951.54 Safari/537.36"}) // Handle the case when options.Headers has a length of zero
	}

	var YahooAPIEngine YahooAPIEngine
	YahooAPIEngine.RowsProvider = options.RowsProvider
	YahooAPIEngine.headers = options.Headers
	YahooAPIEngine.apiURL = "https://www.Yahoo.com"
	YahooAPIEngine.Engine = "Yahoo"
	return &YahooAPIEngine, nil
}

func (bnEngine YahooAPIEngine) Search(request models.IRMExtraSearchRequest, searchId int64) {

	// armo la url de busqueda, y codifico para URL segura
	urlQuery := fmt.Sprintf("%s/search?q=%s", bnEngine.apiURL, url.QueryEscape(request.Query))

	// si vienen fechas especificadas aplico el filtro
	// tbs=cdr:1,cd_min:2/6/2024,cd_max:2/9/2024
	// cdr:1: Indica que se está utilizando un rango de fechas personalizado.
	if request.DateFrom != nil && request.DateTo != nil {
		urlQuery += "&" + fmt.Sprintf("tbs=cdr:1,cd_min:%s,cd_max:%s", url.QueryEscape(request.DateFrom.Format("01/02/2006")), url.QueryEscape(request.DateTo.Format("01/02/2006")))
	}

	// cargo la cantidad solicitadas de Urls de respuestas
	if request.MaxURLs == 0 {
		request.MaxURLs = 10
	}
	urlQuery += fmt.Sprintf("&start=1&num=%d", request.MaxURLs)

	// cargo el Market si viene especificado
	if len(request.Markets) > 0 {
		urlQuery += "&" + fmt.Sprintf("cr=country%s", request.Markets[0])
	}

	urlQuery = "https://search.yahoo.com/search?p=chatgpt&fr=yfp-t&ei=UTF-8&n=10&b=11"

	// creo el request
	req, err := http.NewRequest("GET", urlQuery, nil)
	if err != nil {
		FalseApiLog(fmt.Sprintf("Error al Crear Request -> %v", err.Error()))
		return
	}

	// Agega los headers al request
	for _, header := range bnEngine.headers {
		req.Header.Set(header.Key, header.Value)
	}

	// ejecuto el request
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		FalseApiLog(fmt.Sprintf("Error al Crear Request -> %v", err.Error()))
		return
	}
	defer res.Body.Close()

	// // imprimir en pantalla en modo texto el res.Body
	// bodyText, err := ioutil.ReadAll(res.Body)
	// if err != nil {
	// 	FalseApiLog(fmt.Sprintf("Error al leer el cuerpo de la respuesta -> %v", err.Error()))
	// 	return
	// }

	// fmt.Println(string(bodyText))

	// cargo el documento de respuesta
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		FalseApiLog(fmt.Sprintf("Error al consultar Yahoo Search -> %v", err.Error()))
		return
	}

	// cargo la response (encabezado)
	YahooAPIResponse := &YahooAPIResponse{
		SearchId:      searchId,
		OriginalQuery: request.Query,
		Market:        request.Markets[0],
		MaxURLs:       request.MaxURLs,
		DateFrom:      request.DateFrom,
		DateTo:        request.DateTo,
		Values:        make([]YahooAPIResponseValue, 0),
	}

	// recorro las urls de resultados
	c := 0                 // contador de urls obtenidas
	wg := sync.WaitGroup{} // wait Group para los inserts a la base de datos
	doc.Find("div.dd.fst.algo.algo-sr.relsrch.Sr").Each(func(i int, result *goquery.Selection) {

		// obtengo los valores de la url, titulo y snippet
		link, _ := result.Find("a").First().Attr("href")
		title := result.Find("h3").First().Text()
		snippet := result.Find(".VwiC3b").First().Text()
		c++

		// cargo los valores en la respuesta
		YahooAPIResponse.Values = append(YahooAPIResponse.Values, YahooAPIResponseValue{
			Url:      link,
			Title:    title,
			Snippet:  snippet,
			Position: c})

		/****
		ATENCION!!: Aqui Modificar para guardar en la base de datos
		            Simulo una grabacion, ya que no tengo la base de datos
		****/
		wg.Add(1)
		go func(str string) {
			defer wg.Done() // Marcamos la goroutine como completada al final
			fmt.Println(time.Now().Format("15:04:05"), link)
			time.Sleep(1 * time.Second) // Simulamos una operación que toma tiempo
		}(link)

	})

	// Esperamos a que todas las goroutines terminen
	wg.Wait()

}
