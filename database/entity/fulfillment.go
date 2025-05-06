package entity

import (
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

type OrderItemFulfillmentEntity struct {
	ID                 int
	Order_id           int
	Item_id            int
	Quantity_fulfilled int
	OrderEntity        OrderEntity `gorm:"foreignKey:Order_id"`
	ItemEntity         ItemEntity  `gorm:"foreignKey:Item_id"`
}

func (e OrderItemFulfillmentEntity) GetFilters() Group {
	return Group{
		StringFilterComponent("Название заказа включает", "order_name"),
		StringFilterComponent("Название товара включает", "item_name"),
	}
}

func (e *OrderItemFulfillmentEntity) GetFilteredDb(filters url.Values, db *gorm.DB) *gorm.DB {
	if filters.Has("date_hi") && filters.Get("date_hi") != "" {
		db = db.Where("date < ?", filters.Get("date_hi"))
	}
	if filters.Has("order_name") && filters.Get("order_name") != "" {
		db = db.Joins("OrderEntity").Where("OrderEntity__name LIKE ?", "%"+filters.Get("order_name")+"%")
	}
	if filters.Has("item_name") && filters.Get("item_name") != "" {
		db = db.Joins("ItemEntity").Where("ItemEntity__name LIKE ?", "%"+filters.Get("item_name")+"%")
	}
	return db.Joins("ItemEntity").Joins("OrderEntity")
}

func (e OrderItemFulfillmentEntity) GetDataRow() Group {
	return Group{
		TableDataComponent(html.EscapeString(fmt.Sprintf("%d", e.ID)), Td, fmt.Sprintf("/order_item_fulfillment/%d", e.ID)),
		TableDataComponent(e.OrderEntity.Name, Td, ""),
		TableDataComponent(e.ItemEntity.Name, Td, ""),
		TableDataComponent(fmt.Sprintf("%d", e.Quantity_fulfilled), Td, ""),
	}
}

func (e OrderItemFulfillmentEntity) GetTableHeader() Group {
	return Group{
		TableDataComponent("ID", Th, ""),
		TableDataComponent("Название заказа", Th, ""),
		TableDataComponent("Наименование товара", Th, ""),
		TableDataComponent("Количество предоставлено (единиц)", Th, ""),
	}
}

func (e OrderItemFulfillmentEntity) GetEntityPage(recursive bool) Group {
	return Group{
		LabeledFieldComponent("Количество предоставлено (единиц)", fmt.Sprintf("%d", e.Quantity_fulfilled)),
		If(recursive, Group{
			RelationCardComponent(fmt.Sprintf("Предоставлено в рамках заказа #%d", e.Order_id), &e.OrderEntity),
			RelationCardComponent(fmt.Sprintf("Предоставлен товар #%d (%d шт.)", e.Item_id, e.Quantity_fulfilled), &e.ItemEntity),
		}),
	}
}

func (e OrderItemFulfillmentEntity) GetCreateForm(db *gorm.DB) Group {
	var ord []*OrderEntity
	var res []*ResourceEntity
	db.Table("order").Find(&ord)
	db.Table("resource").Find(&res)
	return Group{
		SelectComponent(ord, "", func(r *OrderEntity) string { return r.Name }, "Выберите заказ, в рамках которого предоставлен товар", "order_id", true, -1),
		SelectComponent(res, "", func(r *ResourceEntity) string { return r.Name }, "Выберите товар", "resource_id", true, -1),
		LabeledInputComponent("number", "", "quantity_fulfilled", "Кол-во предоставлено", "", true),
	}
}

func (e OrderItemFulfillmentEntity) GetReadableName() string {
	return "Предоставление товара в рамках заказа"
}

func (e OrderItemFulfillmentEntity) Validate() bool {
	return true
}

func (e OrderItemFulfillmentEntity) GetId() int {
	return e.ID
}

func (e OrderItemFulfillmentEntity) TableName() string {
	return "order_item_fulfillment"
}

func (e *OrderItemFulfillmentEntity) ValidateAndParseForm(r *http.Request) bool {
	form := r.Form
	if !form.Has("order_id") || !form.Has("item_id") || !form.Has("quantity_fulfilled") {
		return false
	}
	var err error
	e.Order_id, err = strconv.Atoi(form.Get("order_id"))
	if err != nil {
		return false
	}
	e.Item_id, err = strconv.Atoi(form.Get("item_id"))
	if err != nil {
		return false
	}
	e.Quantity_fulfilled, err = strconv.Atoi(form.Get("quantity_fulfilled"))
	if err != nil {
		return false
	}
	return true
}

// func (e *OrderItemFulfillmentEntity) AfterCreate(tx *gorm.DB) (err error) {
// 	e.ResourceEntity.ID = e.Resource_id
// 	res := tx.First(&e.ResourceEntity)
// 	if res.Error != nil {
// 		return res.Error
// 	}
// 	if e.ResourceEntity.Quantity < e.Quantity_spent {
// 		return errors.New("quantity_spent is more then resource quantity")
// 	}
// 	e.ResourceEntity.Quantity -= e.Quantity_spent
// 	res = tx.Updates(&e.ResourceEntity)
// 	if res.Error != nil {
// 		return res.Error
// 	}
// 	return
// }
