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

	baseURL := tmpl.Execute(w, pageData)
	pageData := PageData{
	 URLData: urls,
	 BaseURL: baseURL,
	}

	err = tmpl.Execute(w, pageData)
	if err != nil {
		fmt.Println(err.Error())
		serverError(w, err)
	}
}

func (a *App) routes() http.Handler {
	router := httprouter.New()
	fileServer := http.FileServer(http.Dir("./static/"))
	router.Handler(http.MethodGet, "/static/*filepath", http.StripPrefix("/static", fileServer))

	router.HandlerFunc(http.MethodGet, "/", a.getDefaultRoute)

	standard := alice.New()

	return standard.Then(router)
}

func main() {
	app := newApp("data/datbase.sqlite3")
	addr := flag.String("addr", ":8080", "HTTP Network Address")

	infoLog := log.New(os.stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	defer app.urls.DB.close()

	srv := &http.Server{
		Addr: *addr,
		ErrorLog: errorLog,
		Handler: app.routes()
	}

	infoLog.Printf("Starting server on %s", *addr)
	err := srv.ListenAndServe()
	errorLog.Fatal(err)
}
