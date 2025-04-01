package main

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	. "maragu.dev/gomponents"
	_ "maragu.dev/gomponents/components"
	_ "modernc.org/sqlite"

	"github.com/google/uuid"
	. "github.com/sergeykochiev/curs/backend/entity"
	. "github.com/sergeykochiev/curs/backend/gui"
	. "github.com/sergeykochiev/curs/backend/types"
	. "github.com/sergeykochiev/curs/backend/util"
)

var State StateType

const addr = "localhost:3003"

func checkAuth(rw http.ResponseWriter, r *http.Request) bool {
	if len(r.Cookies()) > 0 && len(r.CookiesNamed("token")) > 0 {
		return true
	}
	rw.Header().Add("Location", "/login")
	return false
}

func GetEntityPage[T HtmlEntity](rw http.ResponseWriter, r *http.Request, paths []string, ent T, db QueryExecutor) (Node, int, error) {
	if !checkAuth(rw, r) {
		return nil, http.StatusSeeOther, errors.New("unauthorized")
	}
	if len(paths) > 3 {
		return nil, http.StatusNotFound, nil
	}
	if len(paths) == 3 {
		node, err := ReturnEntityPage(db, ent, paths[2])
		return node, 0, err
	}
	node, err := ReturnEntityListPage(db, ent)
	return node, 0, err
}

type TypographyHandler struct{}

func GetHandle(rw http.ResponseWriter, r *http.Request) (Node, int) {
	paths := strings.Split(r.RequestURI, "/")
	println("I Received GET")
	var node Node
	var err error
	code := http.StatusOK
	switch paths[1] {
	case "":
		if !checkAuth(rw, r) {
			return nil, http.StatusSeeOther
		}
		node = MainPageComponent()
	case "signup":
		node = UserFormComponent(true)
	case "login":
		node = UserFormComponent(false)
	case "resource":
		var ent ResourceEntity
		node, code, err = GetEntityPage(rw, r, paths, &ent, State.DB)
	case "resource_spending":
		var ent ResourceSpendingEntity
		node, code, err = GetEntityPage(rw, r, paths, &ent, State.DB)
	case "resource_resupply":
		var ent ResourceResupplyEntity
		node, code, err = GetEntityPage(rw, r, paths, &ent, State.DB)
	case "order":
		var ent OrderEntity
		node, code, err = GetEntityPage(rw, r, paths, &ent, State.DB)
	case "create_order":
		var ent OrderEntity
		node = CreateFormComponent(ent.GetReadableName(), ent.GetCreateForm())
	case "create_resupply":
		if !checkAuth(rw, r) {
			return nil, http.StatusSeeOther
		}
		var ent ResourceResupplyEntity
		var res ResourceEntity
		resarr, err := GetRows(State.DB, &res, "")
		if err != nil {
			println("ERROR error fetching resources for resupply form:", err.Error())
			rw.WriteHeader(http.StatusInternalServerError)
			return nil, http.StatusSeeOther
		}
		node = CreateFormComponent(ent.GetReadableName(), ent.GetCreateForm(resarr))
	case "create_spending":
		if !checkAuth(rw, r) {
			return nil, http.StatusSeeOther
		}
		var ent ResourceSpendingEntity
		var res ResourceEntity
		var resarr []*ResourceEntity
		resarr, err = GetRows(State.DB, &res, "")
		if err != nil {
			// TODO handle error
			println("ERROR error fetching resources for spending form: ", err.Error())
			return nil, http.StatusInternalServerError
		}
		var ord OrderEntity
		var ordarr []*OrderEntity
		ordarr, err = GetRows(State.DB, &ord, "")
		if err != nil {
			// TODO handle error
			println("ERROR error fetching orders for spending form: ", err.Error())
			return nil, http.StatusInternalServerError
		}
		node = CreateFormComponent(ent.GetReadableName(), ent.GetCreateForm(ordarr, resarr))
	case "create_resource":
		if !checkAuth(rw, r) {
			return nil, http.StatusSeeOther
		}
		var ent ResourceEntity
		node = CreateFormComponent(ent.GetReadableName(), ent.GetCreateForm())
	default:
		code = http.StatusNotFound
	}
	if err != nil {
		println("E [ GET", paths[1], "]", err.Error())
		if errors.Is(err, sql.ErrNoRows) {
			node = NotFoundPage()
		} else {
			return nil, http.StatusInternalServerError
		}
	} else {
		println("I [ GET", paths[1], "] successful")
	}
	return node, code
}

func CreateResourceResupply(rw http.ResponseWriter, r *http.Request) (int, error) {
	var err error
	if !checkAuth(rw, r) {
		return http.StatusSeeOther, errors.New("unauthorized")
	}
	if !r.Form.Has("resource_id") || !r.Form.Has("quantity_added") || !r.Form.Has("date") {
		return http.StatusBadRequest, errors.New("invalid formdata")
	}
	var ent ResourceResupplyEntity
	ent.Resource_id, err = strconv.Atoi(r.Form.Get("resource_id"))
	if err != nil {
		return http.StatusInternalServerError, PreappendError("failed to process resource_id field", err)
	}
	ent.Quantity_added, err = strconv.Atoi(r.Form.Get("quantity_added"))
	if err != nil {
		return http.StatusInternalServerError, PreappendError("failed to process quantity_added field", err)
	}
	ent.Date = r.Form.Get("date")
	if !ent.Validate() {
		return http.StatusUnprocessableEntity, PreappendError("unvalidated", err)
	}
	tx, err := State.DB.Begin()
	if err != nil {
		return http.StatusInternalServerError, PreappendError("transaction start failed", err)
	}
	if _, err = ent.Insert(tx); err != nil {
		tx.Rollback()
		return http.StatusInternalServerError, PreappendError("inserting entity failed", err)
	}
	var resource ResourceEntity
	if err = GetSingleRow(tx, &resource, ent.Resource_id); err != nil {
		tx.Rollback()
		return http.StatusInternalServerError, PreappendError("getting resource by id failed", err)
	}
	resource.Quantity += ent.Quantity_added
	if _, err = resource.Update(tx); err != nil {
		tx.Rollback()
		return http.StatusInternalServerError, PreappendError("updating resource by id failed", err)
	}
	tx.Commit()
	return http.StatusCreated, nil
}

func CreateOrder(rw http.ResponseWriter, r *http.Request) (int, error) {
	var err error
	if !checkAuth(rw, r) {
		return http.StatusSeeOther, errors.New("unauthorized")
	}
	if !r.Form.Has("name") || !r.Form.Has("client_name") || !r.Form.Has("client_phone") || !r.Form.Has("date_created") {
		return http.StatusBadRequest, errors.New("invalid formdata")
	}
	var ent OrderEntity
	ent.Client_name = r.Form.Get("client_name")
	ent.Client_phone = r.Form.Get("client_phone")
	ent.Name = r.Form.Get("name")
	ent.Date_created = r.Form.Get("date_created")
	ent.Creator_id, err = strconv.Atoi(r.CookiesNamed("token")[0].Value)
	if !ent.Validate() {
		return http.StatusUnprocessableEntity, PreappendError("unvalidated", err)
	}
	if _, err = ent.Insert(State.DB); err != nil {
		return http.StatusInternalServerError, PreappendError("inserting entity failed", err)
	}
	return http.StatusCreated, nil
}

func CreateResource(rw http.ResponseWriter, r *http.Request) (int, error) {
	if !checkAuth(rw, r) {
		return http.StatusSeeOther, errors.New("unauthorized")
	}
	if !r.Form.Has("name") || !r.Form.Has("cost_by_one") {
		return http.StatusBadRequest, errors.New("invalid formdata")
	}
	var ent ResourceEntity
	ent.Name = r.Form.Get("name")
	cost_by_one, err := strconv.Atoi(r.Form.Get("cost_by_one"))
	if err != nil {
		return http.StatusInternalServerError, PreappendError("failed to process client_phone field", err)
	}
	ent.Cost_by_one = float32(cost_by_one)
	if !ent.Validate() {
		return http.StatusUnprocessableEntity, PreappendError("unvalidated", err)
	}
	if _, err = ent.Insert(State.DB); err != nil {
		return http.StatusInternalServerError, PreappendError("inserting entity failed", err)
	}
	return http.StatusCreated, nil
}

func CreateResourceSpending(rw http.ResponseWriter, r *http.Request) (int, error) {
	var err error
	if !checkAuth(rw, r) {
		return http.StatusSeeOther, errors.New("unauthorized")
	}
	if !r.Form.Has("order_id") || !r.Form.Has("resource_id") || !r.Form.Has("quantity_spent") || !r.Form.Has("date") {
		return http.StatusBadRequest, errors.New("invalid formdata")
	}
	var ent ResourceSpendingEntity
	ent.Order_id, err = strconv.Atoi(r.Form.Get("order_id"))
	if err != nil {
		return http.StatusInternalServerError, PreappendError("failed to process order_id field", err)
	}
	ent.Resource_id, err = strconv.Atoi(r.Form.Get("resource_id"))
	if err != nil {
		return http.StatusInternalServerError, PreappendError("failed to process resource_id field", err)
	}
	ent.Quantity_spent, err = strconv.Atoi(r.Form.Get("quantity_spent"))
	if err != nil {
		return http.StatusInternalServerError, PreappendError("failed to process quantity_spent field", err)
	}
	ent.Date = r.Form.Get("date")
	if !ent.Validate() {
		return http.StatusUnprocessableEntity, PreappendError("unvalidated", err)
	}
	tx, err := State.DB.Begin()
	if err != nil {
		return http.StatusInternalServerError, PreappendError("transaction start failed", err)
	}
	if _, err = ent.Insert(tx); err != nil {
		tx.Rollback()
		return http.StatusInternalServerError, PreappendError("inserting entity failed", err)
	}
	var resource ResourceEntity
	if err = GetSingleRow(tx, &resource, ent.Resource_id); err != nil {
		tx.Rollback()
		return http.StatusInternalServerError, PreappendError("getting resource by id failed", err)
	}
	if resource.Quantity < ent.Quantity_spent {
		tx.Rollback()
		return http.StatusBadRequest, errors.New("quantity_spent is more then resource quantity")
	}
	resource.Quantity -= ent.Quantity_spent
	if _, err = resource.Update(tx); err != nil {
		tx.Rollback()
		return http.StatusInternalServerError, PreappendError("updating resource by id failed", err)
	}
	tx.Commit()
	return http.StatusCreated, nil
}

func EndOrder(rw http.ResponseWriter, r *http.Request, id string) (int, error) {
	if !checkAuth(rw, r) {
		return http.StatusSeeOther, errors.New("unauthorized")
	}
	order_id, err := strconv.Atoi(id)
	if err != nil {
		return http.StatusInternalServerError, PreappendError("failed to process id", err)
	}
	if !r.Form.Has("date_ended") {
		return http.StatusBadRequest, errors.New("invalid formdata")
	}
	date_ended := r.Form.Get("date_ended")
	if _, err = State.DB.Exec("update \"order\" set ended = true, date_ended = $1 where id = $2", date_ended, order_id); err != nil {
		return http.StatusInternalServerError, PreappendError("failed to update entity", err)
	}
	return http.StatusOK, nil
}

func Signup(rw *http.ResponseWriter, r *http.Request) (int, error) {
	if !r.Form.Has("password") || !r.Form.Has("repeat_password") || !r.Form.Has("name") {
		return http.StatusBadRequest, errors.New("invalid formdata")
	}
	password := r.Form.Get("password")
	name := r.Form.Get("name")
	if password != r.Form.Get("repeat_password") {
		return http.StatusBadRequest, errors.New("passwords don't match")
	}
	row := State.DB.QueryRow("select * from user where name = $1", name)
	var ent UserEntity
	if !errors.Is(ent.ScanRow(row), sql.ErrNoRows) {
		return http.StatusBadRequest, errors.New("user already exists")
	}
	if _, err := State.DB.Exec("insert into user (name, password, is_admin) values ($1, $2, false)", name, password); err != nil {
		return http.StatusInternalServerError, err
	}
	(*rw).Header().Add("Location", "/login")
	return http.StatusSeeOther, nil
}

func Login(rw *http.ResponseWriter, r *http.Request) (int, error) {
	if !r.Form.Has("password") || !r.Form.Has("name") {
		return http.StatusBadRequest, errors.New("invalid formdata")
	}
	var ent UserEntity
	if err := ent.ScanRow(State.DB.QueryRow("select * from user where name = $1", r.Form.Get("name"))); err != nil {
		return http.StatusInternalServerError, err
	}
	if !ent.CheckPassword(r.Form.Get("password")) {
		return http.StatusBadRequest, errors.New("wrong password")
	}
	var cookie http.Cookie
	cookie.Name = "token"
	cookie.Value = fmt.Sprintf("%d", ent.Id)
	http.SetCookie(*rw, &cookie)
	(*rw).Header().Add("Location", "/")
	return http.StatusSeeOther, nil
}

func PostHandle(rw http.ResponseWriter, r *http.Request) int {
	paths := strings.Split(r.RequestURI, "/")
	println("I Received POST")
	code := http.StatusNotFound
	var err error
	r.ParseForm()
	switch paths[1] {
	case "signup":
		code, err = Signup(&rw, r)
	case "login":
		code, err = Login(&rw, r)
	case "create_order":
		code, err = CreateOrder(rw, r)
	case "create_resource":
		code, err = CreateResource(rw, r)
	case "create_resupply":
		code, err = CreateResourceResupply(rw, r)
	case "create_spending":
		code, err = CreateResourceSpending(rw, r)
	case "end_order":
		if len(paths) != 3 {
			return http.StatusNotFound
		}
		code, err = EndOrder(rw, r, paths[2])
	}
	if err != nil {
		println("E [ POST", paths[1], "]", err.Error())
	} else {
		println("I [ POST", paths[1], "] successful")
	}
	return code
}

func (TypographyHandler) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	code := http.StatusNotFound
	switch r.Method {
	case http.MethodGet:
		var node Node
		node, code = GetHandle(rw, r)
		if node != nil {
			node.Render(rw)
		}
	case http.MethodPost:
		code = PostHandle(rw, r)
	}
	if code != 0 {
		rw.WriteHeader(code)
	}
}

func initdb(db *sql.DB) error {
	data, err := os.ReadFile("schema.sql")
	if err != nil {
		log.Fatal("E cannot read schema file (./schema.sql): ", err.Error())
	}
	_, err = db.Exec(string(data))
	if err != nil {
		println("E initializing db failed: ", err.Error())
	}
	return err
}

func main() {
	var err error
	State.MainDB, err = sql.Open("sqlite", "main.db")
	if err != nil || State.MainDB.Ping() != nil {
		log.Fatal("F cannot connect to main db: ", err.Error())
	}
	databases, err := State.MainDB.Query("select * from databases")
	if err != nil {
		log.Fatal("F cannot query databases list: ", err.Error())
	}
	var d DatabaseRecord
	if databases.Next() {
		if err := databases.Scan(&d.Id, &d.Name, &d.Filepath, &d.Is_initialized); err != nil {
			log.Fatal("F cannot scan database info: ", err.Error())
		}
		State.DB, err = sql.Open("sqlite", d.Filepath)
		if err != nil || State.MainDB.Ping() != nil {
			log.Fatal("F Cannot connect to data db: ", err.Error())
		}
		if d.Is_initialized == 0 {
			if err = initdb(State.DB); err != nil {
				log.Fatal("F Cannot initialize uninitialized data db: ", err.Error())
			}
			State.MainDB.Exec("update databases set is_initialized = 1 where id = $1", d.Id)
		}
	} else {
		dbuuid := uuid.New().String()
		name := "autogenerated_" + dbuuid
		State.MainDB.Exec("insert into databases (name, filepath, is_initialized) values ($1, $2, 0)", name, dbuuid+".db")
		State.DB, err = sql.Open("sqlite", dbuuid+".db")
		if err != nil || State.MainDB.Ping() != nil {
			log.Fatal("F connect to autogenerated data db: ", err.Error())
		}
		if err = initdb(State.DB); err != nil {
			log.Fatal("F initialize uninitialized data db: ", err.Error())
		}
	}

	var Th TypographyHandler
	fmt.Printf("I Listening on http://%s\n", addr)
	http.ListenAndServe(addr, Th)
}
