package main

import (
	"fmt"
	"hangman-web/pendu"
	"html/template"
	"net/http"
)

type User struct {
	Difficulty string
	Username   string
	Success    bool
	//statistiques
	Victoire                 int
	Parties                  int
	Victoire_consecutives    int
	Derniere_partie_victoire bool
	//pendu
	Pendu pendu.Variables_pendu
}

// On initialise le joueur
var Joueur User = User{"none", "none", false, 0, 0, 0, false, pendu.Variables_pendu{}}

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
	http.HandleFunc("/hangman", HangmanHandler)
	//lancement du serveur
	http.ListenAndServe(":80", nil)
}

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseFiles("index.html")
	if Joueur.Pendu.Essaie != 0 { //si le joueur a quitté la page pendant une partie
		Joueur.Pendu = pendu.Variables_pendu{}
		Joueur.Victoire_consecutives = 0
		Joueur.Derniere_partie_victoire = false
	}
	if r.Method != http.MethodPost { //si le joueur n'a pas rempli le formulaire
		if Joueur.Success { //lorsque le joueur a fini une partie et qu'il recommence, il n'a pas à renseigner son nom, on demande la difficulté
			Joueur.Difficulty = r.FormValue("difficulty")
			t.Execute(w, Joueur)
			return //on sort de la fonction pour éviter de réécrire le template
		}
		t.Execute(w, nil)
		return //on sort de la fonction pour éviter de réécrire le template
	}
	if Joueur.Success { //lorque le joueur a déjà renseigné son nom et qu'il choisit une difficulté
		Joueur.Difficulty = r.FormValue("difficulty")
		t.Execute(w, Joueur)
		return
	} else { //si le joueur a rempli le formulaire pour la première fois
		Joueur.Username = r.FormValue("username")
		Joueur.Success = true
		t.Execute(w, Joueur)
	}
}

func HangmanHandler(w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseFiles("hangman.html")
	if Joueur.Pendu.Essaie == 0 { //si le joueur n'a pas de partie en cours ou qu'il a fini une partie
		Joueur.Pendu = pendu.Variables_pendu{}     //on réinitialise les variables
		if r.FormValue("difficulty") == "facile" { //on initialise le pendu en fonction de la difficulté choisie
			Joueur.Pendu.Initialisation("pendu/words.txt")
		} else if r.FormValue("difficulty") == "moyen" { //on initialise le pendu en fonction de la difficulté choisie
			Joueur.Pendu.Initialisation("pendu/words2.txt")
		} else if r.FormValue("difficulty") == "difficile" { //on initialise le pendu en fonction de la difficulté choisie
			Joueur.Pendu.Initialisation("pendu/words3.txt")
		}
		t.Execute(w, Joueur)
	} else { //si le joueur a déjà commencé la partie
		Joueur.Pendu.Revelation_lettre(Joueur.Pendu.Entrée_utilisateur(r.FormValue("lettre"))) //on révèle les lettres correspondantes
		blocage_double := false                                                                //on bloque le double appel de la fonction de fin de partie
		if Joueur.Pendu.Mot_actuel == Joueur.Pendu.Mot_a_trouver {                             //si le joueur a gagné
			Joueur.Pendu.Phrase = "Vous avez gagné ! Le mot était bien : " + Joueur.Pendu.Mot_a_trouver
			//on mets les essaie à 0 pour forcer la fin de la partie
			Joueur.Pendu.Essaie = 0
			//on mets à jour les statistiques
			blocage_double = true
			Joueur.Victoire++
			Joueur.Parties++
			Joueur.Victoire_consecutives++
			Joueur.Derniere_partie_victoire = true
		}
		if Joueur.Pendu.Essaie == 0 && !blocage_double { //si le joueur a perdu
			Joueur.Pendu.Phrase = "Vous avez perdu ! Le mot était : " + Joueur.Pendu.Mot_a_trouver
			//pas besoin de mettre les essaie à 0, on le fait plus haut au lancement d'une nouvelle partie
			//on mets à jour les statistiques
			Joueur.Victoire_consecutives = 0
			Joueur.Parties++
			Joueur.Derniere_partie_victoire = false
		}
		t.Execute(w, Joueur)
	}
}
