package requesthandler

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
    "crypto/rsa"
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
