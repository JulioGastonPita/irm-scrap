package engines

import "time"

type YahooAPIEngineOptions struct {
	RowsProvider any
	Headers      []Header
}

type YahooAPIEngine struct {
	apiURL       string
	headers      []Header
	Engine       string
	RowsProvider any
	//models.ExtraSearchEnginePlugin
}

type YahooAPIRequest struct {
	Query   string `json:"query"`
	Market  string `json:"market"`
	MaxURLs string `json:"maxUrls,omitempty"`
}

type YahooAPIResponse struct {
	SearchId      int64                   `json:"searchId"`
	OriginalQuery string                  `json:"originalQuery"`
	Market        string                  `json:"market"`
	MaxURLs       int                     `json:"maxUrls,omitempty"`
	DateFrom      *time.Time              `json:"dateFrom,omitempty"`
	DateTo        *time.Time              `json:"dateTo,omitempty"`
	Values        []YahooAPIResponseValue `json:"values"`
}

type YahooAPIResponseValue struct {
	Url      string `json:"url"`
	Title    string `json:"title"`
	Snippet  string `json:"snippet"`
	Position int    `json:"possition"`
}
