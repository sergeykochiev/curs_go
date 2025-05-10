package report

import (
	"fmt"

	billgen_types "github.com/sergeykochiev/billgen/types"
	"github.com/sergeykochiev/curs/backend/util"
)

type ItemPopularity struct {
	Name            string
	Last_date       string
	Count_fulfilled int
}

func (ip ItemPopularity) ToTRow() []billgen_types.TDData {
	return []billgen_types.TDData{
		{Value: ip.Name, Align: 1},
		{Value: ip.Last_date, Align: 1},
		{Value: fmt.Sprintf("%d", ip.Count_fulfilled), Align: 1},
	}
}

func (rs ItemPopularity) ToTHead() []billgen_types.THData {
	return []billgen_types.THData{
		{Value: "Название товара/услуги", Width: 0},
		{Value: "Дата последнего предоставления", Width: 0},
		{Value: "Предоставлено (раз)", Width: 0},
	}
}

func (ip ItemPopularity) GetName() string {
	return "Популярность товаров/услуг"
}

func (ip ItemPopularity) GetQuery(is_dl bool, is_dh bool) string {
	return fmt.Sprintf(`select item.name as name, max("order".date_ended) as last_date, count(order_item_fulfillment.id) as count_fulfilled from item left join order_item_fulfillment on item.id = order_item_fulfillment.item_id join "order" on order_item_fulfillment.order_id = "order".id where "order".ended = 1%s%s group by item.id order by count_fulfilled desc`, util.ConditionalArg(is_dl, ` and "order".date_ended > ?`, ""), util.ConditionalArg(is_dh, ` and "order".date_ended < ?`, ""))
}

func (ip *ItemPopularity) Clear() {
	*ip = ItemPopularity{}
}
