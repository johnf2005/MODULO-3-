package main

import (
	"log"
	"net/http"

	"modulo3_go/internal/handlers"
	"modulo3_go/internal/storage"
)

func main() {
	store, err := storage.NewSQLiteStorage("modulo3_go.db")
	if err != nil {
		log.Fatalf("no se pudo inicializar el almacenamiento: %v", err)
	}
	defer store.Close()

	server := handlers.NewServer(store)
	log.Println("Servidor iniciado en http://localhost:8085")
	if err := http.ListenAndServe(":8085", server); err != nil {
		log.Fatalf("error al iniciar el servidor: %v", err)
	}
}
