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

func (e *ResourceEntity) GetDataRow() Group {
	return Group{
		Div(Class("px-[2px] grid place-items-center"), Text(html.EscapeString(fmt.Sprintf("%d", e.ID)))),
		TableCell(e.Name),
		TableCell(e.Date_last_updated),
		TableCell(fmt.Sprintf("%f", e.Cost_by_one)),
		TableCell(fmt.Sprintf("%d", e.Quantity)),
	}
}

func (e *ResourceEntity) GetTableHeader() Group {
	return Group{
		Div(Class("px-[2px] grid place-items-center"), Text("ID")),
		TableCell("Наименование"),
		TableCell("Дата обновления"),
		TableCell("Цена за единицу"),
		TableCell("Количество"),
	}
}

func (e ResourceEntity) GetEntityPage(recursive bool) Group {
	return Group{
		LabeledField("Наименование", e.Name),
		LabeledField("Дата обновления", e.Date_last_updated),
		LabeledField("Цена за единицу", fmt.Sprintf("%f", e.Cost_by_one)),
		LabeledField("Количество", fmt.Sprintf("%d", e.Quantity)),
	}
}

func (e ResourceEntity) GetCreateForm(db *gorm.DB) Group {
	return Group{
		InputComponent("text", "", "name", "Название", "", true),
		InputComponent("text", "", "cost_by_one", "Стоимость за единицу", "", true),
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

func (e *ResourceEntity) GetName() string {
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
