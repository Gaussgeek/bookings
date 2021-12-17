package main

import (
	"encoding/gob"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/Gaussgeek/bookings/internal/config"
	"github.com/Gaussgeek/bookings/internal/driver"
	"github.com/Gaussgeek/bookings/internal/handlers"
	"github.com/Gaussgeek/bookings/internal/helpers"
	"github.com/Gaussgeek/bookings/internal/models"
	"github.com/Gaussgeek/bookings/internal/render"
	"github.com/alexedwards/scs/v2"
)

const portNumber = ":8080"

var app config.AppConfig
var session *scs.SessionManager
var infoLog *log.Logger
var errorLog *log.Logger

// main is the main function
func main() {
	db, err := run()
	if err != nil {
		log.Fatal(err)
	}
	defer db.SQL.Close() //this closes the database

	defer close(app.MailChan) //this closes the email channel

	fmt.Println("Starting mail listener....")
	listenForMail()

	fmt.Println(fmt.Sprintf("Staring application on port %s ", portNumber))

	srv := &http.Server{
		Addr:    portNumber,
		Handler: routes(&app),
	}

	err = srv.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}

func run() (*driver.DB, error) {

	// what am I going to put in the session
	gob.Register(models.Reservation{})
	gob.Register(models.User{})
	gob.Register(models.Room{})
	gob.Register(models.Restriction{})
	gob.Register(map[string]int{})

	//create a channel of type models.MailData
	mailChan := make(chan models.MailData)
	app.MailChan = mailChan

	
	app.InProduction = false // change this to true when in production

	infoLog = log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	app.InfoLog = infoLog

	errorLog = log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
	app.ErrorLog = errorLog

	// set up the session
	session = scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = app.InProduction

	app.Session = session

	// connect to a database
	log.Println("Connecting to database......")
	db, err := driver.ConnectSQL("host=localhost port=5432 dbname=bookings user=andrewl password=")
	if err != nil {
		log.Fatal("Cannot connect to the database! Dying...")
	}
	log.Println("Connected to the database!")

	tc, err := render.CreateTemplateCache()
	if err != nil {
		log.Fatal("cannot create template cache")
		return nil, err
	}

	app.TemplateCache = tc
	app.UseCache = false

	repo := handlers.NewRepo(&app, db)
	handlers.NewHandlers(repo)

	render.NewRenderer(&app)

	helpers.NewHelpers(&app)

	return db, nil
}
