package main

import (
	"fmt"
	"log"
	"net/http"
	"reflect"

	"github.com/go-chi/chi/v5"
	"github.com/sergeykochiev/curs/backend/database"
	. "github.com/sergeykochiev/curs/backend/database/entity"
	"github.com/sergeykochiev/curs/backend/gui"
	"github.com/sergeykochiev/curs/backend/handler"
	"github.com/sergeykochiev/curs/backend/middleware"
	"github.com/sergeykochiev/curs/backend/types"
	"gorm.io/gorm"
)

const addr = "localhost:3003"

func EntityRouterFactory[T interface {
	types.HtmlTemplater
	types.Identifier
	types.FormParser
	types.Filterator
}](db *gorm.DB, entity T) func(r chi.Router) {
	create := func(w http.ResponseWriter, r *http.Request) {
		res := db.Create(entity)
		if res.Error != nil {
			http.Error(w, res.Error.Error(), http.StatusInternalServerError)
			return
		}
	}
	// update := func(w http.ResponseWriter, r *http.Request) {
	// 	res := db.Updates(entity)
	// 	if res.Error != nil {
	// 		http.Error(w, res.Error.Error(), http.StatusInternalServerError)
	// 		return
	// 	}
	// }
	// delete := func(w http.ResponseWriter, r *http.Request) {
	// 	res := db.Delete(entity)
	// 	if res.Error != nil {
	// 		http.Error(w, res.Error.Error(), http.StatusInternalServerError)
	// 		return
	// 	}
	// }
	getAllPage := func(w http.ResponseWriter, r *http.Request) {
		filteredDb := entity.GetFilteredDb(r.URL.Query(), db)
		arr := reflect.MakeSlice(reflect.SliceOf(reflect.TypeOf(entity)), 0, 0).Interface()
		res := filteredDb.Find(&arr)
		if res.Error != nil {
			http.Error(w, res.Error.Error(), http.StatusInternalServerError)
			return
		}
		gui.EntityListPage(entity, arr.([]T)).Render(w)
	}
	getOnePage := func(w http.ResponseWriter, r *http.Request) { gui.EntityPage(entity).Render(w) }
	getCreatePage := func(w http.ResponseWriter, r *http.Request) {
		gui.CreateFormComponent(entity.GetReadableName(), entity.GetCreateForm(db)).Render(w)
	}
	return func(r chi.Router) {
		r.Use(middleware.WithAuthUserIdContext)
		r.Get("/", getAllPage)
		r.Get("/create", getCreatePage)
		r.Route("/", func(r chi.Router) {
			r.Use(middleware.WithFormEntityContextFactory(entity))
			r.Use(middleware.WithEntityValidation)
			r.Post("/", create)
		})
		r.Route("/{id}", func(r chi.Router) {
			r.Use(middleware.WithDbEntityContextFactory(entity, db))
			r.Get("/", getOnePage)
			// r.Delete("/", delete)
		})
		// r.Route("/{id}", func(r chi.Router) {
		// 	r.Use(WithFormEntityContextFactory(entity))
		// 	r.Use(WithEntityValidation)
		// 	r.Patch("/", update)
		// })
	}
}

func main() {
	var err error
	var main_db, db *gorm.DB
	main_db, err = database.ConnectDb("main.db")
	if err != nil {
		log.Fatal("F cannot connect to main db: ", err.Error())
	}
	var dbs []DatabaseEntity
	res := main_db.Table("databases").Find(&dbs)
	if res.Error != nil {
		log.Fatal("F cannot query databases list: ", res.Error.Error())
	}
	if len(dbs) == 0 {
		db = database.CreateNewDataDb(main_db)
	} else {
		db = database.GetOldDataDb(main_db, dbs[0])
	}
	r := chi.NewRouter()
	r.Use(middleware.WithRequestInfoLogging)
	r.Route("/", func(r chi.Router) {
		r.Use(middleware.WithAuthUserIdContext)
		r.Get("/", func(w http.ResponseWriter, r *http.Request) { gui.MainPageComponent().Render(w) })
	})
	r.Route("/signup", func(r chi.Router) {
		r.Use(middleware.WithFormFieldsValidationFactory([]string{"name", "password", "repeat_password"}))
		r.Get("/", func(w http.ResponseWriter, r *http.Request) { gui.UserFormComponent(true).Render(w) })
		r.Post("/", func(w http.ResponseWriter, r *http.Request) { handler.SignupPost(w, r, db) })
	})
	r.Route("/login", func(r chi.Router) {
		r.Use(middleware.WithFormFieldsValidationFactory([]string{"name", "password"}))
		r.Get("/", func(w http.ResponseWriter, r *http.Request) { gui.UserFormComponent(false).Render(w) })
		r.Post("/", func(w http.ResponseWriter, r *http.Request) { handler.LoginPost(w, r, db) })
	})
	r.Route("/order", EntityRouterFactory(db, &OrderEntity{}))
	r.Route("/order/{id}/end", func(r chi.Router) {
		r.Use(middleware.WithAuthUserIdContext)
		r.Post("/", func(w http.ResponseWriter, r *http.Request) { handler.EndOrder(w, r, db) })
	})
	r.Route("/resource", EntityRouterFactory(db, &ResourceEntity{}))
	r.Route("/resource_resupply", EntityRouterFactory(db, &ResourceResupplyEntity{}))
	r.Route("/resource_spending", EntityRouterFactory(db, &ResourceSpendingEntity{}))
	fmt.Printf("I Listening on http://%s\n", addr)
	http.ListenAndServe(addr, r)
}
