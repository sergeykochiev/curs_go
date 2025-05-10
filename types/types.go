package types

import (
	"net/http"
	"net/url"

	"gorm.io/gorm"
	. "maragu.dev/gomponents"
)

type ResourceSpending struct {
	Name          string
	Last_date     string
	One_is_called string
	Count_spent   float32
}

type ItemPopularity struct {
	Name            string
	Last_date       string
	Count_fulfilled int
}

type Preloader interface {
	GetPreloadedDb(db *gorm.DB) *gorm.DB
}

type HtmlTemplater interface {
	GetFilters() Group
	GetTableHeader() Group
	GetDataRow() Group
	GetReadableName() string
	GetEntityPage(recursive bool) Group
	GetCreateForm(db *gorm.DB) Group
	GetEntityPageButtons() Group
}

type Filterator interface {
	GetFilteredDb(values url.Values, db *gorm.DB) *gorm.DB
}

type Validator interface {
	Validate() bool
}

type FormParser interface {
	ValidateAndParseForm(r *http.Request) error
}

type Writable interface {
	Clear()
	SetId(id int64)
}

type Identifier interface {
	TableName() string
	GetId() int64
}
