package main

import (
	"net/http"

	"github.com/Gaussgeek/bookings/internal/config"
	"github.com/Gaussgeek/bookings/internal/handlers"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

func routes(app *config.AppConfig) http.Handler {
	mux := chi.NewRouter()

	mux.Use(middleware.Recoverer)
	mux.Use(NoSurf)
	mux.Use(SessionLoad)

	mux.Get("/", handlers.Repo.Home)
	mux.Get("/about/", handlers.Repo.About)
	mux.Get("/presidential-suites/", handlers.Repo.Presidential)
	mux.Get("/platinum-rooms/", handlers.Repo.Platinum)

	mux.Get("/search-availability/", handlers.Repo.Availability)
	mux.Post("/search-availability/", handlers.Repo.PostAvailability)
	mux.Post("/search-availability-json/", handlers.Repo.AvailabilityJSON)
	mux.Get("/choose-room/{id}/", handlers.Repo.ChooseRoom)
	mux.Get("/book-room", handlers.Repo.BookRoom)

	mux.Get("/contact/", handlers.Repo.Contact)

	mux.Get("/make-reservation/", handlers.Repo.Reservation)
	mux.Post("/make-reservation/", handlers.Repo.PostReservation)
	mux.Get("/reservation-summary/", handlers.Repo.ReservationSummary)

	mux.Get("/user/login", handlers.Repo.ShowLogin)
	mux.Post("/user/login", handlers.Repo.PostShowLogin)
	mux.Get("/user/logout", handlers.Repo.Logout)

	//this directs where to find the static images
	fileServer := http.FileServer(http.Dir("./static/"))
	mux.Handle("/static/*", http.StripPrefix("/static", fileServer))

	mux.Route("/admin", func(mux chi.Router) { //all paths starting with admin, only available to admin
		//mux.Use(Auth)

		mux.Get("/dashboard", handlers.Repo.AdminDashboard)

		mux.Get("/reservations-new", handlers.Repo.AdminNewReservations)
		mux.Get("/reservations-all", handlers.Repo.AdminAllReservations)

		mux.Get("/reservations-calendar", handlers.Repo.AdminReservationsCalendar)
		mux.Post("/reservations-calendar", handlers.Repo.AdminPostReservationsCalendar)

		mux.Get("/process-reservation/{src}/{id}/do", handlers.Repo.AdminProcessReservation)
		mux.Get("/delete-reservation/{src}/{id}/do", handlers.Repo.AdminDeleteReservation)

		mux.Get("/reservations/{src}/{id}/show", handlers.Repo.AdminShowReservation)
		mux.Post("/reservations/{src}/{id}", handlers.Repo.AdminPostShowReservation)

	})

	return mux
}
