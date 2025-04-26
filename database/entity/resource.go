package entity

import (
	"fmt"
	"html"
	"net/url"
	"strconv"

	. "github.com/sergeykochiev/curs/backend/gui"
	"gorm.io/gorm"
	. "maragu.dev/gomponents"
	_ "maragu.dev/gomponents/components"
	. "maragu.dev/gomponents/html"
)

type ResourceEntity struct {
	ID                int
	Name              string
	Date_last_updated string
	Cost_by_one       float32
	Quantity          int
}

func (e *ResourceEntity) GetFilters() Group {
	return Group{
		StringFilterComponent("Название включает", "name"),
		DateFilterComponent("Дата последнего обновления в диапазоне", "date_last_updated"),
	}
}

func (e *ResourceEntity) GetFilteredDb(filters url.Values, db *gorm.DB) *gorm.DB {
	if filters.Has("date_last_updated_lo") && filters.Get("date_last_updated_lo") != "" {
		db = db.Where("date_last_updated > ?", filters.Get("date_last_updated_lo"))
	}
	if filters.Has("date_last_updated_hi") && filters.Get("date_last_updated_hi") != "" {
		db = db.Where("date_last_updated < ?", filters.Get("date_last_updated_hi"))
	}
	if filters.Has("name") && filters.Get("name") != "" {
		db = db.Where("name LIKE ?", "%"+filters.Get("name")+"%")
	}
	return db
}

func (e *ResourceEntity) GetDataRow() Group {
	return Group{
		Div(Class("px-[2px] grid place-items-center"), Text(html.EscapeString(fmt.Sprintf("%d", e.ID)))),
		TableCellComponent(e.Name),
		TableCellComponent(e.Date_last_updated),
		TableCellComponent(fmt.Sprintf("%f", e.Cost_by_one)),
		TableCellComponent(fmt.Sprintf("%d", e.Quantity)),
	}
}

func (e *ResourceEntity) GetTableHeader() Group {
	return Group{
		Div(Class("px-[2px] grid place-items-center"), Text("ID")),
		TableCellComponent("Наименование"),
		TableCellComponent("Дата обновления"),
		TableCellComponent("Цена за единицу"),
		TableCellComponent("Количество"),
	}
}

func (e ResourceEntity) GetEntityPage(recursive bool) Group {
	return Group{
		LabeledFieldComponent("Наименование", e.Name),
		LabeledFieldComponent("Дата обновления", e.Date_last_updated),
		LabeledFieldComponent("Цена за единицу", fmt.Sprintf("%f", e.Cost_by_one)),
		LabeledFieldComponent("Количество", fmt.Sprintf("%d", e.Quantity)),
	}
}

func (e ResourceEntity) GetCreateForm(db *gorm.DB) Group {
	return Group{
		LabeledInputComponent("text", "", "name", "Название", "", true),
		LabeledInputComponent("text", "", "cost_by_one", "Стоимость за единицу", "", true),
	}
}

func (e *ResourceEntity) GetReadableName() string {
	return "Ресурс"
}

func (e *ResourceEntity) GetId() int {
	return e.ID
}

func (e *ResourceEntity) Validate() bool {
	return true
}

func (e *ResourceEntity) TableName() string {
	return "resource"
}

func (e *ResourceEntity) ValidateAndParseForm(form url.Values) bool {
	if !form.Has("name") || !form.Has("cost_by_one") {
		return false
	}
	e.Name = form.Get("name")
	cost_by_one, err := strconv.Atoi(form.Get("cost_by_one"))
	if err != nil {
		return false
	}
	e.Cost_by_one = float32(cost_by_one)
	return true
}
