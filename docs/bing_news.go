package bing_news

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"providers/common"
	"providers/common/models"
	"slices"
	"strings"
	"sync"
	"time"
)

type BingNewsAPIEngineOptions struct {
	RowsProvider common.RowsProvider
	BingApiKey   string
}

// devuelve un BingNewsAPIEngine
func NewBingNewsAPIEngine(options BingNewsAPIEngineOptions) (*BingNewsAPIEngine, error) {
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

type BingNewsAPIEngine struct {
	apiURL string
	apiKey string
	models.ExtraSearchEnginePlugin
}

func (bnEngine BingNewsAPIEngine) Search(request models.IRMExtraSearchRequest, searchId int64) {
	queries, _ := bnEngine.GetQueries(request)
	client := new(http.Client)
	obtainedURLs := 0
	for _, query := range queries {
		// Crear request
		req, err := http.NewRequest("GET", bnEngine.apiURL, nil)
		if err != nil {
			common.ApiLog(fmt.Sprintf("Ocurrió un error consultando en BING -> %+v <- %v", query, err.Error()))
		}
		// Añadir Token BING
		req.Header.Add("Ocp-Apim-Subscription-Key", bnEngine.apiKey)

		// Añadir params a Request
		param := req.URL.Query()
		param.Add("q", query.Query)
		param.Add("sortBy", "Relevance")
		param.Add("mkt", query.Market)
		param.Add("count", query.MaxURLs)
		param.Add("freshness", "Month")
		req.URL.RawQuery = param.Encode()

		// Enviar request a Bing News API
		resp, err := client.Do(req)
		if err != nil {
			common.ApiLog(fmt.Sprintf("Error al consultar Bing Search API -> %v", err.Error()))
			continue
		}

		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			common.ApiLog(fmt.Sprintf("Error al consultar Bing Search API -> %v", err.Error()))
			continue
		}

		ans := new(BingNewsAPIResponse)
		err = json.Unmarshal(body, &ans)
		if err != nil {
			common.ApiLog(fmt.Sprintf("Error al parsear respuesta de Bing Search API -> %v <- %v", body, err.Error()))
			continue
		}

		obtainedURLs += len(ans.Value)
		wg := sync.WaitGroup{}
		for _, bingURL := range ans.Value {
			customUrl := bingURL.toCustomURL(query, searchId)
			wg.Add(1)
			go func(c map[string]interface{}) {
				_, errorPut := bnEngine.RowsProvider.PutEntity("custom_urls", c)
				if errorPut != nil {
					common.ApiLog(fmt.Sprintf("Error en custom_search %d -> %v", searchId, errorPut.Error()))
				} else {
					common.ApiLog(fmt.Sprintf("Guardada correctamente URL para custom_search %d -> %v", searchId, customUrl["url"].(string)))
				}
				wg.Done()
			}(customUrl)

		}
		wg.Wait()
		// Esperar 1 segundos por las deudas
		time.Sleep(1 * time.Second)
	}
	common.ApiLog(fmt.Sprintf("API GDELT devolvió %d urls", obtainedURLs))
}

func (bnEngine BingNewsAPIEngine) GetQueries(req models.IRMExtraSearchRequest) ([]BingNewsAPIRequest, error) {
	var rows common.RowsResponse
	var err error
	if req.EntityType == "" {
		rows, err = bnEngine.RowsProvider.SearchEntity(`specific_search_question`, `1=1`)
		if err != nil {
			return []BingNewsAPIRequest{}, err
		}
	} else {
		rows, err = bnEngine.RowsProvider.SearchEntity(`specific_search_question`, `search_category = ?`, req.EntityType)
		if err != nil {
			return []BingNewsAPIRequest{}, err
		}
	}

	questions := make([]BingNewsAPIRequest, 0)
	for _, row := range rows.Rows {
		if slices.Contains(req.Markets, string(row["market"].([]uint8))) {
			questions = append(questions, BingNewsAPIRequest{
				Query:   strings.Replace(string(row["question"].([]uint8)), "{0}", req.Query, -1),
				Market:  string(row["market"].([]uint8)),
				MaxURLs: fmt.Sprintf("%d", req.MaxURLs),
			})
		}
	}
	return questions, nil
}
