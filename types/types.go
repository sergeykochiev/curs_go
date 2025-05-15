package types

import (
	"net/http"
	"net/url"

	jwt "github.com/golang-jwt/jwt/v5"
	billgen_types "github.com/sergeykochiev/billgen/types"
	"gorm.io/gorm"
	. "maragu.dev/gomponents"
)

type JwtUserDataClaims struct {
	jwt.RegisteredClaims
	UserId int64
}

type TableTemplater interface {
	GetName() string
	ToTHead() []billgen_types.THData
	ToTRow() []billgen_types.TDData
	GetQuery(bool, bool) string
}

type Preloader interface {
	GetPreloadedDb(db *gorm.DB) *gorm.DB
}

type Entity interface {
	HtmlTemplater
	Identifier
	FormParser
	Filterator
	Preloader
	Writable
}

type HtmlTemplater interface {
	GetFilters() Group
	GetTableHeader() Group
	GetDataRow() Group
	GetReadableName() string
	GetEntityPage(bool) Group
	GetCreateForm(*gorm.DB) Group
	GetEntityPageButtons() Group
}

type Filterator interface {
	GetFilteredDb(url.Values, *gorm.DB) *gorm.DB
}

type Validator interface {
	Validate() bool
}

type FormParser interface {
	ValidateAndParseForm(*http.Request) error
}

type Writable interface {
	Clearer
	IdSetter
}

type Clearer interface {
	Clear()
}

type IdSetter interface {
	SetId(int64)
}

type Identifier interface {
	TableName() string
	GetId() int64
}
