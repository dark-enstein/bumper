package main

import (
	"encoding/json"
	"fmt"
	"github.com/adamwasila/go-semver"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"strconv"
	"time"
)

const (
	MODEMAJOR = "major"
	MODEMINOR = "minor"
	MODEPATCH = "patch"
)

var (
	Ver                     = "0.0.1"
	ERR_SUCCESSFUL          = "0"
	ERR_BYTESUNREADABLE     = "1"
	ERR_REQDATAEMPTY        = "2"
	ERR_INTERNALSERVERERROR = "5"
)

type BumpRequest struct {
	Version        string `json:"version"`
	CurrentVersion string `json:"currentVersion"`
	Class          string `json:"class"`
}

type BumpResponse struct {
	StatusCode string `json:"statusCode"`
	NewVersion string `json:"newVersion"`
}

// Version represents the current version of the application
type Version struct {
	Version string `json:"version"`
}

// bumpHandler handles the POST requests to bump the version
func bumpHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Received request to /bump")
	// For simplicity, we will just increment the version string
	w.Header().Set("Date", time.Now().UTC().String())
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Connection", "close")
	w.Header().Set("Server", fmt.Sprintf("bumper:%s\n", Ver))
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Credentials", "true")
	var bReq BumpRequest
	var bResp BumpResponse
	err := json.NewDecoder(r.Body).Decode(&bReq)
	if err != nil {
		log.Printf("Error while decoding request body: %s\n", err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	log.Printf("Message received: %#v\n", bReq)
	var mode int
	switch bReq.Class {
	case MODEMAJOR:
		mode = 1
	case MODEMINOR:
		mode = 2
	case MODEPATCH:
		mode = 3
	}
	log.Printf("Request ver: %s\n", bReq.CurrentVersion)
	v, err := Bump(bReq.CurrentVersion, mode)
	if err != nil {
		log.Printf("Go-Version Error: %s\n", err.Error())
		bResp.StatusCode = ERR_INTERNALSERVERERROR
		_ = json.NewEncoder(w).Encode(bResp)
	}
	log.Printf("Generated next ver: %s\n", v)
	bResp.StatusCode = ERR_SUCCESSFUL
	bResp.NewVersion = v
	b, err := json.Marshal(bResp)
	if err != nil {
		bResp.StatusCode = ERR_INTERNALSERVERERROR
		_ = json.NewEncoder(w).Encode(bResp)
	}
	w.Header().Set("Content-Length", strconv.Itoa(len(b)))
	w.WriteHeader(200)
	_, _ = w.Write(b)
}

// versionHandler handles the GET requests to retrieve the version
func versionHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Received request to /version")
	w.Header().Set("Date", time.Now().UTC().String())
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Connection", "close")
	w.Header().Set("Server", fmt.Sprintf("bumper:%s\n", Ver))
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Credentials", "true")
	var v Version
	v.Version = Ver
	log.Println("Version info:", v.Version)
	b, err := json.Marshal(v)
	if err != nil {
		log.Printf("Go-Version Error: %s\n", err.Error())
		_ = json.NewEncoder(w).Encode(v)
	}
	w.Header().Set("Content-Length", strconv.Itoa(len(b)))
	resp, err := json.Marshal(v)
	if err != nil {
		log.Printf("Go-Version Error: %s\n", err.Error())
		_ = json.NewEncoder(w).Encode(v)
	}
	w.Header().Set("Content-Length", strconv.Itoa(len(resp)))
	w.WriteHeader(200)
	_, _ = w.Write(resp)
}

func Bump(s string, mode int) (string, error) {
	v, err := semver.Parse(s)
	if err != nil {
		return "", fmt.Errorf("Error bumping version: %s\n", err.Error())
	}
	var v2 semver.Version

	switch mode {
	case 1:
		v2, err = v.Bump(semver.NextMajor())
	case 2:
		v2, err = v.Bump(semver.NextMinor())
	case 3:
		v2, err = v.Bump(semver.NextPatch())
	}
	return v2.String(), err
}

func main() {
	r := mux.NewRouter()

	// Define the endpoints
	r.HandleFunc("/bump", bumpHandler).Methods("POST")
	r.HandleFunc("/version", versionHandler).Methods("GET")

	// Start the server
	log.Println("Starting server on :8080")
	srv := &http.Server{
		Handler: r,
		Addr:    "127.0.0.1:8080",
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Fatal(srv.ListenAndServe())
}
