package main

import (
	"crypto/rand"
	"crypto/rsa"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"runtime"

	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"
	billgen "github.com/sergeykochiev/billgen/gen"
	billgen_init "github.com/sergeykochiev/billgen/init"
	"github.com/sergeykochiev/curs/backend/database"
	. "github.com/sergeykochiev/curs/backend/database/entity"
	"github.com/sergeykochiev/curs/backend/database/entity/report"
	"github.com/sergeykochiev/curs/backend/gui"
	"github.com/sergeykochiev/curs/backend/handler"
	"github.com/sergeykochiev/curs/backend/middleware"
	"github.com/sergeykochiev/curs/backend/types"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

// --- idea - https://go.dev/wiki/LockOSThread
func init() {
	runtime.LockOSThread()
}

var main_queue = make(chan func())

//---

const addr = "localhost:3003"

func EntityRouterFactory[T types.Entity](db *gorm.DB, entity T, id_route func(r chi.Router)) func(r chi.Router) {
	preloadedDb := entity.GetPreloadedDb(db)
	return func(r chi.Router) {
		r.Get("/", handler.CreateEntityGetAllPageHandler(entity, preloadedDb))
		r.Route("/create", func(r chi.Router) {
			r.Get("/", handler.CreateEntityCreatePageHandler(entity, db))
			r.Route("/", func(r chi.Router) {
				r.Use(middleware.WithFormEntityContextFactory(entity))
				r.Use(middleware.WithEntityValidation)
				r.Post("/", handler.CreateEntityCreateHandler(entity, preloadedDb))
			})
		})
		r.Route("/{id}", func(r chi.Router) {
			r.Use(middleware.WithDbEntityContextFactory(entity, db))
			r.Get("/", handler.CreateEntityGetPageHandler(entity))
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
	schema.RegisterSerializer("decimal", database.DecimalIdSerializer{})
	if err = billgen_init.Init(); err != nil {
		log.Fatal("F failed to init wkhtmltopdf from billgen: ", err.Error())
	}
	defer billgen_init.Destroy()
	err = godotenv.Load(".company.env")
	if err != nil {
		log.Fatal("F failed to load company dotenv: ", err.Error())
	}
	var key *rsa.PrivateKey
	key, err = rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		log.Fatal("F failed to generate RSA key for JWT: ", err.Error())
	}
	var db *gorm.DB
	db, err = database.Connect("main.db")
	if _, err := os.ReadFile("initialized"); err != nil {
		if err = database.ExecuteFile(db, "schema.sql"); err != nil {
			log.Fatal("F failed to init db: ", err.Error())
		}
		if err = database.ExecuteFile(db, "static_data.sql"); err != nil {
			log.Fatal("F failed to seed db: ", err.Error())
		}
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
	r.Route("/signup", func(r chi.Router) {
		r.Route("/", func(r chi.Router) {
			r.Use(middleware.WithFormFieldsValidationFactory([]string{"name", "password", "repeat_password"}))
			r.Post("/", func(w http.ResponseWriter, r *http.Request) { handler.SignupPost(w, r, db) })
		})
		r.Get("/", func(w http.ResponseWriter, r *http.Request) { gui.UserFormPage(true).Render(w) })
	})
	r.Route("/login", func(r chi.Router) {
		r.Route("/", func(r chi.Router) {
			r.Use(middleware.WithFormFieldsValidationFactory([]string{"name", "password"}))
			r.Post("/", handler.CreateLoginPostHandler(db, key))
		})
		r.Get("/", func(w http.ResponseWriter, r *http.Request) { gui.UserFormPage(false).Render(w) })
	})
	r.Route("/", func(r chi.Router) {
		r.Use(middleware.WithAuthUserContext(db, &key.PublicKey))
		r.Get("/", func(w http.ResponseWriter, r *http.Request) { gui.MainPage().Render(w) })
		r.Route("/order", EntityRouterFactory(db, &OrderEntity{}, func(r chi.Router) {
			r.Get("/bill", func(w http.ResponseWriter, r *http.Request) {
				handler.GenerateOrderBill(w, r, db, billgen.CreateBillPdf, &main_queue)
			})
			r.Get("/invoice", func(w http.ResponseWriter, r *http.Request) {
				handler.GenerateOrderBill(w, r, db, billgen.CreateInvoicePdf, &main_queue)
			})
			r.Get("/end", func(w http.ResponseWriter, r *http.Request) { handler.EndOrder(w, r, db) })
		}))
		r.Route("/resource", EntityRouterFactory(db, &ResourceEntity{}, func(r chi.Router) {}))
		r.Route("/resource_resupply", EntityRouterFactory(db, &ResourceResupplyEntity{}, func(r chi.Router) {}))
		r.Route("/resource_spending", EntityRouterFactory(db, &OrderResourceSpendingEntity{}, func(r chi.Router) {}))
		r.Route("/item", EntityRouterFactory(db, &ItemEntity{}, func(r chi.Router) {}))
		r.Route("/order_item_fulfillment", EntityRouterFactory(db, &OrderItemFulfillmentEntity{}, func(r chi.Router) {}))
		r.Route("/item_resource_need", EntityRouterFactory(db, &ItemResourceNeed{}, func(r chi.Router) {}))
		r.Route("/item_popularity", func(r chi.Router) {
			r.Get("/", func(w http.ResponseWriter, r *http.Request) {
				gui.DatedReportFormPage("Создать отчет о популярности товаров/услуг").Render(w)
			})
			r.Post("/", handler.CreateGenerateDatedReportHandler(db, &main_queue, &report.ItemPopularity{}))
		})
		r.Route("/resource_spendings", func(r chi.Router) {
			r.Get("/", func(w http.ResponseWriter, r *http.Request) {
				gui.DatedReportFormPage("Создать отчет о тратах ресурсов").Render(w)
			})
			r.Post("/", handler.CreateGenerateDatedReportHandler(db, &main_queue, &report.ResourceSpending{}))
		})
	})
	fmt.Printf("I Listening on http://%s\n", addr)
	go http.ListenAndServe(addr, r)
	var quit = make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	for {
		select {
		case f := <-main_queue:
			f()
		case <-quit:
			log.Println("Closing main queue")
			return
		}
	}
}
