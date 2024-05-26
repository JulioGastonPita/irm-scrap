package engines

import "time"

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

type GoogleAPIResponse struct {
	SearchId      int64                    `json:"searchId"`
	OriginalQuery string                   `json:"originalQuery"`
	Market        string                   `json:"market"`
	MaxURLs       int                      `json:"maxUrls,omitempty"`
	DateFrom      *time.Time               `json:"dateFrom,omitempty"`
	DateTo        *time.Time               `json:"dateTo,omitempty"`
	Values        []GoogleAPIResponseValue `json:"values"`
}

type GoogleAPIResponseValue struct {
	Url      string `json:"url"`
	Title    string `json:"title"`
	Snippet  string `json:"snippet"`
	Position int    `json:"possition"`
}
