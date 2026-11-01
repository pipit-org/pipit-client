package main

import (
	_ "embed"
	"flag"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
)

//go:embed index.html
var index []byte

func main() {
	listen := flag.String("listen", ":3000", "dirección donde escucha el cliente")
	api := flag.String("api", "http://localhost:8080", "URL del balanceador de Pipit")
	flag.Parse()

	target, err := url.Parse(*api)
	if err != nil {
		log.Fatalf("URL de la API inválida: %v", err)
	}

	mux := http.NewServeMux()

	// /api/... se reenvía al balanceador. Así el navegador habla siempre con el
	// mismo origen (sin problemas de CORS) y los streams SSE pasan sin buffering.
	mux.Handle("/api/", http.StripPrefix("/api", httputil.NewSingleHostReverseProxy(target)))

	// La página del chat.
	mux.HandleFunc("GET /{$}", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Write(index)
	})

	log.Printf("cliente de Pipit en http://localhost%s (API: %s)", *listen, *api)
	log.Fatal(http.ListenAndServe(*listen, mux))
}
