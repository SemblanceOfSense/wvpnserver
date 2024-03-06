package dbhandler

import (
	"crypto/rsa"
	"encoding/json"
	"os"
	"vpnserver/internal/requesthandler"

	"errors"
	"strconv"
    "crypto/sha256"
    crand "crypto/rand"
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

func EncryptKey(id int, key string) ([]byte, error) {
    j, err := os.ReadFile("/home/semblanceofsense/auth/privkeys/" + strconv.Itoa(id))
    if err != nil { return make([]byte, 0), err }

    publicStruct := &rsa.PublicKey{}
    json.Unmarshal(j, publicStruct)

    encryptedBytes, err := rsa.EncryptOAEP(
	    sha256.New(),
        nil,
	    publicStruct,
	    []byte(key),
	nil)
    if err != nil {
        return make([]byte, 0), err
    }

    return encryptedBytes, nil
}
