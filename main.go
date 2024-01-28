package main

import (
	"log"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/keshavsharma98/chirpy/internal/database"
)

type apiConfig struct {
	fileserverHits int
	DB             *database.DB
}

func main() {
	const filepathRoot = "."
	const port = "8080"

	db, err := database.NewDB("database.json")
	if err != nil {
		log.Fatalln("Error in setting up DB", err)
	}

	apiCfg := apiConfig{
		fileserverHits: 0,
		DB:             db,
	}

	// mux := http.NewServeMux()
	// mux.Handle("/app/", apiCfg.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(filepathRoot)))))
	// mux.HandleFunc("/healthz", handlerReadiness)
	// mux.HandleFunc("/metrics", apiCfg.handlerMetrics)
	// mux.HandleFunc("/reset", apiCfg.handlerReset)

	r := chi.NewRouter()
	fsHandler := apiCfg.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(filepathRoot))))
	r.Handle("/app", fsHandler)
	r.Handle("/app/*", fsHandler)

	// Create a new router for /api namespace
	apiRouter := chi.NewRouter()

	// Move non-website endpoints to /api namespace
	apiRouter.Get("/healthz", handlerReadiness)
	apiRouter.Get("/reset", apiCfg.handlerReset)
	apiRouter.Post("/chirps", apiCfg.handlerCreateChirps)
	apiRouter.Get("/chirps", apiCfg.handlerGetAllChirps)
	apiRouter.Get("/chirps/{id}", apiCfg.handlerGetChirpById)
	apiRouter.Post("/users", apiCfg.handlerCreateUsers)
	apiRouter.Post("/login", apiCfg.handlerLogin)
	// Mount the apiRouter under /api path in the main router
	r.Mount("/api", apiRouter)

	adminRouter := chi.NewRouter()
	adminRouter.Get("/metrics", apiCfg.handlerMetrics)
	r.Mount("/admin", adminRouter)

	corsMux := middlewareCors(r)

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: corsMux,
	}

	log.Printf("Serving files from %s on port: %s\n", filepathRoot, port)
	log.Fatal(srv.ListenAndServe())
}
