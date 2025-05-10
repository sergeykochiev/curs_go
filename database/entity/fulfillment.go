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

type OrderItemFulfillmentEntity struct {
	Id                 decimal.Decimal `gorm:"primaryKey;serializer:decimal"`
	Order_id           decimal.Decimal `gorm:"serializer:decimal"`
	Item_id            decimal.Decimal `gorm:"serializer:decimal"`
	Quantity_fulfilled float32
	OrderEntity        OrderEntity `gorm:"foreignKey:Order_id"`
	ItemEntity         ItemEntity  `gorm:"foreignKey:Item_id"`
}

func (e OrderItemFulfillmentEntity) GetEntityPageButtons() Group {
	return Group{}
}

func (e OrderItemFulfillmentEntity) GetFilters() Group {
	return Group{
		StringFilterComponent("Название заказа включает", "order_name"),
		StringFilterComponent("Название товара включает", "item_name"),
	}
}

func (e *OrderItemFulfillmentEntity) GetPreloadedDb(db *gorm.DB) *gorm.DB {
	return db.Joins("OrderEntity").Joins("ItemEntity")
}

func (e *OrderItemFulfillmentEntity) GetFilteredDb(filters url.Values, db *gorm.DB) *gorm.DB {
	if filters.Has("order_name") && filters.Get("order_name") != "" {
		db = db.Joins("OrderEntity").Where("OrderEntity__name LIKE ?", "%"+filters.Get("order_name")+"%")
	}
	if filters.Has("item_name") && filters.Get("item_name") != "" {
		db = db.Joins("ItemEntity").Where("ItemEntity__name LIKE ?", "%"+filters.Get("item_name")+"%")
	}
	return db
}

func (e OrderItemFulfillmentEntity) GetDataRow() Group {
	return Group{
		TableDataComponent(html.EscapeString(fmt.Sprintf("%d", e.GetId())), Td, fmt.Sprintf("/order_item_fulfillment/%d", e.GetId())),
		TableDataComponent(e.OrderEntity.Name, Td, ""),
		TableDataComponent(e.ItemEntity.Name, Td, ""),
		TableDataComponent(fmt.Sprintf("%f", e.Quantity_fulfilled), Td, ""),
	}
}

func (e OrderItemFulfillmentEntity) GetTableHeader() Group {
	return Group{
		TableDataComponent("Id", Th, ""),
		TableDataComponent("Название заказа", Th, ""),
		TableDataComponent("Наименование товара", Th, ""),
		TableDataComponent("Количество предоставлено (единиц)", Th, ""),
	}
}

func (e OrderItemFulfillmentEntity) GetEntityPage(recursive bool) Group {
	return Group{
		LabeledFieldComponent("Количество предоставлено (единиц)", fmt.Sprintf("%f", e.Quantity_fulfilled)),
		If(recursive, Group{
			RelationCardComponent(fmt.Sprintf("Предоставлено в рамках заказа #%d", e.OrderEntity.GetId()), &e.OrderEntity),
			RelationCardComponent(fmt.Sprintf("Предоставлен товар #%d (%f шт.)", e.ItemEntity.GetId(), e.Quantity_fulfilled), &e.ItemEntity),
		}),
	}
}

func (e OrderItemFulfillmentEntity) GetCreateForm(db *gorm.DB) Group {
	var ord []*OrderEntity
	var res []*ItemEntity
	db.Find(&ord)
	db.Find(&res)
	return Group{
		SelectComponent(ord, "", func(r *OrderEntity) string { return r.Name }, "Выберите заказ, в рамках которого предоставлен товар", "order_id", true, -1),
		SelectComponent(res, "", func(r *ItemEntity) string { return r.Name }, "Выберите товар", "item_id", true, -1),
		LabeledInputComponent("number", "", "quantity_fulfilled", "Кол-во предоставлено", "", true),
	}
}

func (e OrderItemFulfillmentEntity) GetReadableName() string {
	return "Предоставление товара в рамках заказа"
}

func (e OrderItemFulfillmentEntity) Validate() bool {
	return true
}

func (e OrderItemFulfillmentEntity) GetId() int64 {
	return e.Id.IntPart()
}

func (e *OrderItemFulfillmentEntity) Clear() {
	*e = OrderItemFulfillmentEntity{}
}

func (e *OrderItemFulfillmentEntity) SetId(id int64) {
	e.Id = decimal.NewFromInt(id)
}

func (e OrderItemFulfillmentEntity) TableName() string {
	return "order_item_fulfillment"
}

func (e *OrderItemFulfillmentEntity) ValidateAndParseForm(r *http.Request) error {
	form := r.Form
	if !form.Has("order_id") || !form.Has("item_id") || !form.Has("quantity_fulfilled") {
		return errors.New("Invalid fields")
	}
	var err error
	e.Order_id, err = decimal.NewFromString(form.Get("order_id"))
	if err != nil {
		return err
	}
	e.Item_id, err = decimal.NewFromString(form.Get("item_id"))
	if err != nil {
		return err
	}
	quantity_fulfilled, err := strconv.ParseFloat(form.Get("quantity_fulfilled"), 32)
	if err != nil {
		return err
	}
	e.Quantity_fulfilled = float32(quantity_fulfilled)
	return nil
}

// func (e *OrderItemFulfillmentEntity) AfterCreate(tx *gorm.DB) (err error) {
// 	e.ResourceEntity.GetId() = e.Resource_id
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
