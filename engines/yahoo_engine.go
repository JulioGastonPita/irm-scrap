package engines

import (
	"fmt"
	"net/http"
	"net/url"
	"regexp"
	"strings"
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
	YahooAPIEngine.apiURL = "search.yahoo.com"
	YahooAPIEngine.Engine = "Yahoo"
	return &YahooAPIEngine, nil
}

func (bnEngine YahooAPIEngine) Search(request models.IRMExtraSearchRequest, searchId int64) {

	// si no se especifica, el markets
	if request.Markets == nil {
		request.Markets = []string{"US"}
	}

	// armo la url de busqueda, y codifico para URL segura
	urlQuery := fmt.Sprintf("https://%s.%s/search?p=%s", request.Markets[0], bnEngine.apiURL, url.QueryEscape(request.Query))

	// armo la url de busqueda, y codifico para URL segura

	// si vienen fechas especificadas aplico el filtro
	// tbs=cdr:1,cd_min:2/6/2024,cd_max:2/9/2024
	// cdr:1: Indica que se está utilizando un rango de fechas personalizado.
	if request.DateFrom != nil && request.DateTo != nil {
		urlQuery += fmt.Sprintf("&fr=yfp-t&bt=%s&et=%s", request.DateFrom.Format("20060102"), request.DateTo.Format("20060102"))
	}

	// Debido a que no funcion el conteo de resultados, voy haciendo multiples llamadas desde hasta
	// hasta obtener la cantidad de resultados solicitados
	// el tope de resultados por pagina es de 7,
	// no funciona con mas de esto
	maxResultsPage := 7
	paginado := 1
	c := 0 // contador de resultados

	wg := sync.WaitGroup{} // wait Group para los inserts a la base de datos

	for {

		// armo los parametros del paginado
		pageParameters := "&b=" + fmt.Sprintf("%d", paginado) + fmt.Sprintf("&pz=%d", maxResultsPage)

		// creo el request
		req, err := http.NewRequest("GET", urlQuery+pageParameters, nil)
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

		// // imprimir en pant

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

		doc.Find(".algo-sr").Each(func(i int, result *goquery.Selection) {

			// obtengo los valores de la url, titulo y snippet
			link, _ := result.Find("a").First().Attr("href")

			title, _ := result.Find("a").First().Attr("aria-label")

			snippet := result.Find("span.fc-falcon").First().Text()
			c++

			// cargo los valores en la respuesta
			YahooAPIResponse.Values = append(YahooAPIResponse.Values, YahooAPIResponseValue{
				Url:      link,
				Title:    title,
				Snippet:  snippet,
				Position: c})

			// println("URL: ", sanitizeURL(link))
			// println("title: ", cleanHTMLTags(title))
			// println("snippet: ", snippet)
			// println("**************************************")

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

		// aumento el paginado y veo si llegué al maximo
		paginado += maxResultsPage

		// Verifico contra el maximo de resultados solicitados
		if paginado >= request.MaxURLs {
			break
		}

	}

	// Esperamos a que todas las goroutines terminen
	wg.Wait()

}
func sanitizeURL(input string) string {
	startIndex := strings.Index(input, "/RU=")
	endIndex := strings.Index(input, "/RK=")
	if startIndex == -1 || endIndex == -1 {
		return ""
	}
	startIndex += 4
	if startIndex >= endIndex {
		return ""
	}
	decodedURL, err := url.QueryUnescape(input[startIndex:endIndex])
	if err != nil {
		return ""
	}
	return decodedURL
}

func cleanHTMLTags(input string) string {
	// Remove HTML tags using regex
	re := regexp.MustCompile("<[^>]*>")
	cleanedString := re.ReplaceAllString(input, "")

	return cleanedString
}
