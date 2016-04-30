package server

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"ipCircBuffer"
	"net"
	"net/http"
	"nginxParser"
	"posts"
	"shortlink"
	"strconv"
	"strings"
	"time"
)

func getPost(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	post, err := posts.GetPost(vars["Title"])
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			WriteError(w, err, 404)
			return
		}
	}
	WriteOutputError(w, map[string]*posts.Post{"Post": post}, err)
}

func getPostList(w http.ResponseWriter, r *http.Request) {
	WriteJson(w, map[string][]string{"Posts": posts.GetPostList()})
}

func getParagraph(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	post, err := posts.GetPost(vars["Title"])
	if err != nil {
		WriteError(w, err, 400)
		return
	}
	i, err := strconv.Atoi(vars["id"])
	if err != nil {
		WriteError(w, err, 400)
		return
	}

	if i >= len(post.Content) {
		WriteErrorString(w, fmt.Sprintf("%s does not have paragraph %d", vars["Title"], i), 400)
		return
	}

	paragraph := post.Content[i]
	WriteJson(w, map[string]string{"Paragraph": paragraph})
}

func getInfo(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	post, err := posts.GetPost(vars["Title"])
	if err != nil {
		WriteError(w, err, 400)
		return
	}
	info := (*post).Info()
	WriteJson(w, map[string]posts.PostInfo{"Info": info})
}

func getDesktopIP(w http.ResponseWriter, r *http.Request) {

	ip := ipCircBuffer.DesktopIP.Get()
	if ip == nil {
		ip = net.ParseIP("0.0.0.0")
	}
	fmt.Fprint(w, ip.String())
}

func postDesktopIP(w http.ResponseWriter, r *http.Request) {
	stringIP := r.FormValue("IP")
	if stringIP == "-1" {
		ipCircBuffer.DesktopIP.Revert()
		fmt.Fprint(w, "Reverted")
		return
	}
	ipCircBuffer.DesktopIP.Set(stringIP)
	fmt.Fprint(w, stringIP)
}

func clearDesktopIP(w http.ResponseWriter, r *http.Request) {
	ipCircBuffer.RaspberryIP.Clear()
	fmt.Fprint(w, "Cleared")
}

func getRaspberryIP(w http.ResponseWriter, r *http.Request) {

	ip := ipCircBuffer.RaspberryIP.Get()
	if ip == nil {
		ip = net.ParseIP("0.0.0.0")
	}
	fmt.Fprint(w, ip.String())
}

func postRaspberryIP(w http.ResponseWriter, r *http.Request) {
	stringIP := r.FormValue("IP")
	if stringIP == "-1" {
		ipCircBuffer.RaspberryIP.Revert()
		fmt.Fprint(w, "Reverted")
		return
	}
	ipCircBuffer.RaspberryIP.Set(stringIP)
	fmt.Fprint(w, stringIP)
}

func clearRaspberryIP(w http.ResponseWriter, r *http.Request) {
	ipCircBuffer.RaspberryIP.Clear()
	fmt.Fprint(w, "Cleared")
}

func get24HourRequests(w http.ResponseWriter, r *http.Request) {
	requests, err := nginxParser.GetRequests(24*time.Hour, "/var/log/nginx/access.log")
	WriteOutputError(w, map[string]int{"Requests": requests}, err)
}
func redirectToShortlink(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/#shortlink", 302)
}

func makeShortlink(w http.ResponseWriter, r *http.Request) {

	parsedJson := make(map[string]string)
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&parsedJson)
	if err != nil {
		WriteErrorString(w, err.Error()+" in parsing POST body (JSON)", 400)
		return
	}

	link, ok := parsedJson["Link"]
	if !ok {
		WriteErrorString(w, "'Link' not in JSON body", 400)
		return
	}

	newID, err := shortlink.Add(link)
	if err != nil {
		WriteError(w, err, 400)
		return
	}
	WriteJson(w, map[string]string{"ID": newID})
}

func getShortlink(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	linkID := vars["linkID"]
	if linkID == "" {
		http.Redirect(w, r, "/beta#shortlink", 302)
		//WriteErrorString(w, "link ID empty", 400)
		return
	}

	//pads front with 0
	for len(linkID) < 4 {
		linkID = "0" + linkID
	}
	link, err := shortlink.Get(linkID)
	if err != nil {
		WriteError(w, err, 400)
		return
	}
	http.Redirect(w, r, link, 302)

}
