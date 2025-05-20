package entity

import (
	"errors"
	"fmt"
	"html"
	"net/http"
	"net/url"
	"strconv"

	. "github.com/sergeykochiev/curs/backend/gui"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
	. "maragu.dev/gomponents"
	_ "maragu.dev/gomponents/components"
	. "maragu.dev/gomponents/html"
)

type OrderResourceSpendingEntity struct {
	Id             decimal.Decimal `gorm:"primaryKey"`
	Order_id       decimal.Decimal
	Resource_id    decimal.Decimal
	Quantity_spent float64
	Date           string
	OrderEntity    OrderEntity    `gorm:"foreignKey:Order_id"`
	ResourceEntity ResourceEntity `gorm:"foreignKey:Resource_id"`
}

func (e OrderResourceSpendingEntity) GetEntityPageButtons() Group {
	return Group{}
}

func (e OrderResourceSpendingEntity) GetFilters() Group {
	return Group{
		DateFilterComponent("Дата в диапазоне", "date"),
		StringFilterComponent("Название заказа включает", "order_name"),
		StringFilterComponent("Название ресурса включает", "resource_name"),
	}
}

func (e *OrderResourceSpendingEntity) GetPreloadedDb(db *gorm.DB) *gorm.DB {
	return db.Joins("ResourceEntity").Joins("OrderEntity")
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
		TableDataComponent(html.EscapeString(fmt.Sprintf("%d", e.GetId())), Td, fmt.Sprintf("/resource_spending/%d", e.GetId())),
		TableDataComponent(e.OrderEntity.Name, Td, ""),
		TableDataComponent(e.ResourceEntity.Name, Td, ""),
		TableDataComponent(fmt.Sprintf("%f", e.Quantity_spent), Td, ""),
		TableDataComponent(e.Date, Td, ""),
	}
}

func (e OrderResourceSpendingEntity) GetTableHeader() Group {
	return Group{
		TableDataComponent("Id", Th, ""),
		TableDataComponent("Название заказа", Th, ""),
		TableDataComponent("Наименование ресурса", Th, ""),
		TableDataComponent("Количество потрачено (единиц)", Th, ""),
		TableDataComponent("Дата траты", Th, ""),
	}
}

func (e OrderResourceSpendingEntity) GetEntityPage(recursive bool) Group {
	return Group{
		LabeledFieldComponent("Количество потрачено (единиц)", fmt.Sprintf("%f", e.Quantity_spent)),
		LabeledFieldComponent("Дата траты", e.Date),
		If(recursive, Group{
			RelationCardComponent(fmt.Sprintf("Потрачено на заказ #%d", e.OrderEntity.GetId()), &e.OrderEntity),
			RelationCardComponent(fmt.Sprintf("Потрачен ресурс #%d", e.ResourceEntity.GetId()), &e.ResourceEntity),
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

func (e OrderResourceSpendingEntity) GetId() int64 {
	return e.Id.IntPart()
}

func (e *OrderResourceSpendingEntity) SetId(id int64) {
	e.Id = decimal.NewFromInt(id)
}

func (e *OrderResourceSpendingEntity) Clear() {
	*e = OrderResourceSpendingEntity{}
}

func (e OrderResourceSpendingEntity) TableName() string {
	return "order_resource_spending"
}

func (e *OrderResourceSpendingEntity) ValidateAndParseForm(r *http.Request) error {
	form := r.Form
	if !form.Has("order_id") || !form.Has("resource_id") || !form.Has("quantity_spent") || !form.Has("date") {
		return errors.New("Invalid fields")
	}
	var err error
	e.Order_id, err = decimal.NewFromString(form.Get("order_id"))
	if err != nil {
		return err
	}
	e.Resource_id, err = decimal.NewFromString(form.Get("resource_id"))
	if err != nil {
		return err
	}
	quantity_added, err := strconv.ParseFloat(form.Get("quantity_added"), 32)
	if err != nil {
		return err
	}
	e.Quantity_spent = float64(quantity_added)
	e.Date = form.Get("date")
	return nil
}

func (e *OrderResourceSpendingEntity) AfterCreate(tx *gorm.DB) (err error) {
	e.ResourceEntity.Id = e.Resource_id
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
