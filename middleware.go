package main

import (
	"context"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/sergeykochiev/curs/backend/types"
	"gorm.io/gorm"
)

func withRequestInfoLogging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println(r.Method + " " + r.RequestURI)
		next.ServeHTTP(w, r)
	})
}

func withAuthUserIdContext(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if len(r.Cookies()) == 0 || len(r.CookiesNamed("token")) == 0 {
			w.Header().Add("Location", "/login")
			w.WriteHeader(http.StatusSeeOther)
			return
		}
		ctx := context.WithValue(r.Context(), "userid", r.CookiesNamed("token")[0].Value)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func withDbEntityContextFactory(entity types.Identifier, db *gorm.DB) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			id := chi.URLParam(r, "id")
			res := db.First(&entity, id)
			if res.Error != nil {
				http.Error(w, "ID not found", 404)
				return
			}
			ctx := context.WithValue(r.Context(), "entity", entity)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func withFormEntityContextFactory(entity types.FormParser) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			err := r.ParseForm()
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			r.Form.Add("userid", r.Context().Value("userid").(string))
			if !entity.ValidateAndParseForm(r.Form) {
				http.Error(w, "Invalid formdata", http.StatusUnprocessableEntity)
				return
			}
			ctx := context.WithValue(r.Context(), "entity", entity)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func withFormFieldsValidationFactory(fields []string) func(next http.Handler) http.Handler {
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

func withEntityValidation(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		entity := r.Context().Value("entity")
		if entity.(types.Validator).Validate() {
			http.Error(w, "Invalid entity", http.StatusUnprocessableEntity)
			return
		}
		next.ServeHTTP(w, r)
	})
}
