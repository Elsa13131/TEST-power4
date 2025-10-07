package menu

import (
	"fmt"
	"net/http"
	"strconv"      // Pour convertir string en int
	"strings"      // Pour manipuler les strings
)

// STRUCTURE POUR LE JEU
type Jeu struct {
	// Structure qui représente l'état du jeu
	Plateau     [6][7]string  // Tableau 6 lignes × 7 colonnes pour stocker les pions
	TourActuel  string        // "rouge" ou "jaune" pour savoir à qui c'est le tour
	Partie      bool          // true si partie en cours, false si terminée
}

// VARIABLE GLOBALE POUR STOCKER L'ÉTAT DU JEU
var jeuActuel = Jeu{
	// Initialisation du jeu
	Plateau:    [6][7]string{},  // Plateau vide au début
	TourActuel: "rouge",          // Rouge commence
	Partie:     true,             // Partie active
}


func WelcomeHandler(w http.ResponseWriter, r *http.Request) {
	// Redirige vers la page "/next"
	http.Redirect(w, r, "/next", http.StatusSeeOther)
}

// HomeHandler gère la page d'accueil
// Sert le fichier HTML templates.html
func HomeHandler(w http.ResponseWriter, r *http.Request) {
	// ServeFile envoie le fichier HTML au navigateur
	http.ServeFile(w, r, "templates/templates.html")
}

// NextHandler gère la page suivante après redirection
// Redirige maintenant vers la page de jeu dynamique
func NextHandler(w http.ResponseWriter, r *http.Request) {
	/* fmt.Fprintln(w, "Vous êtes sur la page suivante !") */
	// Utilise maintenant GameHandler pour afficher le jeu fonctionnel
	GameHandler(w, r)
}

// FONCTION POUR AFFICHER LA PAGE DE JEU
func GameHandler(w http.ResponseWriter, r *http.Request) {
	// Fonction qui gère l'affichage de la page de jeu
	
	// Générer le HTML avec l'état actuel du jeu
	html := genererHTML()
	// Appelle la fonction qui crée le HTML avec les bonnes couleurs
	
	// Envoyer le HTML au navigateur
	w.Header().Set("Content-Type", "text/html")
	// Dit au navigateur que c'est du HTML
	
	fmt.Fprint(w, html)
	// Écrit le HTML dans la réponse
}

// FONCTION POUR JOUER UN COUP
func JouerHandler(w http.ResponseWriter, r *http.Request) {
	// Fonction qui gère quand un joueur clique sur une colonne
	
	// Extraire le numéro de colonne de l'URL
	path := r.URL.Path
	// Récupère le chemin de l'URL (ex: "/jouer/3")
	
	parts := strings.Split(path, "/")
	// Divise le chemin en parties : ["", "jouer", "3"]
	
	if len(parts) != 3 {
		// Si l'URL n'a pas exactement 3 parties, erreur
		http.Redirect(w, r, "/jeu", http.StatusSeeOther)
		return
	}
	
	colonne, err := strconv.Atoi(parts[2])
	// Convertit "3" en nombre 3
	
	if err != nil || colonne < 1 || colonne > 7 {
		// Si conversion échoue ou colonne invalide
		http.Redirect(w, r, "/jeu", http.StatusSeeOther)
		return
	}
	
	// Jouer le coup
	fmt.Println("DEBUG: Avant jouerCoup, colonne =", colonne-1)  // Debug
	jouerCoup(colonne - 1)  // -1 car les tableaux commencent à 0
	fmt.Println("DEBUG: Après jouerCoup, état du plateau:")  // Debug
	for i := 0; i < 6; i++ {
		fmt.Println("  ", jeuActuel.Plateau[i])  // Debug
	}
	// Appelle la fonction qui place le pion
	
	// Rediriger vers la page de jeu mise à jour
	http.Redirect(w, r, "/jeu", http.StatusSeeOther)
	// Recharge la page pour afficher le nouveau coup
}

// FONCTION POUR NOUVEAU JEU
func NouveauJeuHandler(w http.ResponseWriter, r *http.Request) {
	// Fonction pour recommencer une partie
	
	// Réinitialiser le jeu
	jeuActuel = Jeu{
		Plateau:    [6][7]string{},  // Plateau vide
		TourActuel: "rouge",          // Rouge recommence
		Partie:     true,             // Partie active
	}
	
	// Rediriger vers la page de jeu
	http.Redirect(w, r, "/jeu", http.StatusSeeOther)
}

// FONCTION POUR JOUER UN COUP
func jouerCoup(colonne int) {
	// Fonction qui place un pion dans une colonne
	
	// Chercher la case la plus basse disponible
	for ligne := 5; ligne >= 0; ligne-- {
		// Parcourt de bas en haut (ligne 5 à 0)
		
		if jeuActuel.Plateau[ligne][colonne] == "" {
			// Si la case est vide
			
			jeuActuel.Plateau[ligne][colonne] = jeuActuel.TourActuel
			// Place le pion de la couleur du joueur actuel
			
			// Changer de joueur
			if jeuActuel.TourActuel == "rouge" {
				jeuActuel.TourActuel = "jaune"
			} else {
				jeuActuel.TourActuel = "rouge"
			}
			
			break  // Sortir de la boucle
		}
	}
}

// FONCTION POUR GÉNÉRER LE HTML
func genererHTML() string {
	// Fonction qui crée le HTML avec l'état actuel du jeu
	
	html := `<!DOCTYPE html>
<html lang="fr">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Puissance 4 - Jeu</title>
    <link rel="stylesheet" href="/assets/game.css">
</head>
<body>
    <h1>Puissance 4</h1>
    <p id="tour-joueur">Tour du joueur : ` + strings.Title(jeuActuel.TourActuel) + `</p>
    
    <div class="boutons-colonnes">
        <a href="/jouer/1"><button>Col 1</button></a>
        <a href="/jouer/2"><button>Col 2</button></a>
        <a href="/jouer/3"><button>Col 3</button></a>
        <a href="/jouer/4"><button>Col 4</button></a>
        <a href="/jouer/5"><button>Col 5</button></a>
        <a href="/jouer/6"><button>Col 6</button></a>
        <a href="/jouer/7"><button>Col 7</button></a>
    </div>
    
    <table id="plateau">`
	
	// Générer les lignes du tableau
	for ligne := 0; ligne < 6; ligne++ {
		html += "<tr>"
		for colonne := 0; colonne < 7; colonne++ {
			classe := "case-vide"  // Classe par défaut
			
			if jeuActuel.Plateau[ligne][colonne] == "rouge" {
				classe = "pion-rouge"
			} else if jeuActuel.Plateau[ligne][colonne] == "jaune" {
				classe = "pion-jaune"
			}
			
			html += fmt.Sprintf(`<td id="c%dr%d" class="%s"></td>`, 
				colonne+1, ligne+1, classe)
		}
		html += "</tr>"
	}
	
	html += `</table>
    
    <div class="actions">
        <a href="/nouveau-jeu"><button>Nouveau Jeu</button></a>
        <a href="/"><button>Retour Accueil</button></a>
    </div>
</body>
</html>`
	
	return html
}

// SetupRoutes configure toutes les routes du serveur web
// Cette fonction relie tous les handlers entre eux
func SetupRoutes() {
	// Route principale "/" → Page d'accueil avec le bouton
	http.HandleFunc("/", HomeHandler)
	
	// Route "/welcome" → Gère le clic sur le bouton (redirection)
	http.HandleFunc("/welcome", WelcomeHandler)
	
	// Route "/next" → Page de destination après redirection
	http.HandleFunc("/next", NextHandler)
	
	// NOUVELLES ROUTES POUR LE JEU
	http.HandleFunc("/jeu", GameHandler)
	// Route pour afficher la page de jeu
	
	http.HandleFunc("/jouer/", JouerHandler)  
	// Route pour gérer les coups (toutes les URLs /jouer/1, /jouer/2, etc.)
	
	http.HandleFunc("/nouveau-jeu", NouveauJeuHandler)
	// Route pour recommencer une partie
	
	// Route pour servir les fichiers CSS et autres assets
	// FileServer permet de servir tous les fichiers du dossier assets/
	http.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("assets/"))))
}

// StartServer démarre le serveur web sur le port 8080
// Cette fonction lance le serveur et affiche un message de confirmation
func StartServer() {
	fmt.Println("Serveur démarré sur le port 8080...")
	fmt.Println("Ouvrez votre navigateur sur : http://localhost:8080")
	
	// Démarre le serveur web (bloque le programme jusqu'à arrêt)
	http.ListenAndServe(":8080", nil)
}