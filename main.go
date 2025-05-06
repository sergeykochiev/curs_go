package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"reflect"

	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"
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
}](db *gorm.DB, entity T, id_route func(r chi.Router)) func(r chi.Router) {
	create := func(w http.ResponseWriter, r *http.Request) {
		res := db.Create(entity)
		if res.Error != nil {
			http.Error(w, res.Error.Error(), http.StatusInternalServerError)
			return
		}
		http.Redirect(w, r, "/"+entity.TableName(), http.StatusSeeOther)
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
		r.Use(middleware.WithAuthUserContext(db))
		r.Get("/", getAllPage)
		r.Route("/create", func(r chi.Router) {
			r.Get("/", getCreatePage)
			r.Route("/", func(r chi.Router) {
				r.Use(middleware.WithFormEntityContextFactory(entity))
				r.Use(middleware.WithEntityValidation)
				r.Post("/", create)
			})
		})
		r.Route("/{id}", func(r chi.Router) {
			r.Use(middleware.WithDbEntityContextFactory(entity, db))
			r.Get("/", getOnePage)
			// r.Delete("/", delete)
			r.Route("/", id_route)
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
	err = godotenv.Load(".company.env")
	if err != nil {
		log.Fatal("F failed to load company dotenv: ", err.Error())
	}
	var db *gorm.DB
	db, err = database.ConnectDb("main.db")
	if _, err := os.ReadFile("initialized"); err != nil {
		database.InitDb(db, "schema.sql")
		os.Create("initialized")
	}
	r := chi.NewRouter()
	r.Use(middleware.WithRequestInfoLogging)
	r.Get("/tailwind.js", func(w http.ResponseWriter, r *http.Request) {
		data, err := os.ReadFile("./tailwind.js")
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Write(data)
	})
	r.Route("/", func(r chi.Router) {
		r.Use(middleware.WithAuthUserContext(db))
		r.Get("/", func(w http.ResponseWriter, r *http.Request) { gui.MainPageComponent().Render(w) })
	})
	r.Route("/signup", func(r chi.Router) {
		r.Route("/", func(r chi.Router) {
			r.Use(middleware.WithFormFieldsValidationFactory([]string{"name", "password", "repeat_password"}))
			r.Post("/", func(w http.ResponseWriter, r *http.Request) { handler.SignupPost(w, r, db) })
		})
		r.Get("/", func(w http.ResponseWriter, r *http.Request) { gui.UserFormComponent(true).Render(w) })
	})
	r.Route("/login", func(r chi.Router) {
		r.Route("/", func(r chi.Router) {
			r.Use(middleware.WithFormFieldsValidationFactory([]string{"name", "password"}))
			r.Post("/", func(w http.ResponseWriter, r *http.Request) { handler.LoginPost(w, r, db) })
		})
		r.Get("/", func(w http.ResponseWriter, r *http.Request) { gui.UserFormComponent(false).Render(w) })
	})
	r.Route("/order", EntityRouterFactory(db, &OrderEntity{}, func(r chi.Router) {
		r.Route("/bill", func(r chi.Router) {
			r.Use(middleware.WithFormFieldsValidationFactory([]string{"date", "client_company"}))
			r.Get("/bill", func(w http.ResponseWriter, r *http.Request) { handler.GenerateOrderBill(w, r, db) })
		})
		r.Post("/end", func(w http.ResponseWriter, r *http.Request) { handler.EndOrder(w, r, db) })
	}))
	r.Route("/resource", EntityRouterFactory(db, &ResourceEntity{}, func(r chi.Router) {}))
	r.Route("/resource_resupply", EntityRouterFactory(db, &ResourceResupplyEntity{}, func(r chi.Router) {}))
	r.Route("/resource_spending", EntityRouterFactory(db, &OrderResourceSpendingEntity{}, func(r chi.Router) {}))
	r.Route("/item", EntityRouterFactory(db, &ItemEntity{}, func(r chi.Router) {}))
	r.Route("/order_item_fulfillment", EntityRouterFactory(db, &OrderItemFulfillmentEntity{}, func(r chi.Router) {}))
	fmt.Printf("I Listening on http://%s\n", addr)
	http.ListenAndServe(addr, r)
}
