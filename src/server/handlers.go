package server

import (
	"github.com/gorilla/mux"
	"net/http"
	"session"
	"strconv"
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
