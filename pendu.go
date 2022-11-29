package main

import (
	"fmt"
	"math/rand"
	"os"
	"sort"
	"strings"
)

// variables globales
type variables struct {
	Mot_a_trouver string
	Mot_actuel    string
	Essaie        int
	Liste_lettre  []string
}

func jeux_pendu() {
	var Partie variables
	Partie.Initialisation() //on initialise le jeu
	for Partie.Essaie > 0 { //boucle principale du jeu, s'arrête lorsque l'on perd
		Affichage_mot(Partie)                                  //on affiche le mot actuel
		Affichage_liste_lettre(Partie)                         //on affiche la liste des lettres déjà essayées
		Revelation_lettre(Entrée_utilisateur(Partie), &Partie) //on demande à l'utilisateur de rentrer une lettre
		if Partie.Mot_actuel == Partie.Mot_a_trouver {         //condition de victoire
			fmt.Println("\n\nVous avez gagné !\nLe mot était bien :", Partie.Mot_a_trouver)
			os.Exit(0) //sortie du programme
		}
	}
	fmt.Println("\n\nVous avez perdu !\nLe mot était :", Partie.Mot_a_trouver)
}

func (Partie *variables) Initialisation() {
	Partie.Essaie = 10     //on initialise le nombre d'essaie
	if len(os.Args) != 2 { //vérifie qu'il y a bien un argument
		fmt.Print("\nMerci d'indiquer le nom du fichier texte à utiliser : \nExemple : go run main.go words.txt\n\n")
		os.Exit(1) //sinon, on quitte le programme
	} else {
		Lecture_Fichier(os.Args[1], Partie) //on lit le fichier donné en argument
	}
	Affichage_espace() //on affiche un espace pour faire un affichage propre
	fmt.Println("Bienvenue dans le jeu du pendu !")
	fmt.Println("Bonne chance, vous avez 10 essaies")
	fmt.Println("\nNote : Le programme affiche des lettres dès le lancement, Toutefois il n'affiche pas pour autant toutes les occurences de ces lettres")
}

func Affichage_espace() { //pour faire un affichage propre
	fmt.Print("\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n")
	fmt.Print("\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n")
	fmt.Print("\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n")
}

func Affichage_mot(Partie variables) {
	fmt.Println()
	for _, caractère := range Partie.Mot_actuel {
		fmt.Print(strings.ToUpper(string(caractère)), " ") // Permet d'afficher les lettres en majuscule avec un espace entre chaque
	}
	fmt.Print("\n\n")
}

func Affichage_liste_lettre(Partie variables) {
	if len(Partie.Liste_lettre) == 0 { //si la liste est vide
		return
	}
	fmt.Print("Liste des essais : ")
	for _, lettre := range Partie.Liste_lettre {
		fmt.Print(lettre, " ") //affiche la liste des lettres déjà essayées avec un espace entre chaque
	}
	fmt.Println()
}

func Affichage_pendu(Partie variables) {
	Affichage_espace()                         //on affiche un espace pour faire un affichage propre
	fichier, err := os.ReadFile("hangman.txt") //on lit le fichier
	if err != nil {                            //si il y a une erreur
		fmt.Println("Impossible d'ouvrir le fichier")
		os.Exit(1) //on quitte le programme
	}
	var position int = 10 - (Partie.Essaie + 1) //on calcule la position de la ligne à afficher
	for i := 0; i < 7; i++ {                    //on boucle sur les 7 lignes du fichier
		for j := 0; j < 10; j++ { //on boucle sur les 10 caractères de la ligne
			fmt.Print(string(fichier[position*71+i*10+j])) //on affiche le caractère
		}
	}
	fmt.Println()
}

func Entrée_utilisateur(Partie variables) string { //demande à l'utilisateur de choisir une lettre ou un mot
	var lettre string
	fmt.Print("Choix : ")
	fmt.Scanln(&lettre)      //on récupère l'entrée utilisateur
	if !Est_lettre(lettre) { //vérifie que l'utilisateur a bien entré que des lettres
		fmt.Println("Merci de n'entrer que des lettres minusucules")
		return Entrée_utilisateur(Partie) //on relance la fonction
	}
	if len(lettre) == 1 { //vérifie que l'utilisateur n'a pas entré plus d'une lettre
		for _, lettre_essaye := range Partie.Liste_lettre { //vérifie que l'utilisateur n'a pas déjà essayé cette lettre
			if strings.ToUpper(lettre) == lettre_essaye {
				fmt.Println("Vous avez déjà essayé cette lettre, merci d'en choisir une autre")
				return Entrée_utilisateur(Partie) //on relance la fonction
			}
		}
	}
	return strings.ToLower(lettre) //on retourne la lettre en minuscule
}

func Lecture_Fichier(nom_fichier string, Partie *variables) {
	var mot string
	var liste_mots []string
	fichier, err := os.ReadFile(nom_fichier) //on lit le fichier
	if err != nil {                          //si il y a une erreur
		fmt.Println("fichier introuvable ou illisible")
		os.Exit(1) //on quitte le programme
	}
	for index, caractère := range fichier {
		if caractère == 10 { //lorsque l'on va à la ligne le mot est fini
			liste_mots = append(liste_mots, mot) //on l'ajoute à la liste
			mot = ""                             //on rénitialise mot
		} else { //sinon on ajoute le caractère a mot
			mot += string(caractère)
		}
		if index == len(fichier)-1 { //on vérifie la fin
			liste_mots = append(liste_mots, mot)
		}
		if caractère == 32 || (!Est_lettre(string(caractère)) && caractère != 10) { //si le caractère est un espace
			fmt.Println("Le fichier contient des caractères non autorisés, merci d'utiliser un fichier texte avec uniquement des lettres minuscules")
			os.Exit(1) //on quitte le programme
		}
	}
	rand.Seed(int64(os.Getpid()))                                 //on initialise le générateur de nombre aléatoire
	Partie.Mot_a_trouver = liste_mots[rand.Intn(len(liste_mots))] //on choisit un mot aléatoire
	for i := 0; i < len(Partie.Mot_a_trouver); i++ {
		Partie.Mot_actuel += "_" //on initialise le mot actuel avec des _
	}
	for i := 0; i < len(Partie.Mot_a_trouver)/2-1; i++ {
		Partie.Mot_actuel = strings.Replace(Partie.Mot_actuel, "_", string(Partie.Mot_a_trouver[i]), 1) //on remplace (len(Mot_a_trouver)/2 -1) des _ par des lettres
	}
}

func Revelation_lettre(lettre string, Partie *variables) {
	if len(lettre) != 1 { //vérifie que l'utilisateur a rentré une lettre ou un mot
		if lettre == Partie.Mot_a_trouver { //vérifie que le mot entré est le bon
			Partie.Mot_actuel = Partie.Mot_a_trouver //on met le mot actuel à jour
		} else {
			Partie.Essaie -= 2     //on enlève 2 essaies
			if Partie.Essaie < 0 { //vérifie si il reste encore des essaies
				Partie.Essaie = 0 //on met les essai à 0 au cas ou ils étaient négatif
			}
			Affichage_pendu(*Partie) //on affiche le pendu
			fmt.Println("Votre mot est incorrect, il vous reste", Partie.Essaie, "essaies")
		}
	} else {
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
			Partie.Essaie--          //on enlève un essaie
			Affichage_pendu(*Partie) //on affiche le pendu
			fmt.Println("La lettre n'est pas dans le mot, il vous reste", Partie.Essaie, "essaies :")
		} else {
			Partie.Mot_actuel = mot_temporaire //on met le mot actuel à jour
			if Partie.Essaie != 10 {
				Affichage_pendu(*Partie) //on affiche le pendu
				fmt.Println("La lettre est dans le mot, il vous reste", Partie.Essaie, "essaies :")
			} else {
				Affichage_espace() //on affiche un espace
				fmt.Println("La lettre est dans le mot, il vous reste", Partie.Essaie, "essaies :")
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
