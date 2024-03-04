package main

import (
	"fmt"
	"log"
	"net/http"
	"vpnserver/internal/requesthandler"
    "vpnserver/internal/dbhandler"
)

func (h HelloHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    switch r.URL.Path {
    case "/addpublickey":
        switch r.Method {
        case http.MethodGet:
            fmt.Fprintf(w, "hi")
        case http.MethodPost:
            err := dbhandler.AddPublicKey(requesthandler.PublicKeyRequest(r.Body))
            if err != nil {log.Fatal(err)}
        default:
            http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        }
    case "/addprivatekey":
        switch r.Method {
        case http.MethodPost:
            err := dbhandler.AddPrivKey(requesthandler.PrivateKeyRequest(r.Body))
            if err != nil {log.Fatal(err)}
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
