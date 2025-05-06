package middleware

import (
	"context"
	"log"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/sergeykochiev/curs/backend/database/entity"
	"github.com/sergeykochiev/curs/backend/types"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func WithRequestInfoLogging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println(r.Method + " " + r.RequestURI)
		next.ServeHTTP(w, r)
	})
}

func WithAuthUserContext(db *gorm.DB) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if len(r.Cookies()) == 0 || len(r.CookiesNamed("token")) == 0 {
				http.Redirect(w, r, "/login", http.StatusSeeOther)
				return
			}
			int_id, err := strconv.Atoi(r.CookiesNamed("token")[0].Value)
			if err != nil {
				http.Error(w, "Failed to get userID", 404)
				return
			}
			user := entity.UserEntity{ID: int_id}
			res := db.First(&user)
			if res.Error != nil && res.Error != gorm.ErrRecordNotFound {
				http.Error(w, "Failed to get userID", 404)
				return
			} else if res.Error != nil {
				w.Header().Add("Location", "/login")
				w.WriteHeader(http.StatusSeeOther)
				return
			}
			ctx := context.WithValue(r.Context(), "user", user)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func WithDbEntityContextFactory[T interface {
	types.Identifier
}](entity T, db *gorm.DB) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			id := chi.URLParam(r, "id")
			res := db.Preload(clause.Associations).First(&entity, id)
			if res.Error != nil {
				http.Error(w, "ID not found", 404)
				return
			}
			// err := entity.LoadRelations(db)
			// if err != nil {
			// 	http.Error(w, "Error loading relations", 404)
			// 	return
			// }
			ctx := context.WithValue(r.Context(), "entity", entity)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func WithFormEntityContextFactory(entity types.FormParser) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			err := r.ParseForm()
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			if err = entity.ValidateAndParseForm(r); err != nil {
				http.Error(w, err.Error(), http.StatusUnprocessableEntity)
				return
			}
			ctx := context.WithValue(r.Context(), "entity", entity)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func WithFormFieldsValidationFactory(fields []string) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			err := r.ParseForm()
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			for _, f := range fields {
				if !r.Form.Has(f) {
					http.Error(w, "invalid formdata", http.StatusBadRequest)
				}
			}
			next.ServeHTTP(w, r)
		})
	}
}

func WithEntityValidation(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		entity := r.Context().Value("entity")
		if !entity.(types.Validator).Validate() {
			http.Error(w, "Invalid entity", http.StatusUnprocessableEntity)
			return
		}
		next.ServeHTTP(w, r)
	})
}
