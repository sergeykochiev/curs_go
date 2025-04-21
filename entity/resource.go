package entity

import (
	"database/sql"
	"fmt"
	"html"

	. "github.com/sergeykochiev/curs/backend/gui"
	. "github.com/sergeykochiev/curs/backend/types"
	. "github.com/sergeykochiev/curs/backend/util"
	. "maragu.dev/gomponents"
	_ "maragu.dev/gomponents/components"
	. "maragu.dev/gomponents/html"
)

type ResourceEntity struct {
	Id                int
	Name              string
	Date_last_updated string
	Cost_by_one       float32
	Quantity          int
}

func (e *ResourceEntity) ScanRow(r Scanner) error {
	return r.Scan(&e.Id, &e.Name, &e.Date_last_updated, &e.Cost_by_one, &e.Quantity)
}

func (e *ResourceEntity) GetSelectWhereQuery(where string) string {
	return "select * from \"resource\" " + where
}

func (e *ResourceEntity) Insert(db QueryExecutor) (sql.Result, error) {
	return db.Exec("insert into resource (name, date_last_updated, cost_by_one, quantity) values ($1, $2, $3, $4)", e.Name, GetCurrentTime(), e.Cost_by_one, e.Quantity)
}

func (e *ResourceEntity) Update(db QueryExecutor) (sql.Result, error) {
	return db.Exec("update resource set name = $1, date_last_updated = $2, cost_by_one = $3, quantity = $4 where id = $5", e.Name, GetCurrentTime(), e.Cost_by_one, e.Quantity, e.Id)
}

func (e *ResourceEntity) GetDataRow() Group {
	return Group{
		Div(Class("px-[2px] grid place-items-center"), Text(html.EscapeString(fmt.Sprintf("%d", e.Id)))),
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

func (e ResourceEntity) GetCreateForm() Group {
	return Group{
		InputComponent("text", "", "name", "Название", "", true),
		InputComponent("text", "", "cost_by_one", "Стоимость за единицу", "", true),
	}
}

func (e *ResourceEntity) GetReadableName() string {
	return "Ресурс"
}

func (e *ResourceEntity) GetId() int {
	return e.Id
}

func (e *ResourceEntity) Validate() bool {
	return true
}

func (e *ResourceEntity) GetName() string {
	return "resource"
}
