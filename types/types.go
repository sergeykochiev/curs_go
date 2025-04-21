package types

import (
	"net/url"

	"gorm.io/gorm"
	. "maragu.dev/gomponents"
)

type CrudEntity interface {
	Get()
	GetOne()
	Insert()
	Update()
	Delete()
}

type HtmlEntity interface {
	HtmlTemplater
	Identifier
}

type HtmlTemplater interface {
	GetTableHeader() Group
	GetDataRow() Group
	GetReadableName() string
	GetEntityPage(recursive bool) Group
	GetCreateForm(db *gorm.DB) Group
}

type Validator interface {
	Validate() bool
}

type FormParser interface {
	ValidateAndParseForm(form url.Values) bool
}

type Identifier interface {
	GetName() string
	GetId() int
}
