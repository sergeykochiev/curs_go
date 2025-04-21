package entity

import (
	"errors"
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

type ResourceSpendingEntity struct {
	ID             int
	Order_id       int
	Resource_id    int
	Quantity_spent int
	Date           string
	_Order         OrderEntity
	_Resource      ResourceEntity
}

func (e *ResourceSpendingEntity) GetDataRow() Group {
	return Group{
		Div(Class("px-[2px] grid place-items-center"), Text(html.EscapeString(fmt.Sprintf("%d", e.ID)))),
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

func (e ResourceSpendingEntity) GetCreateForm(db *gorm.DB) Group {
	var ord []*OrderEntity
	var res []*ResourceEntity
	db.Find(&ord)
	db.Find(&res)
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
	return e.ID
}

func (e *ResourceSpendingEntity) GetName() string {
	return "resource_spending"
}

func (e *ResourceSpendingEntity) ValidateAndParseForm(form url.Values) bool {
	if !form.Has("order_id") || !form.Has("resource_id") || !form.Has("quantity_spent") || !form.Has("date") {
		return false
	}
	var err error
	e.Order_id, err = strconv.Atoi(form.Get("order_id"))
	if err != nil {
		return false
	}
	e.Resource_id, err = strconv.Atoi(form.Get("resource_id"))
	if err != nil {
		return false
	}
	e.Quantity_spent, err = strconv.Atoi(form.Get("quantity_spent"))
	if err != nil {
		return false
	}
	e.Date = form.Get("date")
	return true
}

func (e *ResourceSpendingEntity) AfterCreate(tx *gorm.DB) (err error) {
	e._Resource.ID = e.Resource_id
	res := tx.First(&e._Resource)
	if res.Error != nil {
		return res.Error
	}
	if e._Resource.Quantity < e.Quantity_spent {
		return errors.New("quantity_spent is more then resource quantity")
	}
	e._Resource.Quantity -= e.Quantity_spent
	res = tx.Updates(&e._Resource)
	if res.Error != nil {
		return res.Error
	}
	return
}
