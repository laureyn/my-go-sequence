package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	_ "github.com/lib/pq"
)

var db *sql.DB

func handler(w http.ResponseWriter, r *http.Request) {
	var seqValue int
	err := db.QueryRow("SELECT nextval('my_sequence')").Scan(&seqValue)
	if err != nil {
		log.Printf("❌ Erreur SQL: %v", err)
		http.Error(w, "DB Error: "+err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Fprintf(w, "%d", seqValue)
}

func main() {
	// Récupération des variables d’environnement
	dbUser := getEnv("DB_USER", "postgres")
	dbPass := getEnv("DB_PASS", "")
	dbHost := getEnv("DB_HOST", "localhost")
	dbPort := getEnv("DB_PORT", "5432")
	dbName := getEnv("DB_NAME", "postgres")

	// Construction du DSN PostgreSQL
	dsn := fmt.Sprintf(
		"user=%s password=%s host=%s port=%s dbname=%s sslmode=disable",
		dbUser, dbPass, dbHost, dbPort, dbName,
	)

	// Connexion avec logs
	log.Printf("📡 Connexion à PostgreSQL : %s@%s:%s/%s", dbUser, dbHost, dbPort, dbName)
	var err error
	db, err = sql.Open("postgres", dsn)
	for i := 0; i < 10; i++ {
		log.Printf("🔄 Tentative de connexion à la DB (essai %d)...", i+1)
		err = db.Ping()
		if err == nil {
			log.Println("✅ Connexion à la DB réussie")
			break
		}
		log.Printf("⏳ DB pas encore prête : %v", err)
		time.Sleep(3 * time.Second)
	}

	if err != nil {
		log.Fatalf("❌ Échec définitif de connexion à la DB après plusieurs tentatives : %v", err)
	}
	// Timeout de test de connexion
	db.SetConnMaxLifetime(time.Minute * 3)
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)

	// Vérification de la connexion
	err = db.Ping()
	if err != nil {
		log.Fatalf("❌ Échec de ping vers la base : %v", err)
	}
	log.Println("✅ Connexion réussie à PostgreSQL")

	// Serveur HTTP
	http.HandleFunc("/", handler)
	log.Println("🚀 Serveur en écoute sur :80")
	log.Fatal(http.ListenAndServe(":80", nil))
}

func getEnv(key, fallback string) string {
	if v, ok := os.LookupEnv(key); ok {
		return v
	}
	return fallback
}
