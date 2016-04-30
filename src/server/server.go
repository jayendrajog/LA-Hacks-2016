package server

import (
	"db"
	"fmt"
	"github.com/gorilla/mux"
	"ipCircBuffer"
	"log"
	"net/http"
	"posts"
	"time"
)

func Log(handler http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s %s", r.RemoteAddr, r.Method, r.URL)
		handler.ServeHTTP(w, r)
	})
}

func Run(port uint16) {
	//start := time.Now()
	err := posts.Init("posts")

	ipCircBuffer.Init()

	db.Open()
	defer db.Close()

	if err != nil {
		log.Println(err)
	}
	//log.Println("Took %s", time.Now().Sub(start))
	//log.Println(post)
	r := mux.NewRouter()
	r.HandleFunc("/posts", Log(getPostList))
	r.HandleFunc("/posts/{Title}", Log(getPost))
	r.HandleFunc("/posts/{Title}/paragraph/{id:[0-9]+}", Log(getParagraph)).Methods("GET")
	r.HandleFunc("/posts/{Title}/info", Log(getInfo)).Methods("GET")
	r.HandleFunc("/desktopIP", Log(getDesktopIP)).Methods("GET")
	r.HandleFunc("/desktopIP", Log(postDesktopIP)).Methods("POST")
	r.HandleFunc("/desktopIP", Log(clearDesktopIP)).Methods("DELETE")
	r.HandleFunc("/raspberryIP", Log(getRaspberryIP)).Methods("GET")
	r.HandleFunc("/raspberryIP", Log(postRaspberryIP)).Methods("POST")
	r.HandleFunc("/raspberryIP", Log(clearRaspberryIP)).Methods("DELETE")
	r.HandleFunc("/requestCount", Log(get24HourRequests)).Methods("GET")
	r.HandleFunc("/shortlink", Log(redirectToShortlink)).Methods("GET")
	r.HandleFunc("/shortlink/", Log(redirectToShortlink)).Methods("GET")
	r.HandleFunc("/shortlink", Log(makeShortlink)).Methods("POST")
	r.HandleFunc("/shortlink/{linkID}", Log(getShortlink)).Methods("GET")
	// r.HandleFunc("/shortlink/{linkID}", deleteShortlink).Methods("DELETE")

	for {
		log.Printf("Running at 0.0.0.0:%d\n", port)
		log.Println(http.ListenAndServe(fmt.Sprintf("0.0.0.0:%d", port), r))
		time.Sleep(1 * time.Second)
	}
}
