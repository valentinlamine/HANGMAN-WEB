package main

import (
	"fmt"
	"hangman-web/pendu"
	"html/template"
	"net/http"
)

type User struct {
	Difficulte string
	Username   string
	Success    bool
}

var Partie pendu.Variables_pendu

func main() {
	fmt.Println("Lancement du serveur sur le port 80 : http://localhost")
	//gestion des css
	fs := http.FileServer(http.Dir("css"))
	http.Handle("/css/", http.StripPrefix("/css/", fs))
	//gestion des images
	fs = http.FileServer(http.Dir("images"))
	http.Handle("/images/", http.StripPrefix("/images/", fs))
	//gestion du sous dossier positions dans le dossier images
	fs = http.FileServer(http.Dir("images/positions"))
	http.Handle("/images/positions/", http.StripPrefix("/images/positions/", fs))
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
	Joueur := User{"none", "none", false}
	if Joueur.Success {
		Joueur.Difficulte = r.FormValue("difficulte")
		t.Execute(w, Joueur)
		return
	} else {
		Joueur.Username = r.FormValue("username")
		Joueur.Success = true
		t.Execute(w, Joueur)
	}
}

func PenduHandler(w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseFiles("pendu.html")
	if Partie.Essaie == 0 {
		Partie = pendu.Variables_pendu{}
		if r.FormValue("difficulte") == "facile" {
			Partie.Initialisation("pendu/words.txt")
		} else if r.FormValue("difficulte") == "moyen" {
			Partie.Initialisation("pendu/words2.txt")
		} else if r.FormValue("difficulte") == "difficile" {
			Partie.Initialisation("pendu/words3.txt")
		}
		t.Execute(w, Partie)
	} else {
		Partie.Revelation_lettre(Partie.Entrée_utilisateur(r.FormValue("lettre")))
		if Partie.Mot_actuel == Partie.Mot_a_trouver {
			Partie.Phrase = "Vous avez gagné ! Le mot était bien : " + Partie.Mot_a_trouver
			Partie.Essaie = 0
		}
		if Partie.Essaie == 0 {
			Partie.Phrase = "Vous avez perdu ! Le mot était : " + Partie.Mot_a_trouver
		}
		t.Execute(w, Partie)
	}
}
