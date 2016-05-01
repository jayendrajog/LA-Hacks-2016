package server

import (
	"db"
	"face_auth"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"log"
	"math/rand"
	"net/http"
	"time"
	"ws"
)

func Log(handler http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s %s", r.RemoteAddr, r.Method, r.URL)
		handler.ServeHTTP(w, r)
	})
}

func Run(port uint16) {
	//start := time.Now()

	rand.Seed(time.Now().Unix())

	face_auth.Init()

	db.Open()
	defer db.Close()

	//log.Println("Took %s", time.Now().Sub(start))
	//log.Println(post)
	r := mux.NewRouter()
	r.HandleFunc("/session", Log(makeSession)).Methods("POST")
	r.HandleFunc("/sessions", Log(makeSession)).Methods("POST")
	r.HandleFunc("/sessions/{ID:[0-9]+}/nextPicture", Log(getNextPicture)).Methods("GET")
	r.HandleFunc("/sessions/{ID:[0-9]+}/reaction", Log(updateReaction)).Methods("POST")
	r.HandleFunc("/sessions/{ID:[0-9]+}/ws", Log(ws.ServeWs)).Methods("GET")

	r.HandleFunc("/upload", Log(upload)).Methods("POST")
	r.HandleFunc("/photo", Log(checkFace)).Methods("POST")
	r.HandleFunc("/photos", Log(checkFace)).Methods("POST")

	r.HandleFunc("/creds", Log(getCreds)).Methods("GET")

	r.HandleFunc("/myo_password", Log(check_password)).Methods("POST")

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "OPTIONS"},
		AllowCredentials: true,
	})
	handler := c.Handler(r)

	// r.HandleFunc("/shortlink/{linkID}", deleteShortlink).Methods("DELETE")

	for {
		log.Printf("Running at 0.0.0.0:%d\n", port)
		log.Println(http.ListenAndServe(fmt.Sprintf("0.0.0.0:%d", port), handler))
		time.Sleep(1 * time.Second)
	}
}
