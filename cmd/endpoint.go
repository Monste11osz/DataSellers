package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"time"
)

const COOKIE = "sessionId"

func (app *application) signIn(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		files := []string{"./templates/html/sign.html"}
		tmp, err := template.ParseFiles(files...)
		if err != nil {
			app.serverError(w, err)
			return
		}
		err = tmp.Execute(w, nil)
		if err != nil {
			app.serverError(w, err)
			return
		}

	}
	if r.Method == http.MethodPost {
		email := r.FormValue("username")
		password := r.FormValue("password")
		rep, key, err := app.product.CheckPassLogUser(email, password)
		if err != nil {
			http.ServeFile(w, r, "./templates/html/errlog.html")
		}
		if rep > 0 {
			//sessionId := sessionMem.Init(email)
			cookie := &http.Cookie{
				Name:    COOKIE,
				Value:   key,
				Expires: time.Now().Add(1 * time.Minute),
			}
			http.SetCookie(w, cookie)
			http.Redirect(w, r, "/us/form/inputFile", http.StatusSeeOther)
		} else {
			http.ServeFile(w, r, "./templates/html/errlog.html")

		}
	}
}

func (app *application) authentication(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		files := []string{"./templates/html/login.html"}
		tmp, err := template.ParseFiles(files...)
		if err != nil {
			app.serverError(w, err)
			return
		}
		err = tmp.Execute(w, nil)
		if err != nil {
			app.serverError(w, err)
			return
		}

	}
	if r.Method == http.MethodPost {
		email := r.FormValue("username")
		password := r.FormValue("password")
		sessionId := sessionMem.Init(email)
		//fmt.Println(email, password)
		rep, err := app.product.InputInfo(email, password, sessionId)
		if err != nil {
			http.Error(w, "Ошибка регистрации пользователя\n\tПовторите попытку!", http.StatusInternalServerError)
		}
		if rep > 0 {
			//http.Error(w, "Пользователь с таким Email существует\n\tПовторите попытку или осуществите вход!", http.StatusConflict)
			//http.Redirect(w, r, "/signIn", http.StatusSeeOther)
			http.ServeFile(w, r, "./templates/html/err.html")
		} else {
			//sessionId := sessionMem.Init(email)
			//cookie := &http.Cookie{
			//	Name:    COOKIE,
			//	Value:   sessionId,
			//	Expires: time.Now().Add(1 * time.Minute),
			//}
			//http.SetCookie(w, cookie)
			http.Redirect(w, r, "/signIn", http.StatusSeeOther)
		}
	}
}

func (app *application) start(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		files := []string{"./templates/html/input.html"}
		tmp, err := template.ParseFiles(files...)
		if err != nil {
			app.serverError(w, err)
			return
		}
		err = tmp.Execute(w, nil)
		if err != nil {
			app.serverError(w, err)
			return
		}

	} else {
		//if r.Method == http.MethodPost {
		//body, err := ioutil.ReadAll(r.Body)
		//if err != nil {
		//	app.serverError(w, err)
		//	return
		//}
		//var requestData map[string]string
		//if err := json.Unmarshal(body, &requestData); err != nil {
		//	app.serverError(w, err)
		//	return
		//}
		url := r.FormValue("link")
		id := r.FormValue("id")

		//url := requestData["link"]
		//id := requestData["id"]

		tmp, err := app.Download(url)
		if err != nil {
			app.serverError(w, err)
		}
		report, err := app.product.DownloadInfo(tmp, id)
		if err != nil {
			app.serverError(w, err)
			//log.Fatal(err)
		}
		responseJSON, err := json.Marshal(report)
		if err != nil {
			app.serverError(w, err)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(responseJSON)
		return
		//} else {
		//	w.Header().Set("Allow", http.MethodGet)
		//	http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		//	return
	}
}

func (app *application) search(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			app.serverError(w, err)
			return
		}
		var requestData map[string]string
		if err := json.Unmarshal(body, &requestData); err != nil {
			app.serverError(w, err)
			return
		}
		offId := requestData["offerId"]
		sellId := requestData["id"]
		word := requestData["Name"]
		putInf, err := app.product.SearchInfo(offId, sellId, word)
		if err != nil {
			app.serverError(w, err)
		}
		fmt.Println(putInf)
		responseJSON, err := json.Marshal(putInf)
		if err != nil {
			app.serverError(w, err)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(responseJSON)
		return
	} else {
		w.Header().Set("Allow", http.MethodGet)
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}
}

func (app *application) Download(url string) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		app.errorLog.Println(err.Error())
		return "", err
	}
	defer resp.Body.Close()
	tmp, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		app.errorLog.Println(err.Error())
		return "", err
	}
	//fmt.Println(string(tmp))
	return string(tmp), nil
}
