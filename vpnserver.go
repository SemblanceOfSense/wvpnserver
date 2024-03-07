package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"vpnserver/internal/dbhandler"
	"vpnserver/internal/requesthandler"
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
    case "/getpublickey":
        switch r.Method {
        case http.MethodGet:
            j, err := os.ReadFile("/home/semblanceofsense/auth/pubkeys/" + r.Header.Get("UserID"))
            if err != nil { fmt.Println(err) }
            w.Write(j)
        default:
            http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        }
    case "/getprivatekey":
        switch r.Method {
        case http.MethodGet:
            j, err := os.ReadFile("/home/semblanceofsense/auth/privkeys/" + r.Header.Get("UserID"))
            if err != nil { fmt.Println(err) }
            w.Write(j)
        default:
            http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        }
    case "/getwgkey":
        switch r.Method{
        case http.MethodGet:
            key, err := requesthandler.GetVpnKey()
            if err != nil { fmt.Println(err) }
            header, err := strconv.Atoi(r.Header.Get("UserID"))
            if err != nil { fmt.Println(err) }
            j, err := dbhandler.EncryptKey(header, key)
            if err != nil { fmt.Println(err) }
            w.Write(j)
        default:
            http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        }
    case "/addpeer":
        switch r.Method {
        case http.MethodPost:
            thingStruct := &requesthandler.AddServerPeerStruct{}
            j, err := io.ReadAll(r.Body)
            if err != nil { log.Fatal(err) }
            err = json.Unmarshal(j, thingStruct)
            if err != nil { log.Fatal(err) }
            requesthandler.AddServerPeer(*thingStruct)
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
