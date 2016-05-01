package server

import (
	"face_auth"
	"fmt"
	"github.com/gorilla/mux"
	"io"
	"log"
	"net/http"
	"os"
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
	f, err := os.OpenFile("./tempFaces/"+handler.Filename, os.O_WRONLY|os.O_CREATE, 0664)
	if err != nil {
		WriteError(w, err, 500)
		return
	}
	defer f.Close()
	io.Copy(f, file)
	fmt.Fprintf(w, "%v", handler.Header)
}

// upload logic
func checkFace(w http.ResponseWriter, r *http.Request) {
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
	log.Printf("%v", handler.Header)
	filename := fmt.Sprintf("%d.jpg", time.Now().Unix())
	f, err := os.OpenFile("./tempFaces/"+filename, os.O_WRONLY|os.O_CREATE, 0664)
	if err != nil {
		WriteError(w, err, 500)
		return
	}
	defer f.Close()
	io.Copy(f, file)
	name, err := face_auth.CheckFace(filename)
	if err != nil {
		WriteError(w, err, 500)
		return
	}
	// err = os.Remove("./tempFaces/" + filename)
	if err != nil {
		WriteError(w, err, 500)
		return
	}

	WriteJson(w, map[string]interface{}{"Match": name})
}

func getCreds(w http.ResponseWriter, r *http.Request) {
	// website := r.FormValue("website")
	// userID := r.FormValue("userid")
	// secret := r.FormValue("secret")

	WriteJson(w, map[string]interface{}{"Username": "marktai", "Password": "thisistotallysecure"})
}
