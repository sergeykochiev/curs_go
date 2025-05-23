package entity

import (
	"errors"
	"fmt"
	"html"
	"net/http"
	"net/url"
	"strconv"

	. "github.com/sergeykochiev/curs/backend/gui"
	"github.com/sergeykochiev/curs/backend/util"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
	. "maragu.dev/gomponents"
	_ "maragu.dev/gomponents/components"
	. "maragu.dev/gomponents/html"
)

type ItemEntity struct {
	Id                           decimal.Decimal `gorm:"primaryKey"`
	Name                         string
	Cost_by_one                  float64
	One_is_called                string
	OrderItemFulfillmentEntities []OrderItemFulfillmentEntity `gorm:"foreignKey:Item_id"`
	ItemResourceNeeds            []ItemResourceNeed           `gorm:"foreignKey:Item_id"`
}

func (e ItemEntity) GetEntityPageButtons() Group {
	return Group{}
}

func (e *ItemEntity) GetFilters() Group {
	return Group{
		StringFilterComponent("Название включает", "name"),
		StringFilterComponent("Единицей является", "one_is_called"),
	}
}

func (e ItemEntity) GetPreloadedDb(db *gorm.DB) *gorm.DB {
	return db.Preload("OrderItemFulfillmentEntities.OrderEntity").Preload("ItemResourceNeeds.ResourceEntity")
}

func (e *ItemEntity) GetFilteredDb(filters url.Values, db *gorm.DB) *gorm.DB {
	if filters.Has("name") && filters.Get("name") != "" {
		db = db.Where("name LIKE ?", "%"+filters.Get("name")+"%")
	}
	if filters.Has("one_is_called") && filters.Get("one_is_called") != "" {
		db = db.Where("one_is_called = ?", filters.Get("one_is_called"))
	}
	return db
}

func (e ItemEntity) GetDataRow() Group {
	return Group{
		TableDataComponent(html.EscapeString(fmt.Sprintf("%d", e.GetId())), Td, fmt.Sprintf("/item/%d", e.GetId())),
		TableDataComponent(e.Name, Td, ""),
		TableDataComponent(fmt.Sprintf("%f", e.Cost_by_one), Td, ""),
		TableDataComponent(e.One_is_called, Td, ""),
	}
}

func (e ItemEntity) GetTableHeader() Group {
	return Group{
		TableDataComponent("Id", Th, ""),
		TableDataComponent("Название", Th, ""),
		TableDataComponent("Цена за единицу", Th, ""),
		TableDataComponent("Единица", Th, ""),
	}
}

func (e *ItemEntity) GetEntityPage(recursive bool) Group {
	return Group{
		LabeledFieldComponent("Название", e.Name),
		LabeledFieldComponent("Цена за единицу", fmt.Sprintf("%f", e.Cost_by_one)),
		LabeledFieldComponent("Единица", e.One_is_called),
		If(recursive, RelationCardArrComponent("Предоставления", e.OrderItemFulfillmentEntities, func(ent OrderItemFulfillmentEntity) Node {
			return RelationCardCoreComponent(util.GetOneReadableName(ent), util.GetOneHref(ent), Group{
				ent.GetEntityPage(false),
				ent.OrderEntity.GetEntityPage(false),
			})
		})),
		If(recursive, RelationCardArrComponent("Необходимы ресурсы", e.ItemResourceNeeds, func(ent ItemResourceNeed) Node {
			return RelationCardCoreComponent(util.GetOneReadableName(ent), util.GetOneHref(ent), Group{
				ent.GetEntityPage(false),
				ent.ResourceEntity.GetEntityPage(false),
			})
		})),
	}
}

func (e ItemEntity) GetCreateForm(db *gorm.DB) Group {
	return Group{
		LabeledInputComponent("text", "", "name", "Название товара", "", true),
		LabeledInputComponent("number", "", "cost_by_one", "Стоимость за единицу", "", true),
		LabeledInputComponent("text", `По умолчанию - "Единица"`, "one_is_called", "Единица названа", "", false),
	}
}

func (e ItemEntity) GetReadableName() string {
	return "Товар/услуга"
}

func (e ItemEntity) GetId() int64 {
	return e.Id.IntPart()
}

func (e *ItemEntity) Clear() {
	*e = ItemEntity{}
}

func (e *ItemEntity) SetId(id int64) {
	e.Id = decimal.NewFromInt(id)
}

func (e ItemEntity) TableName() string {
	return "item"
}

func (e ItemEntity) Validate() bool {
	return true
}

func (e *ItemEntity) ValidateAndParseForm(r *http.Request) error {
	form := r.Form
	if !form.Has("name") || !form.Has("cost_by_one") || !form.Has("one_is_called") {
		return errors.New("Invalid fields")
	}
	e.Name = form.Get("name")
	e.One_is_called = form.Get("one_is_called")
	cost_by_one, err := strconv.Atoi(form.Get("cost_by_one"))
	if err != nil {
		return err
	}
	e.Cost_by_one = float64(cost_by_one)
	return nil
}
