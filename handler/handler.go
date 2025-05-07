package handler

import (
	"fmt"
	"net/http"
	"time"

	billgen_types "github.com/sergeykochiev/billgen/types"
	"github.com/sergeykochiev/curs/backend/database/entity"
	"github.com/sergeykochiev/curs/backend/util"
	"gorm.io/gorm"
)

func EndOrder(w http.ResponseWriter, r *http.Request, db *gorm.DB) {
	ord := r.Context().Value("entity").(*entity.OrderEntity)
	ord.Date_ended.Valid = true
	ord.Date_ended.String = util.GetCurrentTime()
	ord.Ended = 1
	if res := db.Updates(&ord); res.Error != nil {
		http.Error(w, res.Error.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, fmt.Sprintf("/order/%d", ord.ID), http.StatusSeeOther)
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
	cookie.Value = fmt.Sprintf("%d", user.ID)
	http.SetCookie(w, &cookie)
	w.Header().Add("Location", "/")
	w.WriteHeader(http.StatusSeeOther)
}

func GenerateOrderBill(w http.ResponseWriter, r *http.Request, db *gorm.DB, tf billgen_types.GenFunc, main_q *chan func()) {
	ord := r.Context().Value("entity").(*entity.OrderEntity)
	if ord.Ended != 1 {
		http.Error(w, "Order is not ended", http.StatusBadRequest)
		return
	}
	datetime_ended, err := time.Parse(time.DateTime, ord.Date_ended.String)
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
