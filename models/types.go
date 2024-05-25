package models

import "time"

type IRMExtraSearchRequest struct {
	Query        string     `json:"query"`
	EntityType   string     `json:"entityType"`
	MaxURLs      int        `json:"maxUrls"`
	DateFrom     *time.Time `json:"dateFrom,omitempty"`
	DateTo       *time.Time `json:"dateTo,omitempty"`
	SearchEngine string     `json:"searchEngine"`
	Markets      []string   `json:"markets"`
}
