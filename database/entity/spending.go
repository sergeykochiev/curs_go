package entity

import (
	"errors"
	"fmt"
	"html"
	"net/http"
	"net/url"
	"strconv"

	. "github.com/sergeykochiev/curs/backend/gui"
	"gorm.io/gorm"
	. "maragu.dev/gomponents"
	_ "maragu.dev/gomponents/components"
	. "maragu.dev/gomponents/html"
)

type OrderResourceSpendingEntity struct {
	ID             int
	Order_id       int
	Resource_id    int
	Quantity_spent int
	Date           string
	OrderEntity    OrderEntity    `gorm:"foreignKey:Order_id"`
	ResourceEntity ResourceEntity `gorm:"foreignKey:Resource_id"`
}

func (e OrderResourceSpendingEntity) GetFilters() Group {
	return Group{
		DateFilterComponent("Дата в диапазоне", "date"),
		StringFilterComponent("Название заказа включает", "order_name"),
		StringFilterComponent("Название ресурса включает", "resource_name"),
	}
}

func (e *OrderResourceSpendingEntity) GetFilteredDb(filters url.Values, db *gorm.DB) *gorm.DB {
	if filters.Has("date_lo") && filters.Get("date_lo") != "" {
		db = db.Where("date > ?", filters.Get("date_lo"))
	}
	if filters.Has("date_hi") && filters.Get("date_hi") != "" {
		db = db.Where("date < ?", filters.Get("date_hi"))
	}
	if filters.Has("order_name") && filters.Get("order_name") != "" {
		db = db.Joins("OrderEntity").Where("OrderEntity__name LIKE ?", "%"+filters.Get("order_name")+"%")
	}
	if filters.Has("resource_name") && filters.Get("resource_name") != "" {
		db = db.Joins("ResourceEntity").Where("ResourceEntity__name LIKE ?", "%"+filters.Get("resource_name")+"%")
	}
	return db.Joins("ResourceEntity").Joins("OrderEntity")
}

func (e OrderResourceSpendingEntity) GetDataRow() Group {
	return Group{
		TableDataComponent(html.EscapeString(fmt.Sprintf("%d", e.ID)), Td, fmt.Sprintf("/resource_spending/%d", e.ID)),
		TableDataComponent(e.OrderEntity.Name, Td, ""),
		TableDataComponent(e.ResourceEntity.Name, Td, ""),
		TableDataComponent(fmt.Sprintf("%d", e.Quantity_spent), Td, ""),
		TableDataComponent(e.Date, Td, ""),
	}
}

func (e OrderResourceSpendingEntity) GetTableHeader() Group {
	return Group{
		TableDataComponent("ID", Th, ""),
		TableDataComponent("Название заказа", Th, ""),
		TableDataComponent("Наименование ресурса", Th, ""),
		TableDataComponent("Количество потрачено (единиц)", Th, ""),
		TableDataComponent("Дата траты", Th, ""),
	}
}

func (e OrderResourceSpendingEntity) GetEntityPage(recursive bool) Group {
	return Group{
		LabeledFieldComponent("Количество потрачено (единиц)", fmt.Sprintf("%d", e.Quantity_spent)),
		LabeledFieldComponent("Дата траты", e.Date),
		If(recursive, Group{
			RelationCardComponent(fmt.Sprintf("Потрачено на заказ #%d", e.Order_id), &e.OrderEntity),
			RelationCardComponent(fmt.Sprintf("Потрачен ресурс #%d", e.Resource_id), &e.ResourceEntity),
		}),
	}
}

func (e OrderResourceSpendingEntity) GetCreateForm(db *gorm.DB) Group {
	var ord []*OrderEntity
	var res []*ResourceEntity
	db.Table("order").Find(&ord)
	db.Table("resource").Find(&res)
	return Group{
		SelectComponent(ord, "", func(r *OrderEntity) string { return r.Name }, "Выберите заказ, на который был потрачен ресурс", "order_id", true, -1),
		SelectComponent(res, "", func(r *ResourceEntity) string { return r.Name }, "Выберите ресурс", "resource_id", true, -1),
		LabeledInputComponent("number", "", "quantity_spent", "Кол-во потрачено", "", true),
		LabeledInputComponent("date", "", "date", "Дата траты", "", true),
	}
}

func (e OrderResourceSpendingEntity) GetReadableName() string {
	return "Трата ресурса на заказ"
}

func (e *OrderResourceSpendingEntity) Validate() bool {
	return true
}

func (e OrderResourceSpendingEntity) GetId() int {
	return e.ID
}

func (e OrderResourceSpendingEntity) TableName() string {
	return "order_resource_spending"
}

func (e *OrderResourceSpendingEntity) ValidateAndParseForm(r *http.Request) bool {
	form := r.Form
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

func (e *OrderResourceSpendingEntity) AfterCreate(tx *gorm.DB) (err error) {
	e.ResourceEntity.ID = e.Resource_id
	res := tx.First(&e.ResourceEntity)
	if res.Error != nil {
		return res.Error
	}
	if e.ResourceEntity.Quantity < e.Quantity_spent {
		return errors.New("quantity_spent is more then resource quantity")
	}
	e.ResourceEntity.Quantity -= e.Quantity_spent
	res = tx.Updates(&e.ResourceEntity)
	if res.Error != nil {
		return res.Error
	}
	return
}
