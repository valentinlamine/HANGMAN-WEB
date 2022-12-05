package main

import (
	"fmt"
	"hangman-web/pendu"
	"html/template"
	"net/http"
)

type User struct {
	Difficulte string
	Success    bool
}

func main() {
	fmt.Println("Lancement du serveur sur le port 80 : http://localhost")
	//gestion des css
	fs := http.FileServer(http.Dir("css"))
	http.Handle("/css/", http.StripPrefix("/css/", fs))
	//gestion des templates
	http.HandleFunc("/", IndexHandler)
	http.HandleFunc("/pendu", PenduHandler)
	http.ListenAndServe(":80", nil)
}

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseFiles("index.html")
	if r.Method != http.MethodPost {
		t.Execute(w, nil)
		return
	}
	difficulte := r.FormValue("difficulte")
	user := User{Difficulte: difficulte, Success: true}
	t.Execute(w, user)
}

func PenduHandler(w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseFiles("pendu.html")
	t.Execute(w, nil)
	pendu.Jeux_pendu()
}
