package auth

import (
	"encoding/json"
	"errors"
	"fmt"
	"os/exec"
	"time"

	"github.com/essentialkaos/branca/v2"
	"github.com/evdnx/unixmint/db"
	"github.com/evdnx/unixmint/pkg/crypt"
	"github.com/evdnx/unixmint/pkg/util"
)

var Timer *time.Timer

func Init() error {
	// try to read key
	_, err := db.Read(db.AuthBucket, "crypto_key")
	if err != nil {
		// create new key
		k := crypt.GenerateKey(32)
		err = db.Update(db.AuthBucket, "crypto_key", string(k))
		if err != nil {
			return err
		}
	}

	// branca key
	_, err = db.Read(db.AuthBucket, "branca_key")
	if err != nil {
		// create new key
		k := crypt.GenerateKey(32)
		err = db.Update(db.AuthBucket, "branca_key", string(k))
		if err != nil {
			return err
		}
	}

	return nil
}

func BrancaKey() []byte {
	key, err := db.Read(db.AuthBucket, "branca_key")
	if err != nil {
		return nil
	}

	return []byte(key)
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
	err = db.Update(db.AuthBucket, "password", encryptedPassword)
	if err != nil {
		return err
	}

	// set last login time
	err = db.Update(db.AuthBucket, "last_login", time.Now().UTC().Format(time.RFC3339))
	if err != nil {
		return err
	}

	// logout after 24 hours automatically
	Timer = time.AfterFunc(24*time.Hour, func() { Logout() })

	return nil
}

func Logout() error {
	Timer.Stop()
	Timer = nil

	err := db.Update(db.AuthBucket, "password", "")
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
	pwd, err := db.Read(db.AuthBucket, "password")
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
	brc, err := branca.NewBranca(BrancaKey())
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
	brc, err := branca.NewBranca(BrancaKey())
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
