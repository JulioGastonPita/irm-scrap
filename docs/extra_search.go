package models

import (
	"providers/common"
	"slices"
	"strings"
	"time"
)

//<editor-fold desc="Modelos BBDD">

type ExtraSearchEngine struct {
	ID   int    `gorm:"primaryKey;autoIncrement;column:id" json:"id"`
	Code string `gorm:"type:varchar(20);not null;unique;column:code" json:"code"`
	Name string `gorm:"type:varchar(40);not null;column:name" json:"name"`
}

func (ExtraSearchEngine) TableName() string {
	return "extra_search_engine"
}

type ExtraSearchEnginePlugin struct {
	RP common.RowsProvider
	// Engine es el código del motor de búsqueda (mercurio.extra_search_engine.code)
	Engine string
	// Search se encarga de realizar la búsqueda requerida e insertar en custom_urls las URLs recuperadas.
	Search func(request IRMExtraSearchRequest, searchId int64)
}

func (self *ExtraSearchEnginePlugin) GetQueries(req IRMExtraSearchRequest) ([]IRMExtraSearchQuery, error) {
	rows, err := self.RP.SearchEntity(`specific_search_question`, `search_category = ?`, req.EntityType)
	if err != nil {
		return []IRMExtraSearchQuery{}, err
	}

	questions := make([]IRMExtraSearchQuery, 0)
	for _, row := range rows.Rows {
		if slices.Contains(req.Markets, string(row["market"].([]uint8))) {
			questions = append(questions, IRMExtraSearchQuery{
				Query:  strings.Replace(string(row["question"].([]uint8)), "{0}", req.Query, -1),
				Market: string(row["market"].([]uint8)),
			})
		}
	}
	return questions, nil

	return []IRMExtraSearchQuery{}, nil
}
func (self *ExtraSearchEnginePlugin) GetMarkets() ([]ExtraSearchMarket, error) {
	rows, err := self.RP.RawQuery("SELECT * FROM extra_search_market")
	if err != nil {
		return []ExtraSearchMarket{}, err
	}

	mkts := make([]ExtraSearchMarket, 0)
	for _, r := range rows.Rows {
		mkts = append(mkts, ExtraSearchMarket{
			ID:   r["id"].(int),
			Code: string(r["code"].([]uint8)),
			Name: string(r["name"].([]uint8)),
		})
	}
	return mkts, err
}

type ExtraSearchEntityType struct {
	ID   int    `gorm:"primaryKey;autoIncrement;column:id" json:"id"`
	Code string `gorm:"type:varchar(20);not null;unique;column:code" json:"code"`
	Name string `gorm:"type:varchar(40);not null;column:name" json:"name"`
}

func (ExtraSearchEntityType) TableName() string {
	return "extra_search_entity_type"
}

type ExtraSearchMarket struct {
	ID   int    `gorm:"primaryKey;autoIncrement;column:id" json:"id"`
	Code string `gorm:"type:varchar(6);not null;unique;column:code" json:"code"`
	Name string `gorm:"type:varchar(50);not null;column:name" json:"name"`
}

func (ExtraSearchMarket) TableName() string {
	return "extra_search_market"
}

type ExtraSearch struct {
	ID           int        `gorm:"primaryKey;autoIncrement;column:id" json:"id"`
	SearchDate   time.Time  `gorm:"type:datetime;not null;default:CURRENT_TIMESTAMP;column:search_date" json:"searchDate"`
	Term         string     `gorm:"type:varchar(255);not null;column:term" json:"term"`
	EntityType   string     `gorm:"type:varchar(20);not null;foreignKey:Code;references:extra_search_entity_type;column:entity_type" json:"entityType"`
	NodeID       *int64     `gorm:"type:bigint;foreignKey:ID;references:nodes;column:node_id" json:"nodeId"`
	DateFrom     *time.Time `gorm:"type:date;column:date_from" json:"dateFrom"`
	DateTo       *time.Time `gorm:"type:date;column:date_to" json:"dateTo"`
	SearchEngine string     `gorm:"type:varchar(20);not null;foreignKey:Code;references:extra_search_engine;column:search_engine" json:"searchEngine"`
	Markets      string     `gorm:"type:json;not null;column:markets" json:"markets"`
	UserID       int64      `gorm:"type:bigint;not null;foreignKey:ID;references:users;column:user_id" json:"userId"`
	UserToken    string     `gorm:"type:text;not null;column:user_token" json:"userToken"`
}

func (ExtraSearch) TableName() string {
	return "extra_search"
}

//</editor-fold>

// <editor-fold desc="Modelos DTO">

type IRMExtraSearchRequest struct {
	Query        string     `json:"query"`
	EntityType   string     `json:"entityType"`
	MaxURLs      int        `json:"maxUrls"`
	DateFrom     *time.Time `json:"dateFrom,omitempty"`
	DateTo       *time.Time `json:"dateTo,omitempty"`
	SearchEngine string     `json:"searchEngine"`
	Markets      []string   `json:"markets"`
}

type IRMExtraSearchQuery struct {
	Query  string `json:"query"`
	Market string `json:"market"` // es-ES    en-US etc
}

type IRMExtraSearchResponse struct {
	ID              int        `json:"id"`
	SearchDate      time.Time  `json:"searchDate"`
	Term            string     `json:"term"`
	EntityType      string     `json:"entityType"`
	NodeID          *int64     `json:"nodeId,omitempty"`
	DateFrom        *time.Time `json:"dateFrom,omitempty"`
	DateTo          *time.Time `json:"dateTo,omitempty"`
	SearchEngine    string     `json:"searchEngine"`
	Markets         []string   `json:"markets"`
	PercentProgress int        `json:"percentProgress"`
}

type IRMExtraSearch struct {
	IRMExtraSearchResponse
	Articles []common.ArticleDataDTO `json:"articles"`
}

//</editor-fold>
