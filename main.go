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
