package main

import (
	"time"

	ge "github.com/JulioGastonPita/irm-scrap.git/engines"
	"github.com/JulioGastonPita/irm-scrap.git/models"
)

func main() {

	/* ********************* YAHOO *********************
	 ************************************************************************************************************ */
	miYahooEngineOptions := ge.YahooAPIEngineOptions{
		Headers: []ge.Header{
			{Key: "User-Agent", Value: "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/101.0.4951.54 Safari/537.36"},
		},
	}

	//  creo el
	miYahooEngine, _ := ge.NewYahooAPIEngine(miYahooEngineOptions)

	dateFrom, _ := time.Parse("2006-01-02", "2023-07-01")
	dateTo, _ := time.Parse("2006-01-02", "2024-07-01")

	//   Cargos algunas opciones de busquedas
	//   para el caso de los Markets, solo se puede especificar uno
	//   debido al requerimiento de Google, en caso de que se necesite mas un un Market ( ISO 3166-1 alfa-2)
	//   debería combinarse los resultados
	search := models.IRMExtraSearchRequest{Query: "ليونل مسي",
		DateFrom: &dateFrom,
		DateTo:   &dateTo,
		Markets:  []string{"AE"}}

	miYahooEngine.Search(search, 12)

	/* ********************* YAHOO *********************
	 ************************************************************************************************************ */

	//  Test
	//	 Opciones de Google,
	//  los headers son opcionales, el valor cargado en el ejemplo es el de omision
	miGoogleEngineOptions := ge.GoogleAPIEngineOptions{
		Headers: []ge.Header{
			{Key: "User-Agent", Value: "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/101.0.4951.54 Safari/537.36"},
		},
	}

	//  creo el
	miGoogleEngine, _ := ge.NewGoogleAPIEngine(miGoogleEngineOptions)

	dateFrom, _ = time.Parse("2006-01-02", "2023-07-01")
	dateTo, _ = time.Parse("2006-01-02", "2024-07-01")

	//   Cargos algunas opciones de busquedas
	//   para el caso de los Markets, solo se puede especificar uno
	//   debido al requerimiento de Google, en caso de que se necesite mas un un Market ( ISO 3166-1 alfa-2)
	//   debería combinarse los resultados
	search = models.IRMExtraSearchRequest{Query: "ليونل مسي",
		DateFrom: &dateFrom,
		DateTo:   &dateTo,
		Markets:  []string{"AE"}}

	miGoogleEngine.Search(search, 12)

	// Test YANDEX
	//Opciones de Google,
	// los headers son opcionales, el valor cargado en el ejemplo es el de omision
	miYandexEngineOptions := ge.YandexAPIEngineOptions{
		Headers: []ge.Header{
			{Key: "User-Agent", Value: "Mozilla/5.0 (Linux; Android 6.0; Nexus 5 Build/MRA58N) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/125.0.0.0 Mobile Safari/537.36"},
		},
	}

	//creo el
	miYandexEngine, _ := ge.NewYandexAPIEngine(miYandexEngineOptions)

	dateFrom, _ = time.Parse("2006-01-02", "2023-07-01")
	dateTo, _ = time.Parse("2006-01-02", "2024-07-01")

	//  Cargos algunas opciones de busquedas
	//  para el caso de los Markets, solo se puede especificar uno
	//  debido al requerimiento de Google, en caso de que se necesite mas un un Market ( ISO 3166-1 alfa-2)
	//  debería combinarse los resultados
	search = models.IRMExtraSearchRequest{Query: "Лионель Месси",
		DateFrom: &dateFrom,
		DateTo:   &dateTo,
		Markets:  []string{"AE"}}

	miYandexEngine.Search(search, 12)

}
