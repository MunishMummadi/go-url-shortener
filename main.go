package main

import (
	"database/sql"
	"flag"
	"fmt"
	"go-url-shortener/internals/models"
	"html/template"
	"log"
	"net/http"
	"os"

	"github.com/julienschmidt/httprouter"
	"github.com/justinas/alice"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
	"golang.org/x/text/number"

	_ "modernc.org/sqlite"
)

type PageData struct {
	BaseURL, Error string
	URLData []*ShortenerData
}

type App struct {
	urls *models.ShortenreDataModel
}

func serverError(w http.ResponseWriter, err error) {
	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

func newApp(dbFile string) App {
	db,err := sql.Open("sqlite", dbFile)
	if err != nil {
		log.Fatal(err)
	}

	if err = db.Ping(); err != nil {
		log.Fatal(err)
	}

	return App{urls: &models.ShortenerDataModel{DB: db}}
}

var functions = template.FuncMap{
	"formatClicks": formatClicks,
}

func formatClicks(clicks int) string {
	p := message.NewPrinter(language.English)
	return p.Sprintf("%v", number.Decimal(clicks))
}

func (a *App) getDefaultRoute(w http.ResponseWriter, r *http.Request) {
	tmplFile := "./templates/default.html"
	tmpl, err := template.New("default.html").Funcs(functions).ParserFiles(tmplFile)
	if err != nil {
		fmt.Println(err.Error())
		serverError(w, err)
		return
	}

	urls, err := a.urls.Latest()
	if err != nil {
		fmt.Printf("Could not retrieve all URLs, because %s. \n", err)
		return
	}
}
