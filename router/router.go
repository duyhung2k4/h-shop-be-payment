package router

import (
	"app/controller"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/go-chi/render"
)

func Router() http.Handler {
	app := chi.NewRouter()

	app.Use(middleware.RequestID)
	app.Use(middleware.RealIP)
	app.Use(middleware.Logger)
	app.Use(middleware.Recoverer)

	cors := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	})
	app.Use(cors.Handler)

	// middlewares := middlewares.NewMiddlewares()

	paymentController := controller.NewPaymentController()

	app.Route("/payment/api/v1", func(r chi.Router) {
		r.Get("/ping", func(w http.ResponseWriter, r *http.Request) {
			res := controller.Response{
				Data:    "OK",
				Status:  200,
				Message: "OK",
				Error:   nil,
			}

			render.JSON(w, r, res)
		})
		r.Route("/protected", func(protected chi.Router) {
			// protected.Use(jwtauth.Verifier(config.GetJWT()))
			// protected.Use(jwtauth.Authenticator(config.GetJWT()))
			// protected.Use(middlewares.ValidateExpAccessToken())

			protected.Route("/payment", func(payment chi.Router) {
				payment.Post("/create-bill-online", paymentController.CreateBillOnline)
			})

		})
	})

	log.Println("Sevice h-shop-be-account starting success!")

	return app
}
