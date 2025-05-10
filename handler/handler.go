package handler

import (
	"database/sql"
	"fmt"
	"net/http"
	"time"

	billgen "github.com/sergeykochiev/billgen/gen"
	billgen_types "github.com/sergeykochiev/billgen/types"
	"github.com/sergeykochiev/curs/backend/database/entity"
	"github.com/sergeykochiev/curs/backend/templates"
	"github.com/sergeykochiev/curs/backend/types"
	"github.com/sergeykochiev/curs/backend/util"
	"gorm.io/gorm"
)

func EndOrder(w http.ResponseWriter, r *http.Request, db *gorm.DB) {
	ord := r.Context().Value("entity").(*entity.OrderEntity)
	date_ended := util.GetCurrentDate()
	var ord_res_spe_arr []entity.OrderResourceSpendingEntity
	var ord_res_spe entity.OrderResourceSpendingEntity
	for _, ord_ite_ful := range ord.OrderItemFulfillmentEntities {
		for _, ite_res_nee := range ord_ite_ful.ItemEntity.ItemResourceNeeds {
			// if ite_res_nee.Quantity_needed > ite_res_nee.ResourceEntity.Quantity {
			// 	http.Error(w, fmt.Sprintf(`Failed to end order: not enough resource "%s" in storage`), http.StatusBadRequest)
			// 	return
			// }
			ord_res_spe = entity.OrderResourceSpendingEntity{
				Order_id:       ord.Id,
				Resource_id:    ite_res_nee.Resource_id,
				Quantity_spent: ite_res_nee.Quantity_needed,
				Date:           date_ended,
			}
			ord_res_spe_arr = append(ord_res_spe_arr, ord_res_spe)
		}
	}
	tx := db.Begin()
	if len(ord_res_spe_arr) != 0 {
		if res := tx.Omit("Id").Create(&ord_res_spe_arr); res.Error != nil {
			http.Error(w, res.Error.Error(), http.StatusInternalServerError)
			return
		}
	}
	ord.Date_ended = sql.NullString{
		Valid: true, String: date_ended,
	}
	ord.Ended = true
	if res := tx.Updates(&ord); res.Error != nil {
		http.Error(w, res.Error.Error(), http.StatusInternalServerError)
		return
	}
	tx.Commit()
	http.Redirect(w, r, fmt.Sprintf("/order/%d", ord.GetId()), http.StatusSeeOther)
}

func LoginPost(w http.ResponseWriter, r *http.Request, db *gorm.DB) {
	var user entity.UserEntity
	if res := db.Where("name = ?", r.Form.Get("name")).First(&user); res.Error != nil {
		http.Error(w, res.Error.Error(), http.StatusInternalServerError)
		return
	}
	if !user.CheckPassword(r.Form.Get("password")) {
		http.Error(w, "Wrong password", http.StatusBadRequest)
		return
	}
	var cookie http.Cookie
	cookie.Name = "token"
	cookie.Value = fmt.Sprintf("%d", user.GetId())
	http.SetCookie(w, &cookie)
	w.Header().Add("Location", "/")
	w.WriteHeader(http.StatusSeeOther)
}

func CreateGenerateDatedReportHandler[T interface {
	types.TableTemplater
}](db *gorm.DB, main_q *chan func(), dst T) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()
		is_date_lo := r.Form.Has("date_lo") && r.Form.Get("date_lo") != ""
		is_date_hi := r.Form.Has("date_hi") && r.Form.Get("date_hi") != ""
		var dates []interface{}
		if is_date_lo {
			dates = append(dates, r.Form.Get("date_lo"))
		}
		if is_date_hi {
			dates = append(dates, r.Form.Get("date_hi"))
		}
		var dsta = util.MakeArrayOf(dst)
		if res := db.Raw(dst.GetQuery(is_date_lo, is_date_hi), dates...).Scan(&dsta); res.Error != nil {
			http.Error(w, "Error getting data: "+res.Error.Error(), http.StatusInternalServerError)
			return
		}
		core_heading := fmt.Sprintf("%s%s%s", dst.GetName(), util.ConditionalArg(is_date_lo, " с %s", ""), util.ConditionalArg(is_date_hi, " по %s", ""))
		heading := fmt.Sprintf(core_heading, dates...)
		w.Header().Add("Content-Disposition", fmt.Sprintf(`attachment; filename*=UTF-8''"%s.pdf"`, heading))
		if err := util.RunOnQ(main_q, func() error {
			var tddaa = make([][]billgen_types.TDData, len(dsta))
			for i, dsti := range dsta {
				tddaa[i] = dsti.ToTRow()
			}
			return billgen.CreatePdfFromHtml(templates.TablePage(heading, dst.ToTHead(), tddaa, []billgen_types.TDData{}), w)
		}); err != nil {
			http.Error(w, "Error generating .pdf: "+err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

func GenerateOrderBill(w http.ResponseWriter, r *http.Request, db *gorm.DB, tf billgen_types.GenFunc, main_q *chan func()) {
	ord := r.Context().Value("entity").(*entity.OrderEntity)
	if !ord.Ended {
		http.Error(w, "Order is not ended", http.StatusBadRequest)
		return
	}
	datetime_ended, err := time.Parse(time.DateOnly, ord.Date_ended.String)
	if err != nil {
		http.Error(w, "Failed to parse time", http.StatusBadRequest)
		return
	}
	report_number := "bill_number"
	date_ended := fmt.Sprintf("%d %s %d", datetime_ended.Day(), util.GetRussianMonthGenitive(int(datetime_ended.Month())), datetime_ended.Year())
	w.Header().Add("Content-Disposition", fmt.Sprintf(`attachment; filename*=UTF-8''"Счет № %s %s.pdf"`, report_number, ord.Company_name.String))
	if err = util.RunOnQ(main_q, func() error {
		return tf(w, util.GetCompanyInfoFromEnv(), ord.GetBIL(db), ord.Company_name.String, report_number, date_ended)
	}); err != nil {
		http.Error(w, "Error generating .pdf: "+err.Error(), http.StatusInternalServerError)
		return
	}
}

func SignupPost(w http.ResponseWriter, r *http.Request, db *gorm.DB) {
	password := r.Form.Get("password")
	name := r.Form.Get("name")
	if password != r.Form.Get("repeat_password") {
		http.Error(w, "Passwords don't match", http.StatusBadRequest)
		return
	}
	var user entity.UserEntity
	res := db.Where("name = ?", name).First(&user)
	if res.Error != nil && res.Error != gorm.ErrRecordNotFound {
		http.Error(w, res.Error.Error(), http.StatusInternalServerError)
		return
	}
	if res.RowsAffected > 0 {
		http.Error(w, "User already exists", http.StatusBadRequest)
		return
	}
	user.Name = name
	user.Password = password
	if res := db.Create(&user); res.Error != nil {
		http.Error(w, res.Error.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Add("Location", "/login")
	w.WriteHeader(http.StatusSeeOther)
}
