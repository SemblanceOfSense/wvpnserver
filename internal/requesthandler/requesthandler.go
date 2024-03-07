package requesthandler

import (
	"bytes"
	"crypto"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
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

type AddServerPeerStruct struct {
    id int
    pubkey string
    signature []byte
}

func AddServerPeer(thingStruct AddServerPeerStruct) error {
    msg := []byte(thingStruct.pubkey)

    msgHash := sha256.New()
    _, err := msgHash.Write(msg)
    if err != nil {
	    return err
    }
    msgHashSum := msgHash.Sum(nil)

    fmt.Println(thingStruct.id)
    j, err := os.ReadFile("/home/semblanceofsense/auth/pubkeys/" + strconv.Itoa(thingStruct.id))
    if err != nil { return err }

    publicStruct := &PublicKeyRequestStruct{}
    fmt.Println(string(j))
    json.Unmarshal(j, &thingStruct)

    err = rsa.VerifyPSS(&publicStruct.Publickey, crypto.SHA256, msgHashSum, thingStruct.signature, nil)
    if err != nil {
        return err
    }

    req, err := http.NewRequest("POST", "https://140.82.19.210:8080/addpeer", bytes.NewReader([]byte(msg)))
    if err != nil {
        return err
    }

    client := http.Client{Timeout: 10 * time.Second}
    res, err := client.Do(req)
    if err != nil {
        return err
    }
    log.Printf("status Code: %d", res.StatusCode)

    return nil
}
