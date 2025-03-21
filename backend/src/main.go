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

	_ "github.com/lib/pq"
	_ "modernc.org/sqlite"

	. "github.com/sergeykochiev/curs/backend/entity"
	. "github.com/sergeykochiev/curs/backend/gui"
	. "github.com/sergeykochiev/curs/backend/types"
	. "github.com/sergeykochiev/curs/backend/util"
)

var Tstate TStateType

const addr = "localhost:3003"

func returnEntityListPage[T interface {
	HtmlTemplater
	ActiveRecorder
	Identifier
}](db *sql.DB, ent T) (Node, error) {
	arr, err := GetRows(db, ent, "")
	if err != nil {
		println("Failed to get: ", err.Error())
		return nil, err
	}
	return DataPageComponent(ent, arr, Tstate.DB), nil
}

func handleSignup(db *sql.DB, name string, password string, repeat string) error {
	if password != repeat {
		return errors.New("Passwords don't match")
	}
	row := db.QueryRow("select * from public.user where name = $1", name)
	var ent UserEntity
	err := ent.ScanRow(row)
	if !errors.Is(err, sql.ErrNoRows) {
		return err
	}
	_, err = db.Exec("insert into public.user (name, password, is_admin) values ($1, $2, false)", name, password)
	return err
}

func handleLogin(db *sql.DB, name string, password string) (int, error) {
	var ent UserEntity
	err := ent.ScanRow(db.QueryRow("select * from public.user where name = $1", name))
	if err != nil {
		return -1, err
	}
	if !ent.CheckPassword(password) {
		return ent.Id, errors.New("Wrong password")
	}
	return ent.Id, nil
}

func checkIfAuthorized(r *http.Request) bool {
	return len(r.Cookies()) > 0 && len(r.CookiesNamed("token")) > 0
}

type TypographyHandler struct{}

func (TypographyHandler) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	paths := strings.Split(r.RequestURI, "/")
	checkAuth := func() bool {
		if !checkIfAuthorized(r) {
			rw.Header().Add("Location", "/login")
			rw.WriteHeader(http.StatusSeeOther)
			return false
		}
		return true
	}
	switch r.Method {
	case http.MethodGet:
		println("Received get request")
		var node Node
		var err error
		switch paths[1] {
		case "":
			if !checkAuth() {
				return
			}
			node = MainPageComponent()
		case "signup":
			node = UserFormComponent(true)
		case "login":
			node = UserFormComponent(false)
		case "resources":
			if !checkAuth() {
				return
			}
			var ent ResourceEntity
			node, err = returnEntityListPage(Tstate.DB, &ent)
		case "resource_spendings":
			if !checkAuth() {
				return
			}
			var ent ResourceSpendingEntity
			node, err = returnEntityListPage(Tstate.DB, &ent)
		case "end_order":
			if !checkAuth() {
				return
			}
			var ent OrderEntity
			var arr []*OrderEntity
			arr, err = GetRows(Tstate.DB, &ent, "where ended = false")
			node = EndOrderComponent(arr)
		case "resource_resupplies":
			if !checkAuth() {
				return
			}
			var ent ResourceResupplyEntity
			node, err = returnEntityListPage(Tstate.DB, &ent)
		case "orders":
			if !checkAuth() {
				return
			}
			var ent OrderEntity
			node, err = returnEntityListPage(Tstate.DB, &ent)
		case "create_order":
			if !checkAuth() {
				return
			}
			node = CreateOrderFormComponent()
		default:
			rw.WriteHeader(http.StatusNotFound)
			return
		}
		if err != nil {
			// TODO handle this
			rw.WriteHeader(http.StatusInternalServerError)
			return
		}
		node.Render(rw)
	case http.MethodPost:
		println("Received post request")
		if len(paths) > 2 {
			rw.WriteHeader(http.StatusNotFound)
			return
		}
		r.ParseForm()
		switch paths[1] {
		case "signup":
			err := handleSignup(Tstate.DB, r.Form.Get("name"), r.Form.Get("password"), r.Form.Get("repeat_password"))
			if err != nil {
				// TODO handle error
				println("ERROR error signing up:", err.Error())
				rw.WriteHeader(http.StatusInternalServerError)
				return
			}
			rw.Header().Add("Location", "/login")
			rw.WriteHeader(http.StatusSeeOther)
		case "login":
			id, err := handleLogin(Tstate.DB, r.Form.Get("name"), r.Form.Get("password"))
			if err != nil {
				// TODO handle error
				println("ERROR error logging in:", err.Error())
				rw.WriteHeader(http.StatusInternalServerError)
				return
			}
			var cookie http.Cookie
			cookie.Name = "token"
			cookie.Value = fmt.Sprintf("%d", id)
			http.SetCookie(rw, &cookie)
			rw.Header().Add("Location", "/")
			rw.WriteHeader(http.StatusSeeOther)
		case "create_order":
			if !checkAuth() {
				return
			}
			if !r.Form.Has("name") || !r.Form.Has("client_name") || !r.Form.Has("client_phone") || !r.Form.Has("date_created") {
				rw.WriteHeader(http.StatusBadRequest)
				return
			}
			var ent OrderEntity
			ent.Client_name = r.Form.Get("client_name")
			ent.Client_phone = r.Form.Get("client_phone")
			ent.Name = r.Form.Get("name")
			ent.Date_created = r.Form.Get("date_created")
			var err error
			ent.Creator_id, err = strconv.Atoi(r.CookiesNamed("token")[0].Value)
			if !ent.Validate() {
				rw.WriteHeader(http.StatusUnprocessableEntity)
				return
			}
			if _, err = ent.Insert(Tstate.DB); err != nil {
				println(err.Error())
				rw.WriteHeader(http.StatusInternalServerError)
				return
			}
		case "end_order":
			if !r.Form.Has("id") || !r.Form.Has("date_ended") {
				rw.WriteHeader(http.StatusBadRequest)
				return
			}
			id, err := strconv.Atoi(r.Form.Get("id"))
			if err != nil {
				rw.WriteHeader(http.StatusInternalServerError)
				return
			}
			if _, err = Tstate.DB.Exec("update public.order set ended = true, date_ended = $1 where id = $2", r.Form.Get("date_ended"), id); err != nil {
				println(err.Error())
				rw.WriteHeader(http.StatusInternalServerError)
				return
			}
		default:
			rw.WriteHeader(http.StatusNotFound)
		}
	default:
		rw.WriteHeader(http.StatusNotFound)
	}
}

func main() {
	conStr := "dbname=test user=postgres password=postgres host=localhost port=5432 sslmode=disable"
	if len(os.Args) <= 1 {
		log.Fatal("unsufficient cmd args count")
	}
	var err error
	Tstate.DB, err = sql.Open("postgres", conStr)
	if err != nil {
		log.Fatal("cannot connect to db: ", err.Error())
	}
	if err = Tstate.DB.Ping(); err != nil {
		log.Fatal("cannot connect to db: ", err.Error())
	}
	switch os.Args[1] {
	case "initdb":
		{
			data, err := os.ReadFile("schema.sql")
			if err != nil {
				log.Fatal("cannot read schema file (./schema.sql): ", err.Error())
			}
			_, err = Tstate.DB.Exec(string(data))
			if err != nil {
				println("error initializing db: ", err.Error())
			}
			os.Exit(0)
		}
	case "http":
		{
			var Th TypographyHandler
			fmt.Printf("INFO listening on http://%s\n", addr)
			http.ListenAndServe(addr, Th)
		}
	default:
		log.Fatal("unknown 1st cmd arg")
	}
}
