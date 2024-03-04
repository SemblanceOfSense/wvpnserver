package dbhandler

import (
	"encoding/json"
    "os"
	"vpnserver/internal/requesthandler"

    "errors"
    "strconv"
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
            file, err = os.Create("/home/semblanceofsense/auth/privkeys" + strconv.Itoa(body.UserID))
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
