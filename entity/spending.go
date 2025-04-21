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

type ResourceSpendingEntity struct {
	Id             int
	Order_id       int
	Resource_id    int
	Quantity_spent int
	Date           string
	_Order         OrderEntity
	_Resource      ResourceEntity
}

func (e *ResourceSpendingEntity) ScanRow(r Scanner) error {
	return r.Scan(&e.Id, &e.Order_id, &e.Resource_id, &e.Quantity_spent, &e.Date, &e._Order.Id, &e._Order.Name, &e._Order.Client_name, &e._Order.Client_phone, &e._Order.Date_created, &e._Order.Creator_id, &e._Order.Date_ended, &e._Order.Ended, &e._Resource.Id, &e._Resource.Name, &e._Resource.Date_last_updated, &e._Resource.Cost_by_one, &e._Resource.Quantity)
}

func (e *ResourceSpendingEntity) GetSelectWhereQuery(where string) string {
	return "select * from resource_spending left join \"order\" on \"order\".id = resource_spending.order_id left join resource on resource.id = resource_spending.resource_id " + where
}

func (e *ResourceSpendingEntity) Insert(db QueryExecutor) (sql.Result, error) {
	return db.Exec("insert into resource_spending (order_id, resource_id, quantity_spent, date) values ($1, $2, $3, $4)", e.Order_id, e.Resource_id, e.Quantity_spent, e.Date)
}

// TODO implement me
func (e *ResourceSpendingEntity) Update(db QueryExecutor) (sql.Result, error) {
	return db.Exec("")
}

func (e *ResourceSpendingEntity) GetDataRow() Group {
	return Group{
		Div(Class("px-[2px] grid place-items-center"), Text(html.EscapeString(fmt.Sprintf("%d", e.Id)))),
		TableCell(e._Order.Name),
		TableCell(e._Resource.Name),
		TableCell(fmt.Sprintf("%d", e.Quantity_spent)),
		TableCell(e.Date),
	}
}

func (e *ResourceSpendingEntity) GetTableHeader() Group {
	return Group{
		Div(Class("px-[2px] grid place-items-center"), Text("ID")),
		TableCell("Название заказа"),
		TableCell("Наименование ресурса"),
		TableCell("Количество потрачено (единиц)"),
		TableCell("Дата траты"),
	}
}

func (e ResourceSpendingEntity) GetEntityPage(recursive bool) Group {
	return Group{
		LabeledField("Количество потрачено (единиц)", fmt.Sprintf("%d", e.Quantity_spent)),
		LabeledField("Дата траты", e.Date),
		If(recursive, Group{
			RelationCard(fmt.Sprintf("Потрачено на заказ #%d", e.Order_id), &e._Order),
			RelationCard(fmt.Sprintf("Потрачен ресурс #%d", e.Resource_id), &e._Resource),
		}),
	}
}

func (e ResourceSpendingEntity) GetCreateForm(ord []*OrderEntity, res []*ResourceEntity) Group {
	return Group{
		SelectComponent(ord, "", func(r *OrderEntity) string { return r.Name }, "Выберите заказ, на который был потрачен ресурс", "order_id", true, -1),
		SelectComponent(res, "", func(r *ResourceEntity) string { return r.Name }, "Выберите ресурс", "resource_id", true, -1),
		InputComponent("number", "", "quantity_spent", "Кол-во потрачено", "", true),
		InputComponent("date", "", "date", "Дата траты", "", true),
	}
}

func (e *ResourceSpendingEntity) GetReadableName() string {
	return "Трата ресурса"
}

func (e *ResourceSpendingEntity) Validate() bool {
	return true
}

func (e *ResourceSpendingEntity) GetId() int {
	return e.Id
}

func (e *ResourceSpendingEntity) GetName() string {
	return "resource_spending"
}
