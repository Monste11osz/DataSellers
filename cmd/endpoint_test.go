package main

import (
	"bytes"
	"log"
	"net/http"
	"net/http/httptest"
	"restforavito/pkg/postgres"
	"sync"
	"testing"
)

func TestStart(t *testing.T) {
	db, err := openDB("user=postgres password=qwerty123 dbname=postgres sslmode=disable host=localhost port=5432")
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	app := &application{
		product: &postgres.ProductMod{DB: db}}
	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		defer wg.Done()

		fileContent := []byte(`{"link": "https://tmpfiles.org/dl/4341598/-11--11.csv1.csv", "id": "1000"}`)
		req, err := http.NewRequest("POST", "/inputFile", bytes.NewBuffer(fileContent))
		if err != nil {
			log.Fatal(err)
			return
		}
		rr := httptest.NewRecorder()
		app.start(rr, req)
		if status := rr.Code; status != http.StatusOK {
			t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusOK)
		}
	}()

	go func() {
		defer wg.Done()

		fileContents := []byte(`{"link": "https://tmpfiles.org/dl/4341530/-11--11.csv2.csv", "id": "1000"}`)
		reqs, err := http.NewRequest("POST", "/inputFile", bytes.NewBuffer(fileContents))
		if err != nil {
			log.Fatal(err)
			return
		}
		rrs := httptest.NewRecorder()
		app.start(rrs, reqs)
		if status := rrs.Code; status != http.StatusOK {
			t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusOK)
		}
	}()
	wg.Wait()

}

func TestSearch(t *testing.T) {
	db, err := openDB("user=postgres password=qwerty123 dbname=postgres sslmode=disable host=localhost port=5432")
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	app := &application{
		product: &postgres.ProductMod{DB: db}}
	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		defer wg.Done()
		fileContent := []byte(`{"offerId": "1", "id": " ", "Name": "a"}`)
		req, err := http.NewRequest("POST", "/products/search", bytes.NewBuffer(fileContent))
		if err != nil {
			log.Fatal(err)
			return
		}
		rr := httptest.NewRecorder()
		app.search(rr, req)
		if status := rr.Code; status != http.StatusOK {
			t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusOK)
		}
	}()

	go func() {
		defer wg.Done()
		fileContents := []byte(`{"offerId": "1", "id": "1000", "Name": "a"}`)
		reqs, err := http.NewRequest("POST", "/products/search", bytes.NewBuffer(fileContents))
		if err != nil {
			log.Fatal(err)
			return
		}
		rrr := httptest.NewRecorder()
		app.search(rrr, reqs)
		if status := rrr.Code; status != http.StatusOK {
			t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusOK)
		}
	}()
	wg.Wait()

}
