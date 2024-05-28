package engines

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"sync"
	"time"

	"github.com/JulioGastonPita/irm-scrap.git/models"

	"github.com/PuerkitoBio/goquery"
)

func NewYandexAPIEngine(options YandexAPIEngineOptions) (*YandexAPIEngine, error) {

	if len(options.Headers) == 0 {
		options.Headers = append(options.Headers, Header{Key: "User-Agent", Value: "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/101.0.4951.54 Safari/537.36"}) // Handle the case when options.Headers has a length of zero
	}

	var yandexAPIEngine YandexAPIEngine
	yandexAPIEngine.RowsProvider = options.RowsProvider
	yandexAPIEngine.headers = options.Headers
	yandexAPIEngine.apiURL = "https://yandex.ru"
	yandexAPIEngine.Engine = "YANDEX"
	return &yandexAPIEngine, nil
}

func (bnEngine YandexAPIEngine) Search(request models.IRMExtraSearchRequest, searchId int64) {

	// armo la url de busqueda, y codifico para URL segura
	urlQuery := fmt.Sprintf("%s/search?text=%s&search_source=dzen_desktop_safe&msid=1716747569842032-4082553620064597634-rtc-%s-suggest-dzen-ru-vla-2-BAL&suggest_reqid=964597634171674756976350635643099", bnEngine.apiURL, url.QueryEscape(request.Query), url.QueryEscape(request.Query))

	//	urlQuery = "https://yandex.ru/search/?text=Лионель&search_source= dzen_desktop_safe&lr= 10133"

	// si vienen fechas especificadas aplico el filtro
	// tbs=cdr:1,cd_min:2/6/2024,cd_max:2/9/2024
	// cdr:1: Indica que se está utilizando un rango de fechas personalizado.
	// if request.DateFrom != nil && request.DateTo != nil {
	// 	urlQuery += "&" + fmt.Sprintf("tbs=cdr:1,cd_min:%s,cd_max:%s", url.QueryEscape(request.DateFrom.Format("01/02/2006")), url.QueryEscape(request.DateTo.Format("01/02/2006")))
	// }

	// // cargo la cantidad solicitadas de Urls de respuestas
	// if request.MaxURLs == 0 {
	// 	request.MaxURLs = 10
	// }
	// urlQuery += fmt.Sprintf("&start=1&num=%d", request.MaxURLs)

	// // cargo el Market si viene especificado
	// if len(request.Markets) > 0 {
	// 	urlQuery += "&" + fmt.Sprintf("cr=country%s", request.Markets[0])
	// }

	// creo el request
	req, err := http.NewRequest("GET", urlQuery, nil)
	if err != nil {
		FalseApiLog(fmt.Sprintf("Error al Crear Request -> %v", err.Error()))
		return
	}

	// Agega los headers al request
	// for _, header := range bnEngine.headers {
	// 	req.Header.Set(header.Key, header.Value)
	// }

	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7")
	req.Header.Set("Accept-Encoding", "gzip, deflate, br, zstd")
	req.Header.Set("Accept-Language", "es-ES,es;q=0.9,en;q=0.8,gl;q=0.7,it;q=0.6,la;q=0.5,pl;q=0.4,pt;q=0.3")
	req.Header.Set("Device-Memory", "8")
	req.Header.Set("Downlink", "10")
	req.Header.Set("Dpr", "1")
	req.Header.Set("Ect", "4g")
	req.Header.Set("Priority", "u=0, i")
	req.Header.Set("Rtt", "200")
	req.Header.Set("Sec-Ch-Ua", "\"Google Chrome\";v=\"125\", \"Chromium\";v=\"125\", \"Not.A/Brand\";v=\"24\"")
	req.Header.Set("Sec-Ch-Ua-Arch", "\"\"")
	req.Header.Set("Sec-Ch-Ua-Bitness", "\"64\"")
	req.Header.Set("Sec-Ch-Ua-Full-Version", "\"125.0.6422.112\"")
	req.Header.Set("Sec-Ch-Ua-Full-Version-List", "\"Google Chrome\";v=\"125.0.6422.112\", \"Chromium\";v=\"125.0.6422.112\", \"Not.A/Brand\";v=\"24.0.0.0\"")
	req.Header.Set("Sec-Ch-Ua-Mobile", "?1")
	req.Header.Set("Sec-Ch-Ua-Model", "\"Nexus 5\"")
	req.Header.Set("Sec-Ch-Ua-Platform", "\"Android\"")
	req.Header.Set("Sec-Ch-Ua-Platform-Version", "\"6.0\"")
	req.Header.Set("Sec-Ch-Ua-Wow64", "?0")
	req.Header.Set("Sec-Fetch-Dest", "document")
	req.Header.Set("Sec-Fetch-Mode", "navigate")
	req.Header.Set("Sec-Fetch-Site", "none")
	req.Header.Set("Sec-Fetch-User", "?1")
	req.Header.Set("Upgrade-Insecure-Requests", "1")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Linux; Android 6.0; Nexus 5 Build/MRA58N) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/125.0.0.0 Mobile Safari/537.36")
	req.Header.Set("Viewport-Width", "309")

	cookies := []*http.Cookie{
		{
			Name:    "_yasc",
			Value:   "QKuPHYJ0+wgZ5huOFmwtgoPivxjUU5YHqJKJDdNBtz0EFLu+rifGLZTmoVIdyeR7eK4Z",
			Path:    "/",
			Domain:  "yandex.ru",
			Secure:  true,
			Expires: parseTime("Fri, 26 May 2034 12:03:13 GMT"),
		},
		{
			Name:     "i",
			Value:    "tfkyYaR/BTPGWQeAzND7flgfSuQYAPkhtOKL289ZyCMOYK0Eb6vuVEV91G9D6xnDHnozRviRBGhMEB1PcaKkj0wfrLk=",
			Path:     "/",
			Domain:   "yandex.ru",
			Secure:   true,
			HttpOnly: true,
			Expires:  parseTime("Wed, 27 May 2026 23:49:07 GMT"),
		},
		{
			Name:    "is_gdpr",
			Value:   "0",
			Path:    "/",
			Domain:  "yandex.ru",
			Secure:  true,
			Expires: parseTime("Wed, 27 May 2026 23:49:07 GMT"),
		},
		{
			Name:    "is_gdpr_b",
			Value:   "CMWYJRC6/gEoAg==",
			Path:    "/",
			Domain:  "yandex.ru",
			Secure:  true,
			Expires: parseTime("Wed, 27 May 2026 23:49:07 GMT"),
		},
		{
			Name:     "receive-cookie-deprecation",
			Value:    "1",
			Path:     "/",
			Domain:   "yandex.ru",
			Secure:   true,
			HttpOnly: true,
			Expires:  parseTime("Tue, 27 May 2025 23:49:07 GMT"),
		},
		{
			Name:     "set cookie: spravka",
			Value:    "dD0xNjg1MzE2OTA4O2k9MTg2LjEwOC4xMjcuMTgwO0Q9OUUxNzYwNjU0QjJGNzRBODZBNDg0RkFEOTYxRjJGNDZDQkVFMDc4NENGQTg1REFCNjUzMDNEMDZENjlCNTIzM0YxOTIxNDQxMERCNjhBNzM7dT0xNjg1MzE2OTA4MTMyNTQ1ODMxO2g9YzZlMWQwZWI0MzE1MTkwMGM2YWYyMGEzYjQ4OWNmZmQ=",
			Path:     "/",
			Domain:   "yandex.ru",
			Secure:   true,
			HttpOnly: true,
			Expires:  parseTime("Tue, 01 Jul 2025 23:35:08 GMT"),
		},
		{
			Name:    "yandexuid",
			Value:   "4313490991716853747",
			Path:    "/",
			Domain:  "yandex.ru",
			Secure:  true,
			Expires: parseTime("Wed, 27 May 2026 23:49:07 GMT"),
		},
		{
			Name:     "yashr",
			Value:    "1335658141716853747",
			Path:     "/",
			Domain:   "yandex.ru",
			Secure:   true,
			HttpOnly: true,
			Expires:  parseTime("Tue, 27 May 2025 23:49:07 GMT"),
		},
		{
			Name:   "yp",
			Value:  "1722081793.atds.1",
			Path:   "/",
			Domain: "yandex.ru",
			Secure: true,
		},
		{
			Name:   "ys",
			Value:  "wprid.1716897793034390-16752849602344619355-balancer-l7leveler-kubr-yp-sas-57-BAL",
			Path:   "/",
			Domain: "yandex.ru",
			Secure: true,
		},
	}

	for _, cookie := range cookies {
		//fmt.Println(cookie)
		req.AddCookie(cookie)
	}

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		FalseApiLog(fmt.Sprintf("Error al Crear Request -> %v", err.Error()))
		return
	}
	defer res.Body.Close()

	// imprimir en pantalla en modo texto el res.Body
	bodyText, err := ioutil.ReadAll(res.Body)
	if err != nil {
		FalseApiLog(fmt.Sprintf("Error al leer el cuerpo de la respuesta -> %v", err.Error()))
		return
	}

	fmt.Println(string(bodyText))

	// cargo el documento de respuesta
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		FalseApiLog(fmt.Sprintf("Error al consultar YANDEX Search -> %v", err.Error()))
		return
	}

	// cargo la response (encabezado)
	yandexAPIResponse := &YandexAPIResponse{
		SearchId:      searchId,
		OriginalQuery: request.Query,
		Market:        request.Markets[0],
		MaxURLs:       request.MaxURLs,
		DateFrom:      request.DateFrom,
		DateTo:        request.DateTo,
		Values:        make([]YandexAPIResponseValue, 0),
	}

	// recorro las urls de resultados
	c := 0                 // contador de urls obtenidas
	wg := sync.WaitGroup{} // wait Group para los inserts a la base de datos

	doc.Find("div.VanillaReact.OrganicTitle.OrganicTitle_multiline.Typo.Typo_text_l.organic__title-wrapper").Each(func(i int, result *goquery.Selection) {

		// obtengo los valores de la url, titulo y snippet
		link, _ := result.Find("a").First().Attr("href")
		title := result.Find("h3").First().Text()
		snippet := result.Find(".VwiC3b").First().Text()
		c++

		// cargo los valores en la respuesta
		yandexAPIResponse.Values = append(yandexAPIResponse.Values, YandexAPIResponseValue{
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

func parseTime(timeStr string) time.Time {
	layout := "Mon, 02 Jan 2006 15:04:05 GMT"
	t, err := time.Parse(layout, timeStr)
	if err != nil {
		fmt.Println("Error parsing time:", err)
	}
	return t
}
