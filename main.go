package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync/atomic"

	"github.com/Chin-mayyy/Chirpy/internal/database"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {
	//Getting the evironment variables
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}
	dbURL := os.Getenv("DB_URL")
	pf := os.Getenv("PLATFORM")

	secret := os.Getenv("JWTSECRET")

	//Setting up a connection to the database.
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal(err)
	}

	dbQueries := database.New(db)

	//Creating a new HTTP request multiplexer which allocates URL to the most appropriate handler.
	mux := http.NewServeMux()
	const port = "8080"
	const filepathRoot = "."

	apicfg := apiConfig{
		fileserverHits: atomic.Int32{},
		db:             dbQueries,
		platform:       pf,
		JWTsecret:      secret,
	}

	//Serves file from assets directory.
	fs := http.FileServer(http.Dir("."))

	mux.Handle("/app/", apicfg.middlewareMetricsInc(http.StripPrefix("/app", fs)))
	mux.HandleFunc("GET /admin/metrics", apicfg.handlerMetrics)
	mux.HandleFunc("POST /admin/reset", apicfg.handlerReset)
	mux.HandleFunc("GET /api/healthz", handlerReadiness)
	mux.HandleFunc("POST /api/chirps", apicfg.handlerCreateChirp)
	mux.HandleFunc("POST /api/users", apicfg.handlerAcceptEmail)
	mux.HandleFunc("POST /api/login", apicfg.handlerGetUser)
	mux.HandleFunc("POST /api/revoke", apicfg.handlerRevokes)
	mux.HandleFunc("POST /api/refresh", apicfg.handlerRefreshes)
	mux.HandleFunc("GET /api/chirps", apicfg.handlerGetChirps)
	mux.HandleFunc("GET /api/chirps/{id}", apicfg.handlerGetChirp)

	fmt.Printf("Serving files from %s on port: %s\n", filepathRoot, port)

	svr := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	//Starting the server.
	log.Fatal(svr.ListenAndServe())
}

type apiConfig struct {
	fileserverHits atomic.Int32
	db             *database.Queries
	platform       string
	JWTsecret      string
}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Middleware triggered for:", r.URL.Path)
		cfg.fileserverHits.Add(1)
		next.ServeHTTP(w, r)
	})
}

func (cfg *apiConfig) handlerMetrics(w http.ResponseWriter, r *http.Request) {
	html := fmt.Sprintf(`<html>
  		<body>
    			<h1>Welcome, Chirpy Admin</h1>
      			<p>Chirpy has been visited %d times!</p>
        	</body>
        </html>`, cfg.fileserverHits.Load())

	w.Header().Add("Content-Type", "text/html")

	w.WriteHeader(200)

	w.Write([]byte(html))
}
