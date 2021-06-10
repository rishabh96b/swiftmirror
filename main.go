package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	mirror "github.com/rishabh96b/swiftmirror/pkg/mirror"
)

type response struct {
	FastestURL string        `json:"fastest_url"`
	Latency    time.Duration `json:"latency"`
}

func getFastestURL(urlList []string) response {
	urlChannel := make(chan string)
	latencyChannel := make(chan time.Duration)

	for _, v := range urlList {
		mirrorURL := v
		go func() {
			start := time.Now()
			_, err := http.Get(mirrorURL + "/README")
			latency := time.Now().Sub(start) / time.Millisecond
			if err == nil {
				urlChannel <- mirrorURL
				latencyChannel <- latency
			}
		}()
	}

	return response{
		<-urlChannel,
		<-latencyChannel,
	}
}

func main() {
	fmt.Println("This is working")
	http.HandleFunc("/mirror", func(rw http.ResponseWriter, r *http.Request) {
		log.Println("Requested Path is: ", r.URL.Path)
		mirrorList, err := json.Marshal(mirror.DebianMirrorList)
		if err != nil {
			log.Println("Cannot marshal Debian mirror list")
			mirrorList = []byte(`{"mirror":"not availabe"}`)
		}
		rw.WriteHeader(http.StatusOK)
		rw.Header().Set("Content-Type", "application/json")
		rw.Write(mirrorList)

	})
	http.HandleFunc("/fast-mirror", func(rw http.ResponseWriter, r *http.Request) {
		log.Println("Requested Path is: ", r.URL.Path)
		response := getFastestURL(mirror.DebianMirrorList[:])
		respJSON, _ := json.Marshal(response)
		rw.WriteHeader(http.StatusOK)
		rw.Header().Set("Content-Type", "application/json")
		rw.Write(respJSON)

	})
	log.Print("Serving on port 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
