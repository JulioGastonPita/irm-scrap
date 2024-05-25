package bing_news

import (
	"net/url"
	"providers/common"
	"strings"
	"time"
)

type BingNewsAPIRequest struct {
	Query   string `json:"query"`
	Market  string `json:"market"`
	MaxURLs string `json:"maxUrls,omitempty"`
}

type BingNewsAPIResponse struct {
	ReadLink     string `json:"readLink"`
	QueryContext struct {
		OriginalQuery string `json:"originalQuery"`
		AdultIntent   bool   `json:"adultIntent"`
	} `json:"queryContext"`
	TotalEstimatedMatches int `json:"totalEstimatedMatches"`
	Sort                  []struct {
		Name       string `json:"name"`
		ID         string `json:"id"`
		IsSelected bool   `json:"isSelected"`
		URL        string `json:"url"`
	} `json:"sort"`
	Value []BingNewsAPIResponseValue `json:"value"`
}

type BingNewsAPIResponseValue struct {
	Name  string `json:"name"`
	URL   string `json:"url"`
	Image struct {
		Thumbnail struct {
			ContentUrl string `json:"thumbnail"`
			Width      int    `json:"width"`
			Height     int    `json:"height"`
		} `json:"thumbnail"`
	} `json:"image"`
	Description string `json:"description"`
	Provider    []struct {
		Type string `json:"_type"`
		Name string `json:"name"`
	} `json:"provider"`
	DatePublished string `json:"datePublished"`
}

func (bnResponse BingNewsAPIResponseValue) toCustomURL(query BingNewsAPIRequest, searchId int64) map[string]interface{} {
	seen, err := common.IRMParseDate(bnResponse.DatePublished)
	if err != nil {
		seen = time.Now()
	}
	language, country := bnResponse.extractCountry(query.Market)
	return map[string]interface{}{
		"url":            bnResponse.URL,
		"title":          bnResponse.Name,
		"language":       bnResponse.parseLanguage(language),
		"source_country": bnResponse.parseCountry(country),
		"domain":         bnResponse.extractDomain(bnResponse.URL),
		"seen":           seen,
		"search_id":      searchId,
		"image_url":      bnResponse.Image.Thumbnail.ContentUrl,
	}
}

func (bnResponse BingNewsAPIResponseValue) extractCountry(market string) (string, string) {
	s := strings.Split(market, "-")
	return strings.ToLower(s[0]), strings.ToLower(s[1])
}

func (bnResponse BingNewsAPIResponseValue) parseLanguage(code string) string {
	switch strings.ToLower(code) {
	case "es":
		return "Spanish"
	case "en":
		return "English"
	case "ko":
		return "Korean"
	case "ru":
		return "Russian"
	case "zh":
		return "Chinese"
	default:
		return "Unknown"
	}
}

func (bnResponse BingNewsAPIResponseValue) parseCountry(code string) string {
	switch strings.ToLower(code) {
	case "es":
		return "Spain"
	case "us":
		return "United States"
	case "kr":
		return "Korea"
	case "ru":
		return "Russian Federation"
	case "cn":
		return "People's republic of China"
	default:
		return "Unknown"
	}
}

func (bnResponse BingNewsAPIResponseValue) extractDomain(urlStr string) string {
	parsedURL, err := url.Parse(urlStr)
	if err != nil {
		return "Error al recuperar dominio"
	}
	return strings.TrimPrefix(parsedURL.Hostname(), "www.")
}
