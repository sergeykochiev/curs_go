package handler

import (
	"database/sql"
	"fmt"
	"net/http"

	billgen "github.com/sergeykochiev/billgen"
	"github.com/sergeykochiev/curs/backend/database/entity"
	"github.com/sergeykochiev/curs/backend/util"
	"gorm.io/gorm"
)

func EndOrder(w http.ResponseWriter, r *http.Request, db *gorm.DB) {
	order := r.Context().Value("entity").(entity.OrderEntity)
	if !r.Form.Has("date_ended") {
		http.Error(w, "Invalid formdata", http.StatusBadRequest)
		return
	}
	order.Date_ended = sql.NullString{String: r.Form.Get("date_ended"), Valid: true}
	order.Ended = 1
	if res := db.Updates(&order); res.Error != nil {
		http.Error(w, res.Error.Error(), http.StatusInternalServerError)
		return
	}
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

func GenerateOrderBill(w http.ResponseWriter, r *http.Request, db *gorm.DB) {
	date := r.Form.Get("date")
	client_company := r.Form.Get("client_company")
	ci := util.GetCompanyInfoFromEnv()
	bil := r.Context().Value("entity").(entity.OrderEntity).GetBIL()
	billgen.CreateBillPdf(ci, bil, client_company, date, "")
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
