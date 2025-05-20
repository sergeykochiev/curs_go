package entity

import (
	"errors"
	"fmt"
	"html"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	. "github.com/sergeykochiev/curs/backend/gui"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
	. "maragu.dev/gomponents"
	_ "maragu.dev/gomponents/components"
	. "maragu.dev/gomponents/html"
)

type ItemResourceNeed struct {
	Id              decimal.Decimal `gorm:"primaryKey"`
	Resource_id     decimal.Decimal
	Item_id         decimal.Decimal
	Quantity_needed float64
	ResourceEntity  ResourceEntity `gorm:"foreignKey:Resource_id"`
	ItemEntity      ItemEntity     `gorm:"foreignKey:Item_id"`
}

func (e ItemResourceNeed) GetEntityPageButtons() Group {
	return Group{}
}

func (e ItemResourceNeed) GetFilters() Group {
	return Group{
		StringFilterComponent("Название ресурса включает", "resource_name"),
		StringFilterComponent("Название товара включает", "item_name"),
	}
}

func (e *ItemResourceNeed) GetPreloadedDb(db *gorm.DB) *gorm.DB {
	return db.Joins("ResourceEntity").Joins("ItemEntity")
}

func (e *ItemResourceNeed) GetFilteredDb(filters url.Values, db *gorm.DB) *gorm.DB {
	if filters.Has("resource_name") && filters.Get("resource_name") != "" {
		db = db.Joins("ResourceEntity").Where("ResourceEntity__name LIKE ?", "%"+filters.Get("resource_name")+"%")
	}
	if filters.Has("item_name") && filters.Get("item_name") != "" {
		db = db.Joins("ItemEntity").Where("ItemEntity__name LIKE ?", "%"+filters.Get("item_name")+"%")
	}
	return db
}

func (e ItemResourceNeed) GetDataRow() Group {
	return Group{
		TableDataComponent(html.EscapeString(fmt.Sprintf("%d", e.GetId())), Td, fmt.Sprintf("/item_resource_need/%d", e.GetId())),
		TableDataComponent(e.ResourceEntity.Name, Td, ""),
		TableDataComponent(e.ItemEntity.Name, Td, ""),
		TableDataComponent(fmt.Sprintf("%f", e.Quantity_needed), Td, ""),
	}
}

func (e ItemResourceNeed) GetTableHeader() Group {
	return Group{
		TableDataComponent("Id", Th, ""),
		TableDataComponent("Название ресурса", Th, ""),
		TableDataComponent("Наименование товара", Th, ""),
		TableDataComponent("Количество предоставлено (единиц)", Th, ""),
	}
}

func (e ItemResourceNeed) GetEntityPage(recursive bool) Group {
	return Group{
		LabeledFieldComponent("Количество необходимо (единиц)", fmt.Sprintf("%f", e.Quantity_needed)),
		If(recursive, Group{
			RelationCardComponent(fmt.Sprintf("Необходим ресурс #%d (%f %s)", e.ResourceEntity.GetId(), e.Quantity_needed, strings.ToLower(e.ResourceEntity.One_is_called)), &e.ResourceEntity),
			RelationCardComponent(fmt.Sprintf("Необходимо для товара #%d", e.ItemEntity.GetId()), &e.ItemEntity),
		}),
	}
}

func (e ItemResourceNeed) GetCreateForm(db *gorm.DB) Group {
	var ord []*ResourceEntity
	var res []*ItemEntity
	db.Find(&ord)
	db.Find(&res)
	return Group{
		SelectComponent(ord, "", func(r *ResourceEntity) string { return r.Name }, "Выберите ресурс, который необходим для производства товара", "Resource_id", true, -1),
		SelectComponent(res, "", func(r *ItemEntity) string { return r.Name }, "Выберите товар", "item_id", true, -1),
		LabeledInputComponent("number", "", "quantity_needed", "Кол-во необходимо", "", true),
	}
}

func (e ItemResourceNeed) GetReadableName() string {
	return "Необходимость ресурса на товар"
}

func (e ItemResourceNeed) Validate() bool {
	return true
}

func (e ItemResourceNeed) GetId() (int int64) {
	return e.Id.IntPart()
}

func (e *ItemResourceNeed) Clear() {
	*e = ItemResourceNeed{}
}

func (e *ItemResourceNeed) SetId(id int64) {
	e.Id = decimal.NewFromInt(id)
}

func (e ItemResourceNeed) TableName() string {
	return "item_resource_need"
}

func (e *ItemResourceNeed) ValidateAndParseForm(r *http.Request) error {
	form := r.Form
	if !form.Has("Resource_id") || !form.Has("item_id") || !form.Has("quantity_needed") {
		return errors.New("Invalid fields")
	}
	var err error
	e.Resource_id, err = decimal.NewFromString(form.Get("Resource_id"))
	if err != nil {
		return err
	}
	e.Item_id, err = decimal.NewFromString(form.Get("item_id"))
	if err != nil {
		return err
	}
	quantity_needed, err := strconv.ParseFloat(form.Get("quantity_needed"), 32)
	if err != nil {
		return err
	}
	e.Quantity_needed = float64(quantity_needed)
	return nil
}

// func (e *ItemResourceNeed) AfterCreate(tx *gorm.DB) (err error) {
// 	e.ResourceEntity.Id = e.Resource_id
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
