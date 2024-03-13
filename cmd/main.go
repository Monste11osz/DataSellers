package main

import (
	"database/sql"
	"flag"
	_ "github.com/lib/pq"
	"log"
	"net/http"
	"os"
	"restforavito/cmd/session"
	"restforavito/pkg/postgres"
)

type application struct {
	product  *postgres.ProductMod
	ans      *postgres.Answer
	Prod     postgres.Prod
	Lis      postgres.ListProd
	infoLog  *log.Logger
	errorLog *log.Logger
}

var sessionMem *session.Session

func main() {
	dsn := flag.String("dsn", "user=postgres password=qwerty123 dbname=postgres sslmode=disable host=localhost port=5432", "Строка подключения к PostgreSQL")
	flag.Parse()
	addr := flag.String("addr", ":4002", "server")
	flag.Parse()

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	db, err := openDB(*dsn)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	sessionMem = session.NewSession()
	app := &application{
		infoLog:  infoLog,
		errorLog: errorLog,
		product:  &postgres.ProductMod{DB: db},
		ans:      &postgres.Answer{},
		Prod:     postgres.Prod{},
		Lis:      postgres.ListProd{},
	}
	srv := &http.Server{
		Addr:    *addr,
		Handler: app.router(),
	}

	err = srv.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	if err = db.Ping(); err != nil {
		log.Fatal(err)
		return nil, err
	}
	return db, err
}
