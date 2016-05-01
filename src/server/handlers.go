package server

import (
	"face_auth"
	"fmt"
	"github.com/gorilla/mux"
	"io"
	"log"
	"net/http"
	"os"
	"passwords"
	"session"
	"strconv"
	"time"
	"ws"
)

func makeSession(w http.ResponseWriter, r *http.Request) {
	s, err := session.New()

	if err != nil {
		WriteError(w, err, 500)
		return
	}

	WriteJson(w, map[string]interface{}{"SessionID": s.SessionID})
}

func getNextPicture(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	sessionID, err := strconv.Atoi(vars["ID"])
	if err != nil {
		WriteErrorString(w, "Invalid ID", 400)
		return
	}

	s, ok := session.Sessions[uint32(sessionID)]
	if !ok {
		WriteErrorString(w, "Session ID not found", 404)
		return
	}

	url, err := s.NextPictureUrl()
	if err != nil {
		WriteError(w, err, 500)
		return
	}
	ws.BroadcastEvent(uint32(sessionID), "Update", s.CurrentPID)

	WriteJson(w, map[string]interface{}{"URL": url})
}

func updateReaction(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	sessionID, err := strconv.Atoi(vars["ID"])
	if err != nil {
		WriteErrorString(w, "Invalid ID", 400)
		return
	}

	value, err := strconv.ParseFloat(r.FormValue("value"), 64)
	if err != nil {
		WriteErrorString(w, "Invalid value", 400)
		return
	}

	s, ok := session.Sessions[uint32(sessionID)]
	if !ok {
		WriteErrorString(w, "Session ID not found", 404)
		return
	}

	err = s.UpdateReaction(value)
	if err != nil {
		WriteError(w, err, 500)
		return
	}

	w.WriteHeader(200)
}

// upload logic
func upload(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(32 << 30)
	if err != nil {
		WriteError(w, err, 400)
		return
	}
	file, handler, err := r.FormFile("file")
	if err != nil {
		WriteError(w, err, 500)
		return
	}
	defer file.Close()
	f, err := os.OpenFile("./faces/"+handler.Filename, os.O_WRONLY|os.O_CREATE, 0664)
	if err != nil {
		WriteError(w, err, 500)
		return
	}
	defer f.Close()
	io.Copy(f, file)
	fmt.Fprintf(w, "%v", handler.Header)
}

// // upload logic
func newUser(w http.ResponseWriter, r *http.Request) {

	name := r.FormValue("name")
	if name == "" {
		WriteErrorString(w, "name not in query values", 400)
		return
	}

	passwordString := r.FormValue("password")
	if passwordString == "" {
		passwordString = "156"
	}

	passwordint, err := strconv.Atoi(passwordString)
	if err != nil {
		WriteError(w, err, 400)
		return
	}
	password := uint(passwordint)

	err = r.ParseMultipartForm(32 << 20)
	if err != nil {
		WriteError(w, err, 400)
		return
	}
	file, handler, err := r.FormFile("file")
	if err != nil {
		WriteError(w, err, 500)
		return
	}
	defer file.Close()
	log.Printf("%v", handler.Header)
	filename := fmt.Sprintf("%d.jpg", time.Now().Unix())
	filepath := "./faces/" + filename

	f, err := os.OpenFile(filepath, os.O_WRONLY|os.O_CREATE, 0664)
	if err != nil {
		WriteError(w, err, 500)
		return
	}
	defer f.Close()
	io.Copy(f, file)

	fmt.Fprintf(w, "%v", handler.Header)

	userID, err := face_auth.NewUser(name, password, filename)
	if err != nil {
		WriteError(w, err, 500)
		return
	}

	WriteJson(w, map[string]interface{}{"UserID": userID})

}

// upload logic
func checkFace(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(32 << 20)
	if err != nil {
		WriteError(w, err, 400)
		return
	}
	file, handler, err := r.FormFile("file")
	if err != nil {
		WriteError(w, err, 500)
		return
	}
	defer file.Close()
	log.Printf("%v", handler.Header)
	filename := fmt.Sprintf("%d.jpg", time.Now().Unix())
	f, err := os.OpenFile("./tempFaces/"+filename, os.O_WRONLY|os.O_CREATE, 0664)
	if err != nil {
		WriteError(w, err, 500)
		return
	}
	defer f.Close()
	io.Copy(f, file)
	id, name, err := face_auth.CheckFace(filename)
	if err != nil {
		WriteError(w, err, 500)
		return
	}
	// err = os.Remove("./tempFaces/" + filename)
	if err != nil {
		WriteError(w, err, 500)
		return
	}

	WriteJson(w, map[string]interface{}{"Match": name, "Name": name, "UserID": id})
}

func getCreds(w http.ResponseWriter, r *http.Request) {
	domain := r.FormValue("domain")
	if domain == "" {
		WriteErrorString(w, "domain not in query values", 400)
		return
	}

	ip := r.Header.Get("X-Real-IP")
	creds, err := passwords.GetCreds(domain, ip)
	if err != nil {
		WriteError(w, err, 400)
		return
	}

	WriteJson(w, map[string]interface{}{"Username": creds[0], "Password": creds[1]})
}

func makeCreds(w http.ResponseWriter, r *http.Request) {
	domain := r.FormValue("domain")
	if domain == "" {
		WriteErrorString(w, "domain not in query values", 400)
		return
	}

	username := r.FormValue("username")
	if username == "" {
		WriteErrorString(w, "username not in query values", 400)
		return
	}

	password := r.FormValue("password")
	if password == "" {
		WriteErrorString(w, "password not in query values", 400)
		return
	}

	ip := r.Header.Get("X-Real-IP")
	err := passwords.MakeCreds(domain, username, password, ip)
	if err != nil {
		WriteError(w, err, 500)
		return
	}

	w.WriteHeader(200)
}

func check_password(w http.ResponseWriter, r *http.Request) {
	userIDint, err := strconv.Atoi(r.FormValue("userid"))
	if err != nil {
		WriteError(w, err, 400)
		return
	}
	userID := uint(userIDint)

	passwordint, err := strconv.Atoi(r.FormValue("password"))
	if err != nil {
		WriteError(w, err, 400)
		return
	}
	password := uint(passwordint)

	ip := r.Header.Get("X-Real-IP")

	verified, err := passwords.VerifyPassword(userID, password, ip)
	if err != nil {
		WriteError(w, err, 400)
		return
	}

	if !verified {
		WriteErrorString(w, "Not verified", 400)
		return
	}

	w.WriteHeader(200)
}
