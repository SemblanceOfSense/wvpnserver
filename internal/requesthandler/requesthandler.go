package requesthandler

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
    "crypto/rsa"
    "net/http"
)

type PublicKeyRequestStruct struct {
    Publickey rsa.PublicKey
    UserID int
}

type PrivateKeyRequestStruct struct {
    Ciphertext []byte
    Iv []byte
    Salt []byte
    UserID int
}

func PublicKeyRequest(v io.ReadCloser) PublicKeyRequestStruct {
    dec := json.NewDecoder(v)

    var request PublicKeyRequestStruct
    err := dec.Decode(&request)
    if err != nil {
        fmt.Println("Failed decode")
        log.Fatal(err)
    }

    return request
}

func PrivateKeyRequest(v io.ReadCloser) PrivateKeyRequestStruct {
    dec := json.NewDecoder(v)

    var request PrivateKeyRequestStruct
    err := dec.Decode(&request)
    if err != nil {
        fmt.Println("Failed decode")
        log.Fatal(err)
    }

    return request
}

func GetVpnKey() (string, error) {
    req, err := http.NewRequest("GET", "http://140.82.19.210:8080/getpubkey", nil)
    if err != nil { return "", err }

    resp, err := http.DefaultClient.Do(req)
    if err != nil { return "", err }
    defer resp.Body.Close()
    body, err := io.ReadAll(resp.Body)

    return string(body), nil
}
