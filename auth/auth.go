package auth

import (
	"errors"
	"fmt"
	"os/exec"
	"time"

	"github.com/essentialkaos/branca/v2"
	"github.com/evdnx/unixmint/constants"
	"github.com/evdnx/unixmint/db"
	"github.com/evdnx/unixmint/pkg/crypt"
	"github.com/evdnx/unixmint/pkg/util"
	"github.com/goccy/go-json"
)

func Init() error {
	// try to read key
	_, err := db.Read(constants.AuthBucket, "key")
	if err != nil {
		// create new key
		k := crypt.GenerateKey(32)
		err = db.Update(constants.AuthBucket, "key", string(k))
		if err != nil {
			return err
		}
	}

	return nil
}

func Login(username, password string) error {
	command := fmt.Sprint("echo ", password, " | su - ", username, " -c ", `"echo 1"`)
	cmd := exec.Command("sh", "-c", command)
	err := cmd.Run()
	if err != nil {
		return err
	}

	// encrypt password
	encryptedPassword, err := crypt.Encrypt(password)
	if err != nil {
		return err
	}

	// write encrypted password to db
	err = db.Update(constants.AuthBucket, "password", encryptedPassword)
	if err != nil {
		return err
	}

	return nil
}

func Logout() error {
	err := db.Update(constants.AuthBucket, "password", "")
	if err != nil {
		return err
	}

	return nil
}

func Username() string {
	return ""
}

func Password() string {
	// read password
	pwd, err := db.Read(constants.AuthBucket, "password")
	if err != nil {
		return ""
	}

	// decrypt password
	password, err := crypt.Decrypt(pwd)
	if err != nil {
		return ""
	}

	return password
}

func EncodeToken(payload any, ttlHours uint32) (string, int64, error) {
	brc, err := branca.NewBranca([]byte("TODO"))
	if err != nil {
		return "", 0, err
	}

	ttlSeconds := ttlHours * 60 * 60

	payloadBytes := util.InterfaceToByte(payload)
	token, err := brc.EncodeToString(payloadBytes)
	if err != nil {
		return "", 0, err
	}

	millis := ttlSeconds * 1000
	expiresAt := time.Now().UnixMilli() + int64(millis)
	return token, expiresAt, nil
}

func DecodeToken(token string, ttlHours uint32, data any) (rawPayload []byte, e error) {
	brc, err := branca.NewBranca([]byte("TODO"))
	if err != nil {
		return nil, err
	}

	brancaToken, err := brc.DecodeString(token)
	if err != nil {
		return nil, err
	}

	ttlSeconds := ttlHours * 60 * 60

	if brancaToken.IsExpired(ttlSeconds) {
		return nil, errors.New("auth token is expired")
	}

	if data != nil {
		err = json.Unmarshal(brancaToken.Payload(), data)
		if err != nil {
			return nil, err
		}
	}

	return brancaToken.Payload(), nil
}
