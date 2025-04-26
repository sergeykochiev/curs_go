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
	ResourceEntity ResourceEntity `gorm:"foreignKey:Resource_id"`
}

func (e *ResourceResupplyEntity) GetFilters() Group {
	return Group{
		DateFilterComponent("Дата в диапазоне", "date"),
		StringFilterComponent("Название ресурса включает", "resource_name"),
	}
}

func (e *ResourceResupplyEntity) GetFilteredDb(filters url.Values, db *gorm.DB) *gorm.DB {
	if filters.Has("date_lo") && filters.Get("date_lo") != "" {
		db = db.Where("date > ?", filters.Get("date_lo"))
	}
	if filters.Has("date_hi") && filters.Get("date_hi") != "" {
		db = db.Where("date < ?", filters.Get("date_hi"))
	}
	if filters.Has("resource_name") && filters.Get("resource_name") != "" {
		db = db.Where("ResourceEntity__name LIKE ?", "%"+filters.Get("resource_name")+"%")
	}
	return db.Joins("ResourceEntity")
}

func (e *ResourceResupplyEntity) GetDataRow() Group {
	return Group{
		Div(Class("px-[2px] grid place-items-center"), Text(html.EscapeString(fmt.Sprintf("%d", e.ID)))),
		TableCellComponent(e.ResourceEntity.Name),
		TableCellComponent(fmt.Sprintf("%d", e.Quantity_added)),
		TableCellComponent(e.Date),
		TableCellComponent(fmt.Sprintf("%f", e.ResourceEntity.Cost_by_one)),
	}
}

func (e *ResourceResupplyEntity) GetTableHeader() Group {
	return Group{
		Div(Class("px-[2px] grid place-items-center"), Text("ID")),
		TableCellComponent("Ресурс"),
		TableCellComponent("Количество добавлено (единиц)"),
		TableCellComponent("Дата поставки"),
		TableCellComponent("Цена за один"),
	}
}

func (e ResourceResupplyEntity) GetEntityPage(recursive bool) Group {
	return Group{
		LabeledFieldComponent("Количество добавлено (единиц)", fmt.Sprintf("%d", e.Quantity_added)),
		LabeledFieldComponent("Дата поставки", e.Date),
		If(recursive, Group{
			RelationCardComponent(fmt.Sprintf("Поставлен ресурс #%d", e.Resource_id), &e.ResourceEntity),
		}),
	}
}

func (e ResourceResupplyEntity) GetCreateForm(db *gorm.DB) Group {
	var res []*ResourceEntity
	db.Table("resource").Find(&res)
	return Group{
		SelectComponent(res, "", func(r *ResourceEntity) string { return r.Name }, "Выберите ресурс", "resource_id", true, -1),
		LabeledInputComponent("number", "", "quantity_added", "Кол-во добавлено", "", true),
		LabeledInputComponent("date", "", "date", "Дата поставки", "", true),
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

func (e *ResourceResupplyEntity) TableName() string {
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
	e.ResourceEntity.ID = e.Resource_id
	res := tx.First(&e.ResourceEntity)
	if res.Error != nil {
		return res.Error
	}
	e.ResourceEntity.Quantity += e.Quantity_added
	res = tx.Updates(&e.ResourceEntity)
	if res.Error != nil {
		return res.Error
	}
	return
}
