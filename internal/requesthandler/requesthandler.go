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
    Id int
    Pubkey string
    Signature []byte
}

type AddPeerRequestStruct struct {
    Publickey string
}

func AddServerPeer(thingStruct AddServerPeerStruct) error {
    msg := []byte(thingStruct.Pubkey)

    msgHash := sha256.New()
    _, err := msgHash.Write(msg)
    if err != nil {
	    return err
    }
    msgHashSum := msgHash.Sum(nil)

    fmt.Println(thingStruct.Id)
    j, err := os.ReadFile("/home/semblanceofsense/auth/pubkeys/" + strconv.Itoa(thingStruct.Id))
    if err != nil { return err }
    fmt.Println(string(j))

    publicStruct := &PublicKeyRequestStruct{}
    err = json.Unmarshal(j, publicStruct)
    if err != nil { return err }

    err = rsa.VerifyPSS(&publicStruct.Publickey, crypto.SHA256, msgHashSum, thingStruct.Signature, nil)
    if err != nil {
        return err
    }

    requestBody := &AddPeerRequestStruct{Publickey: thingStruct.Pubkey}

    jsonData, err := json.Marshal(requestBody)
    if err != nil {
        panic(err)
    }

    req, err := http.NewRequest("POST", "http://140.82.19.210:8080/addpeer", bytes.NewReader(jsonData))
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
