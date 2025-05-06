package types

import (
	"net/http"
	"net/url"

	"gorm.io/gorm"
	. "maragu.dev/gomponents"
)

type Relator interface {
	LoadRelations(db *gorm.DB) error
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

type Identifier interface {
	TableName() string
	GetId() int
}
