package entity

import (
	"database/sql"
	"fmt"
	"html"

	. "github.com/sergeykochiev/curs/backend/gui"
	. "github.com/sergeykochiev/curs/backend/types"
	. "maragu.dev/gomponents"
	_ "maragu.dev/gomponents/components"
	. "maragu.dev/gomponents/html"
)

type ResourceResupplyEntity struct {
	Id             int
	Resource_id    int
	Quantity_added int
	Date           string
	_Resource      ResourceEntity
}

func (e *ResourceResupplyEntity) ScanRow(r Scanner) error {
	return r.Scan(&e.Id, &e.Resource_id, &e.Quantity_added, &e.Date, &e._Resource.Id, &e._Resource.Name, &e._Resource.Date_last_updated, &e._Resource.Cost_by_one, &e._Resource.Quantity)
}

func (e *ResourceResupplyEntity) GetSelectWhereQuery(where string) string {
	return "select * from resource_resupply left join resource on resource_resupply.resource_id = resource.id " + where
}

func (e *ResourceResupplyEntity) Insert(db QueryExecutor) (sql.Result, error) {
	return db.Exec("insert into resource_resupply (resource_id, quantity_added, date) values ($1, $2, $3)", e.Resource_id, e.Quantity_added, e.Date)
}

// TODO implement me
func (e *ResourceResupplyEntity) Update(db QueryExecutor) (sql.Result, error) {
	return db.Exec("")
}

func (e *ResourceResupplyEntity) GetDataRow() Group {
	return Group{
		Div(Class("px-[2px] grid place-items-center"), Text(html.EscapeString(fmt.Sprintf("%d", e.Id)))),
		TableCell(e._Resource.Name),
		TableCell(fmt.Sprintf("%d", e.Quantity_added)),
		TableCell(e.Date),
		TableCell(fmt.Sprintf("%f", e._Resource.Cost_by_one)),
	}
}

func (e *ResourceResupplyEntity) GetTableHeader() Group {
	return Group{
		Div(Class("px-[2px] grid place-items-center"), Text("ID")),
		TableCell("Ресурс"),
		TableCell("Количество добавлено (единиц)"),
		TableCell("Дата поставки"),
		TableCell("Цена за один"),
	}
}

func (e ResourceResupplyEntity) GetEntityPage(recursive bool) Group {
	return Group{
		LabeledField("Количество добавлено (единиц)", fmt.Sprintf("%d", e.Quantity_added)),
		LabeledField("Дата поставки", e.Date),
		If(recursive, Group{
			RelationCard(fmt.Sprintf("Потрачен ресурс #%d", e.Resource_id), &e._Resource),
		}),
	}
}

func (e ResourceResupplyEntity) GetCreateForm(res []*ResourceEntity) Group {
	return Group{
		SelectComponent(res, "", func(r *ResourceEntity) string { return r.Name }, "Выберите ресурс", "resource_id", true, -1),
		InputComponent("number", "", "quantity_added", "Кол-во добавлено", "", true),
		InputComponent("date", "", "date", "Дата поставки", "", true),
	}
}

func (e *ResourceResupplyEntity) GetReadableName() string {
	return "Поставка ресурса"
}

func (e *ResourceResupplyEntity) Validate() bool {
	return true
}

func (e *ResourceResupplyEntity) GetId() int {
	return e.Id
}

func (e *ResourceResupplyEntity) GetName() string {
	return "resource_resupply"
}
