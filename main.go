package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	log "github.com/sirupsen/logrus"
)

type baseInfo struct {
	Cluster  string `json:"cluster"`
	Hostname string `json:"base"`
	crashed  bool
}

var info baseInfo
var baseCrashed bool

func httpHandler(w http.ResponseWriter, r *http.Request) {
	if info.crashed {
		w.WriteHeader(http.StatusGone)
		return
	}

	switch method := r.Method; method {
	case http.MethodDelete:
		crashBase(w)
	case http.MethodGet:
		printInfo(w)
	default:
		notAllowed(w)
	}
}

func crashBase(w http.ResponseWriter) {
	info.crashed = true

	w.WriteHeader(http.StatusInternalServerError)
	msg := fmt.Sprintf("Panic: base %s is under attack!", info.Hostname)
	w.Write([]byte(msg))
}

func printInfo(w http.ResponseWriter) {
	jsonInfo, err := json.Marshal(info)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(jsonInfo))
	w.Write([]byte("\n"))
}

func notAllowed(w http.ResponseWriter) {
	w.WriteHeader(http.StatusForbidden)
}

func main() {
	cluster := os.Getenv("CLUSTER_NAME")
	if cluster == "" {
		log.Fatal("Failed to get cluster name")
	}
	hostname, err := os.Hostname()
	if err != nil {
		log.Fatal("Failed to get hostname")
	}
	info = baseInfo{
		Cluster:  cluster,
		Hostname: hostname,
	}

	http.HandleFunc("/v1/info", httpHandler)
	log.Info("Starting server on port 8000")
	log.Fatal(http.ListenAndServe(":8000", nil))
}
