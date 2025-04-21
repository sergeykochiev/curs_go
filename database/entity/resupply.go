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

type ResourceResupplyEntity struct {
	ID             int
	Resource_id    int
	Quantity_added int
	Date           string
	_Resource      ResourceEntity
}

func (e *ResourceResupplyEntity) FetchForSelect(db *gorm.DB)

func (e *ResourceResupplyEntity) GetDataRow() Group {
	return Group{
		Div(Class("px-[2px] grid place-items-center"), Text(html.EscapeString(fmt.Sprintf("%d", e.ID)))),
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

func (e ResourceResupplyEntity) GetCreateForm(db *gorm.DB) Group {
	var res []*ResourceEntity
	db.Find(&res)
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
	return e.ID
}

func (e *ResourceResupplyEntity) GetName() string {
	return "resource_resupply"
}

func (e *ResourceResupplyEntity) ValidateAndParseForm(form url.Values) bool {
	if !form.Has("resource_id") || !form.Has("quantity_added") || !form.Has("date") {
		return false
	}
	var err error
	e.Resource_id, err = strconv.Atoi(form.Get("resource_id"))
	if err != nil {
		return false
	}
	e.Quantity_added, err = strconv.Atoi(form.Get("quantity_added"))
	if err != nil {
		return false
	}
	e.Date = form.Get("date")
	return true
}

func (e *ResourceResupplyEntity) AfterCreate(tx *gorm.DB) (err error) {
	e._Resource.ID = e.Resource_id
	res := tx.First(&e._Resource)
	if res.Error != nil {
		return res.Error
	}
	e._Resource.Quantity += e.Quantity_added
	res = tx.Updates(&e._Resource)
	if res.Error != nil {
		return res.Error
	}
	return
}
