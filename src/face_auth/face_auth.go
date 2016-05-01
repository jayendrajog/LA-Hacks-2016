package face_auth

import (
	"bytes"
	"db"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"passwords"
)

const DETECT_URL = "https://api.projectoxford.ai/face/v1.0/detect?returnFaceId=true&returnFaceLandmarks=false&returnFaceAttributes=age"
const VERIFY_URL = "https://api.projectoxford.ai/face/v1.0/verify"

var FaceAPIHeaders = map[string]string{
	"Content-Type":              "application/json",
	"Ocp-Apim-Subscription-Key": "38c44ac804c44f6e97673d815163a1db",
}

var FACES_DB = map[uint]string{
	1: "http://sparck.co/faces/ADAM.jpg",
	2: "http://sparck.co/faces/JAY.jpg",
	3: "http://sparck.co/faces/ANTHONY.jpg",
}

var Faces_ids = make(map[uint]string)

var Next_id uint = 0

var client *http.Client

func Init() {
	client = &http.Client{}
	log.Println(PopulateFaces_ids())
}

func PopulateFaces_ids() error {
	rows, err := db.Db.Query("SELECT userID, url, faceID FROM photos")
	if err != nil {
		return err
	}

	var userID uint
	var url string
	var faceID string
	var count = 0
	respChan := make(chan UintString)

	for rows.Next() {
		if err := rows.Scan(&userID, &url, &faceID); err != nil {
			return err
		}
		if faceID == "" {
			count += 1
			go func(userID uint, url string, respChan chan UintString) {
				faceID, err := GetFaceID(url)
				if err != nil {
					log.Println(err)
					return
				}

				updateFaceID, err := db.Db.Prepare("UPDATE photos SET faceID=? WHERE userID=? AND url=?")
				if err != nil {
					log.Println(err)
					return
				}
				_, err = updateFaceID.Exec(faceID, userID, url)
				respChan <- UintString{userID, faceID}
			}(userID, url, respChan)
		} else {
			Faces_ids[userID] = faceID
		}

		if userID > Next_id {
			Next_id = userID
		}
	}

	Next_id += 1
	for {
		if count == 0 {
			break
		}

		face_id := <-respChan
		Faces_ids[face_id.Uint] = face_id.String
		count -= 1
	}
	return nil
}

type ErrorMessage struct {
	Code    string
	Message string
}

type FaceMap struct {
	FaceId         string
	FaceRectangle  map[string]int
	FaceAttributes map[string]float64
}

type FaceResponse struct {
	Faces []FaceMap
	Error ErrorMessage
}

func NewFaceResponse() FaceResponse {
	var ret FaceResponse
	ret.Faces = make([]FaceMap, 0)
	return ret
}

type VerifyResponse struct {
	IsIdentical bool
	Confidence  float64
	Error       ErrorMessage
}

func GetFaceID(url string) (string, error) {
	mapBody := map[string]string{"url": url}

	jsonBytes, err := json.Marshal(mapBody)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest("POST", DETECT_URL, bytes.NewBuffer(jsonBytes))
	if err != nil {
		return "", err
	}

	for key, value := range FaceAPIHeaders {
		req.Header.Set(key, value)
	}

	resp, err := client.Do(req)

	if resp != nil {
		defer resp.Body.Close()
	}

	if err != nil {
		return "", nil
	}

	faceResponse := make([]FaceMap, 0)

	decoder := json.NewDecoder(resp.Body)

	if resp.StatusCode != 200 {
		return "", errors.New("Non 200 response from face identify")
	}

	err = decoder.Decode(&faceResponse)
	if err != nil {
		return "", err
	}

	if len(faceResponse) == 0 {
		return "", errors.New("No face found")
	}

	return faceResponse[0].FaceId, nil
}

type UintString struct {
	Uint   uint
	String string
}

func GetFaceIDWorker(user uint, url string, responseChan chan UintString) {
	log.Println("generating id for ", user)
	id, err := GetFaceID(url)
	if err != nil {
		log.Println(err)
		return
	}
	responseChan <- UintString{user, id}
}

func VerifyFace(faceId1, faceId2 string) (float64, error) {
	mapBody := map[string]string{"faceId1": faceId1, "faceId2": faceId2}

	jsonBytes, err := json.Marshal(mapBody)
	if err != nil {
		return 0, err
	}

	req, err := http.NewRequest("POST", VERIFY_URL, bytes.NewBuffer(jsonBytes))
	if err != nil {
		return 0, err
	}

	for key, value := range FaceAPIHeaders {
		req.Header.Set(key, value)
	}

	resp, err := client.Do(req)

	if resp != nil {
		defer resp.Body.Close()
	}

	if err != nil {
		return 0, nil
	}

	var verifyResponse VerifyResponse

	decoder := json.NewDecoder(resp.Body)

	err = decoder.Decode(&verifyResponse)
	if err != nil {
		return 0, err
	}

	if resp.StatusCode != 200 {
		return 0, errors.New(verifyResponse.Error.Message)
	}

	return verifyResponse.Confidence, nil
}

type SimilarityResponse struct {
	user1      uint
	user2      uint
	similarity float64
}

func VerifyFaceWorker(user1, user2 uint, faceId1, faceId2 string, responseChan chan SimilarityResponse) {
	log.Println("calculating similarity for ", user1, faceId1, user2, faceId2)
	similarity, err := VerifyFace(faceId1, faceId2)
	if err != nil {
		log.Println(err)
		return
	}

	var similarityResponse SimilarityResponse
	similarityResponse.user1 = user1
	similarityResponse.user2 = user2
	similarityResponse.similarity = similarity
	responseChan <- similarityResponse
}

func ManyGetFacesID(nameUrls map[uint]string) {
	faceIdChan := make(chan UintString)
	doneChan := make(chan bool)

	count := 0

	for user, url := range nameUrls {
		go GetFaceIDWorker(user, url, faceIdChan)
		count += 1
	}

	go func() {
		for {
			face_id := <-faceIdChan
			Faces_ids[face_id.Uint] = face_id.String
			count -= 1
			if count == 0 {
				doneChan <- true
				break
			}
		}
	}()

	<-doneChan
}

func ManyVerifyFace(testuser uint, faces_ids map[uint]string) map[uint]float64 {

	similaritiesChan := make(chan SimilarityResponse)

	doneChan := make(chan bool)
	count := 0
	for user, url := range faces_ids {
		if user == testuser {
			continue
		}
		go VerifyFaceWorker(testuser, user, faces_ids[testuser], url, similaritiesChan)
		count += 1
	}

	similarities := make(map[uint]float64)

	go func() {
		for {
			similarityResponse := <-similaritiesChan

			similarities[similarityResponse.user2] = similarityResponse.similarity

			count -= 1
			if count == 0 {
				doneChan <- true
				break
			}
		}
	}()

	<-doneChan

	return similarities
}

func CheckFace(filename string) (uint, string, error) {
	id, err := GetFaceID("http://sparck.co/tempFaces/" + filename)
	if err != nil {
		return 0, "", err
	}
	tempFaces_ids := make(map[uint]string)
	for k, v := range Faces_ids {
		tempFaces_ids[k] = v
	}

	tempFaces_ids[0] = id

	similarities := ManyVerifyFace(0, tempFaces_ids)

	var maxUser uint = 0
	var maxVal float64 = 0
	for user, val := range similarities {
		if val > maxVal {
			maxUser = user
			maxVal = val
		}
	}
	log.Println(similarities)

	if maxVal > 0.65 {
		name, err := passwords.GetName(maxUser)
		if err != nil {
			return 0, "", err
		}
		return maxUser, name, nil
	}

	return 0, "", nil
}

// func NewUser(filename string) (uint, string, error) {
// 	id, err := GetFaceID("http://sparck.co/tempFaces/" + filename)
// 	if err != nil {
// 		return 0, "", err
// 	}
// 	tempFaces_ids := make(map[uint]string)
// 	for k, v := range Faces_ids {
// 		tempFaces_ids[k] = v
// 	}

// 	tempFaces_ids[0] = id

// 	similarities := ManyVerifyFace(0, tempFaces_ids)

// 	var maxUser uint = 0
// 	var maxVal float64 = 0
// 	for user, val := range similarities {
// 		if val > maxVal {
// 			maxUser = user
// 			maxVal = val
// 		}
// 	}
// 	log.Println(similarities)

// 	if maxVal > 0.65 {
// 		name, err := passwords.GetName(maxUser)
// 		if err != nil {
// 			return 0, "", err
// 		}
// 		return maxUser, name, nil
// 	}

// 	return 0, "", nil
// }
