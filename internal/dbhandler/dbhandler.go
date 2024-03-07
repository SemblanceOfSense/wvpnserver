package dbhandler

import (
	crand "crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/json"
	"errors"
	"hash"
	"io"
	"os"
	"strconv"
	"vpnserver/internal/requesthandler"
)

func AddPublicKey(body requesthandler.PublicKeyRequestStruct) error {
    var file *os.File
    if _, err := os.Stat("/home/semblanceofsense/auth/pubkeys/" + strconv.Itoa(body.UserID)); errors.Is(err, os.ErrNotExist) {
        file, err = os.Create("/home/semblanceofsense/auth/pubkeys/" + strconv.Itoa(body.UserID))
        if err != nil {
            return err
        }
    } else {
        return errors.New("PubKey Already Added")
    }

    marshal, error := json.Marshal(body)
    if error != nil {return error}

    _, error = file.WriteString(string(marshal))
    if error != nil {return error}

    return nil
}

func AddPrivKey(body requesthandler.PrivateKeyRequestStruct) error {
    var file *os.File
        if _, err := os.Stat("/home/semblanceofsense/auth/privkeys/" + strconv.Itoa(body.UserID)); errors.Is(err, os.ErrNotExist) {
            file, err = os.Create("/home/semblanceofsense/auth/privkeys/" + strconv.Itoa(body.UserID))
            if err != nil {
                return err
            }
        } else {
            return errors.New("PubKey Already Added")
        }

        marshal, error := json.Marshal(body)
        if error != nil {return error}

        _, error = file.WriteString(string(marshal))
        if error != nil {return error}

        return nil
}

type EncryptedWgKey struct {
    sha hash.Hash
    rand io.Reader
    cipherText []byte
}

func EncryptKey(id int, key string) ([]byte, error) {
    j, err := os.ReadFile("/home/semblanceofsense/auth/pubkeys/" + strconv.Itoa(id))
    if err != nil { return make([]byte, 0), err }

    publicStruct := &rsa.PublicKey{}
    json.Unmarshal(j, publicStruct)

    sha := sha256.New()
    rand := crand.Reader

    encryptedBytes, err := rsa.EncryptOAEP(
	    sha,
        rand,
        publicStruct,
	    []byte(key),
	    nil,
    )
    if err != nil {
        return make([]byte, 0), err
    }

    return encryptedBytes, nil
}
