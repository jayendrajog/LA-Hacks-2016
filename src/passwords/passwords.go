package passwords

import (
	"db"
	"encoding/json"
	"encrypt"
	"errors"
	"io/ioutil"
	"log"
)

var ipIDs = make(map[string]uint)
var key string

func Init() error {

	keys := make(map[string]string)
	content, err := ioutil.ReadFile("./creds/key.json")
	if err != nil {
		return err
	}
	err = json.Unmarshal(content, &keys)
	if err != nil {
		return err
	}

	key = keys["key"]

	return nil
}

func GetName(userID uint) (string, error) {
	var name string

	err := db.Db.QueryRow("SELECT name FROM myo_passwords WHERE userID=?", userID).Scan(&name)
	if err != nil {
		return "", err
	}

	return name, nil
}

func VerifyPassword(userID, password uint, remoteAddr string) (bool, error) {
	var ref_pass uint

	err := db.Db.QueryRow("SELECT password FROM myo_passwords WHERE userID=?", userID).Scan(&ref_pass)
	if err != nil {
		return false, err
	}

	if ref_pass != password {
		return false, nil
	}

	ipIDs[remoteAddr] = userID
	log.Printf("%s has user:%d\n", remoteAddr, userID)
	return true, nil
}

func GetCreds(domain, remoteAddr string) ([2]string, error) {

	var ret [2]string

	var userID uint
	var ok bool

	if userID, ok = ipIDs[remoteAddr]; !ok {
		return ret, errors.New("Not logged in")
	}

	var username string
	var encrypted_password string

	err := db.Db.QueryRow("SELECT username, password FROM credentials WHERE userID=? AND domain=?", userID, domain).Scan(&username, &encrypted_password)

	password, err := encrypt.Decrypt(key, encrypted_password)
	if err != nil {
		return ret, err
	}

	return [2]string{username, password}, nil
}

func MakeCreds(domain, username, password, remoteAddr string) error {

	var userID uint
	var ok bool

	if userID, ok = ipIDs[remoteAddr]; !ok {
		return errors.New("Not logged in")
	}

	encrypted_password, err := encrypt.Encrypt(key, password)
	if err != nil {
		return err
	}

	newCred, err := db.Db.Prepare("INSERT INTO credentials VALUES(?, ?, ?, ?)")
	if err != nil {
		return err
	}
	_, err = newCred.Exec(userID, domain, username, encrypted_password)
	if err != nil {
		return err
	}

	return nil
}
