package pendu

import (
	"fmt"
	"math/rand"
	"os"
	"sort"
	"strings"
	"time"
)

// variables globales
type Variables_pendu struct {
	Mot_a_trouver string
	Mot_actuel    string
	Essaie        int
	Liste_lettre  []string
	Phrase        string
}

func Jeux_pendu() {
	var Partie Variables_pendu
	Partie.Initialisation("words.txt") //on initialise le jeu
	for Partie.Essaie > 0 {            //boucle principale du jeu, s'arrête lorsque l'on perd
		Partie.Affichage_mot()                                       //on affiche le mot actuel
		Partie.Affichage_liste_lettre()                              //on affiche la liste des lettres déjà essayées
		Partie.Revelation_lettre(Partie.Entrée_utilisateur("salut")) //on demande à l'utilisateur de rentrer une lettre
		if Partie.Mot_actuel == Partie.Mot_a_trouver {               //condition de victoire
			fmt.Println("\n\nVous avez gagné !\nLe mot était bien :", Partie.Mot_a_trouver)
			os.Exit(0) //sortie du programme
		}
	}
	fmt.Println("\n\nVous avez perdu !\nLe mot était :", Partie.Mot_a_trouver)
}

func (Partie *Variables_pendu) Initialisation(fichier string) { //initialise le jeu
	Partie.Lecture_Fichier(fichier) //on lit le fichier donné en argument
}

func (Partie *Variables_pendu) Affichage_mot() {
	fmt.Println()
	for _, caractère := range Partie.Mot_actuel {
		fmt.Print(strings.ToUpper(string(caractère)), " ") // Permet d'afficher les lettres en majuscule avec un espace entre chaque
	}
}

func (Partie *Variables_pendu) Affichage_liste_lettre() {
	if len(Partie.Liste_lettre) == 0 { //si la liste est vide
		return
	}
	fmt.Print("Liste des essais : ")
	for _, lettre := range Partie.Liste_lettre {
		fmt.Print(lettre, " ") //affiche la liste des lettres déjà essayées avec un espace entre chaque
	}
	fmt.Println()
}

func (Partie *Variables_pendu) Affichage_pendu() {
	fichier, err := os.ReadFile("pendu/hangman.txt") //on lit le fichier
	if err != nil {                                  //si il y a une erreur
		Partie.Phrase = "Impossible d'ouvrir le fichier hangman.txt"
		print(Partie.Phrase)
		os.Exit(1) //on quitte le programme
	}
	var position int = 10 - (Partie.Essaie + 1) //on calcule la position de la ligne à afficher
	for i := 0; i < 7; i++ {                    //on boucle sur les 7 lignes du fichier
		for j := 0; j < 10; j++ { //on boucle sur les 10 caractères de la ligne
			Partie.Phrase = string(fichier[position*71+i*10+j])
		}
	}
}

func (Partie *Variables_pendu) Entrée_utilisateur(lettre string) string { //demande à l'utilisateur de choisir une lettre ou un mot
	if !Est_lettre(lettre) { //vérifie que l'utilisateur a bien entré que des lettres
		Partie.Phrase = "Merci de n'entrer que des lettres minusucules"
		return ""
	}
	if len(lettre) == 1 { //vérifie que l'utilisateur n'a pas entré plus d'une lettre
		for _, lettre_essaye := range Partie.Liste_lettre { //vérifie que l'utilisateur n'a pas déjà essayé cette lettre
			if strings.ToUpper(lettre) == lettre_essaye {
				Partie.Phrase = "Vous avez déjà essayé cette lettre, merci d'en choisir une autre"
				return ""
			}
		}
	}
	return strings.ToLower(lettre) //on retourne la lettre en minuscule
}

func (Partie *Variables_pendu) Lecture_Fichier(nom_fichier string) {
	var mot string
	var liste_mots []string
	fichier, err := os.ReadFile(nom_fichier) //on lit le fichier
	if err != nil {                          //si il y a une erreur
		Partie.Phrase = "Impossible d'ouvrir le fichier " + nom_fichier
		println(Partie.Phrase)
		println(err)
		os.Exit(1) //on quitte le programme
	}
	for index, caractère := range fichier {
		if caractère == 10 { //si on rencontre un retour à la ligne
			liste_mots = append(liste_mots, mot) //on l'ajoute à la liste
			mot = ""                             //on rénitialise mot
		} else if caractère != 13 { //si on ne rencontre pas un retour à la ligne
			mot += string(caractère)
		}
		if index == len(fichier)-1 { //on vérifie la fin
			liste_mots = append(liste_mots, mot)
		}
		if caractère == 32 || (!Est_lettre(string(caractère)) && caractère != 10 && caractère != 13) { //vérifie que le fichier ne contient que des lettres minuscules et des apostrophes
			Partie.Phrase = "Le fichier contient des caractères non autorisés, merci d'utiliser un fichier texte avec uniquement des lettres minuscules"
			print(Partie.Phrase)
			os.Exit(1) //on quitte le programme
		}
		Partie.Essaie = 10 //on initialise le nombre d'essaie
	}
	//utiliser un seed aléatoire permet d'éviter que le mot soit toujours le même lors de la même exécution du programme
	rand.Seed(time.Now().UnixNano())                              //on utilise le temps actuel comme seed
	Partie.Mot_a_trouver = liste_mots[rand.Intn(len(liste_mots))] //on choisit un mot aléatoire
	for i := 0; i < len(Partie.Mot_a_trouver); i++ {
		Partie.Mot_actuel += "_" //on initialise le mot actuel avec des _
	}
	for i := 0; i < len(Partie.Mot_a_trouver)/2-1; i++ {
		Partie.Mot_actuel = strings.Replace(Partie.Mot_actuel, "_", string(Partie.Mot_a_trouver[i]), 1) //on remplace (len(Mot_a_trouver)/2 -1) des _ par des lettres
	}
}

func (Partie *Variables_pendu) Revelation_lettre(lettre string) {
	if len(lettre) > 1 { //vérifie que l'utilisateur a rentré une lettre ou un mot
		if lettre == Partie.Mot_a_trouver { //vérifie que le mot entré est le bon
			Partie.Mot_actuel = Partie.Mot_a_trouver //on met le mot actuel à jour
		} else {
			Partie.Essaie -= 2     //on enlève 2 essaies
			if Partie.Essaie < 0 { //vérifie si il reste encore des essaies
				Partie.Essaie = 0 //on met les essai à 0 au cas ou ils étaient négatif
			}
			Partie.Affichage_pendu() //on affiche le pendu
			Partie.Phrase = "Votre mot est incorrect"
		}
	} else if len(lettre) == 1 { //vérifie que l'utilisateur a rentré une lettre
		var mot_temporaire string
		for index, caractère := range Partie.Mot_a_trouver { //on parcourt le mot à trouver
			if string(caractère) == lettre { //si la lettre est dans le mot
				mot_temporaire += lettre //on ajoute la lettre au mot temporaire
			} else {
				mot_temporaire += string(Partie.Mot_actuel[index]) //sinon on ajoute le caractère du mot actuel
			}
		}
		Partie.Liste_lettre = append(Partie.Liste_lettre, strings.ToUpper(lettre)) //on ajoute la lettre à la liste des lettres essayées
		sort.Strings(Partie.Liste_lettre)                                          //on trie la liste
		if mot_temporaire == Partie.Mot_actuel {                                   //si le mot temporaire est égal au mot actuel
			Partie.Essaie-- //on enlève un essaie
			Partie.Phrase = "La lettre n'est pas dans le mot"
		} else {
			Partie.Mot_actuel = mot_temporaire //on met le mot actuel à jour
			if Partie.Essaie != 10 {
				Partie.Phrase = "La lettre est dans le mot"
			} else {
				Partie.Phrase = "La lettre est dans le mot"
			}
		}
	}
}

func Est_lettre(str string) bool { //vérifie que la chaine de caractère ne contient que des lettres
	if len(str) == 0 {
		return false
	}
	for _, lettre := range str {
		if lettre < 'a' || lettre > 'z' {
			return false
		}
	}
	return true
}
