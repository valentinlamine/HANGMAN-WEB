package main

import (
	"html/template"
	"net/http"
)

type User struct {
	Difficulte string
	Success    bool
}

// le but est d'intégrer le code de pendu.go dans le code de main.go
func main() {
	println("Le serveur est lancé sur le port 80")
	template_choix_difficulte := template.Must(template.ParseFiles("index.html")) //on parse le fichier html
	fs := http.FileServer(http.Dir("css_pierre"))                                 //on définit le dossier css
	http.Handle("/css_pierre/", http.StripPrefix("/css_pierre/", fs))             //on définit le dossier css
	//integrer les images présente dans les dossier images
	img := http.FileServer(http.Dir("images"))
	http.Handle("/images/", http.StripPrefix("/images/", img))
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) { //on définit la page d'accueil
		if r.Method != http.MethodPost { //si la méthode n'est pas POST
			template_choix_difficulte.Execute(w, nil) //on affiche la page d'accueil
			return
		} //sinon on récupère les données du formulaire
		difficulte := r.FormValue("difficulte")
		println(difficulte)
		details := User{
			Difficulte: r.FormValue("difficulte"),
			Success:    true,
		}
		template_choix_difficulte.Execute(w, details)
	})
	http.ListenAndServe(":80", nil)
}
