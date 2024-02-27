package main

import (
	"fmt"
	"log"
	"net/http"
	"vpnserver/internal/requesthandler"
)

func (h HelloHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    switch r.URL.Path {
    case "/addpublickey":
        switch r.Method {
        case http.MethodGet:
            fmt.Fprintf(w, "hi")
        case http.MethodPost:
            fmt.Print(requesthandler.PublicKeyRequest(r.Body).Publickey.E)
        default:
            http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        }
    case "/addprivatekey":
        switch r.Method {
        case http.MethodPost:
            fmt.Print(string(requesthandler.PrivateKeyRequest(r.Body).Ciphertext))
        default:
            http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        }
    default:
        http.NotFound(w, r)
        return
    }
}

type HelloHandler struct {}

func main() {
    http.Handle("/", &HelloHandler{})

    log.Fatal(http.ListenAndServe(":8080", nil))
}
