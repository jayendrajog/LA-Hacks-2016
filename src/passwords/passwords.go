package passwords

import (
	"db"
	"errors"
	"log"
)

var idIps = make(map[uint]string)

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

	idIps[userID] = remoteAddr
	log.Printf("%d has ip:%s\n", userID, remoteAddr)
	return true, nil
}

func GetCreds(userID uint, domain, remoteAddr string) ([2]string, error) {

	var ret [2]string

	if ref_remoteAddr, ok := idIps[userID]; !ok || ref_remoteAddr != remoteAddr {
		return ret, errors.New("ID does not match")
	}

	var username string
	var password string

	err := db.Db.QueryRow("SELECT username, password FROM credentials WHERE userID=? AND domain=?", userID, domain).Scan(&username, &password)

	if err != nil {
		return ret, err
	}

	return [2]string{username, password}, nil
}
