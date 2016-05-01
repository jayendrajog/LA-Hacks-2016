package passwords

import (
	"db"
	"encrypt"
	"errors"
	"log"
)

var ipIDs = make(map[string]uint)

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

	password, err := encrypt.Decrypt("0ELBGZt6AZf9U6Qc6SteS3tPJ9lpeTFf", encrypted_password)
	if err != nil {
		return ret, err
	}

	if err != nil {
		return ret, err
	}

	return [2]string{username, password}, nil
}
