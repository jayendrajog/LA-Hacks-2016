package face_auth

import (
	"bytes"
	"encoding/json"
	"errors"
	"log"
	"net/http"
)

const DETECT_URL = "https://api.projectoxford.ai/face/v1.0/detect?returnFaceId=true&returnFaceLandmarks=false&returnFaceAttributes=age"
const VERIFY_URL = "https://api.projectoxford.ai/face/v1.0/verify"

var FaceAPIHeaders = map[string]string{
	"Content-Type":              "application/json",
	"Ocp-Apim-Subscription-Key": "38c44ac804c44f6e97673d815163a1db",
}

var FACES_DB = map[string]string{
	"JAY":   "http://sparck.co/faces/JAY.jpg",
	"ADAM":  "http://sparck.co/faces/ADAM.jpg",
	"ANNA":  "http://sparck.co/faces/anna.jpg",
	"JAHAN": "http://sparck.co/faces/JAHAN.jpg",
}

var Faces_ids = make(map[string]string)

var client *http.Client

func Init() {
	client = &http.Client{}
	ManyGetFacesID(FACES_DB)
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

	return faceResponse[0].FaceId, nil
}

func GetFaceIDWorker(name, url string, responseChan chan [2]string) {
	log.Println("generating id for ", name)
	id, err := GetFaceID(url)
	if err != nil {
		log.Println(err)
		return
	}
	responseChan <- [2]string{name, id}
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
	name1      string
	name2      string
	similarity float64
}

func VerifyFaceWorker(name1, name2, faceId1, faceId2 string, responseChan chan SimilarityResponse) {
	log.Println("calculating similarity for ", name1, faceId1, name2, faceId2)
	similarity, err := VerifyFace(faceId1, faceId2)
	if err != nil {
		log.Println(err)
		return
	}

	var similarityResponse SimilarityResponse
	similarityResponse.name1 = name1
	similarityResponse.name2 = name2
	similarityResponse.similarity = similarity
	responseChan <- similarityResponse
}

func ManyGetFacesID(nameUrls map[string]string) {
	faceIdChan := make(chan [2]string)
	doneChan := make(chan bool)

	count := 0

	for name, url := range FACES_DB {
		go GetFaceIDWorker(name, url, faceIdChan)
		count += 1
	}

	go func() {
		for {
			face_id := <-faceIdChan
			Faces_ids[face_id[0]] = face_id[1]
			count -= 1
			if count == 0 {
				doneChan <- true
				break
			}
		}
	}()

	<-doneChan
}

func ManyVerifyFace(testname string, faces_ids map[string]string) map[string]float64 {

	similaritiesChan := make(chan SimilarityResponse)

	doneChan := make(chan bool)
	count := 0
	for name, url := range faces_ids {
		if name == testname {
			continue
		}
		go VerifyFaceWorker(testname, name, faces_ids[testname], url, similaritiesChan)
		count += 1
	}

	similarities := make(map[string]float64)

	go func() {
		for {
			similarityResponse := <-similaritiesChan

			similarities[similarityResponse.name2] = similarityResponse.similarity

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

func CheckFace(testname string) (string, error) {
	id, err := GetFaceID("http://sparck.co/tempFaces/" + testname)
	if err != nil {
		return "", err
	}
	tempFaces_ids := make(map[string]string)
	for k, v := range Faces_ids {
		tempFaces_ids[k] = v
	}

	tempFaces_ids[testname] = id

	similarities := ManyVerifyFace(testname, tempFaces_ids)

	maxName := ""
	var maxVal float64 = 0
	for name, val := range similarities {
		if val > maxVal {
			maxName = name
			maxVal = val
		}
	}
	log.Println(similarities)

	if maxVal > 0.7 {
		return maxName, nil
	}

	return "", nil

}
