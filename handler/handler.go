package handler

import (
	"database/sql"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/sergeykochiev/curs/backend/database/entity"
	"gorm.io/gorm"
)

func EndOrder(w http.ResponseWriter, r *http.Request, db *gorm.DB) {
	id := chi.URLParam(r, "id")
	int_id, err := strconv.Atoi(id)
	if err != nil {
		http.Error(w, "Invalid formdata", http.StatusBadRequest)
		return
	}
	if !r.Form.Has("date_ended") {
		http.Error(w, "Invalid formdata", http.StatusBadRequest)
		return
	}
	if res := db.Updates(&entity.OrderEntity{ID: int_id, Date_ended: sql.NullString{String: r.Form.Get("date_ended"), Valid: true}}); res.Error != nil {
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

func SignupPost(w http.ResponseWriter, r *http.Request, db *gorm.DB) {
	password := r.Form.Get("password")
	name := r.Form.Get("name")
	if password != r.Form.Get("repeat_password") {
		http.Error(w, "Passwords don't match", http.StatusBadRequest)
		return
	}
	var user entity.UserEntity
	res := db.Where("name = ?", name).First(&user)
	if res.Error != nil {
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
