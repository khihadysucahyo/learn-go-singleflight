package main

import (
	"fmt"
	"log"
	"net/http"
	"time"
	"golang.org/x/sync/singleflight"
)

func main() {
	var requestGroup singleflight.Group

	http.HandleFunc("/normal", func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		status, err := githubStatus()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		log.Printf("github handler request: processing time: %+v | status %q", time.Since(start), status)
		fmt.Fprintf(w, "Github status: %q", status)

	})

	http.HandleFunc("/singleflight", func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		v, err, shared := requestGroup.Do("github", func() (interface{}, error) {
			return githubStatus()
		})

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		status := v.(string)

		log.Printf("github handler request: processing time: %+v | status %q | shared %t", time.Since(start), status, shared)
		fmt.Fprintf(w, "Github status: %q", status)

	})

	log.Println("running server on port :8080")
	http.ListenAndServe(":8080", nil)

}

func githubStatus() (string, error) {
	time.Sleep(1 * time.Second)
	log.Println("call github status")
	res, err := http.Get("https://api.github.com")
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return "", fmt.Errorf("github response %s", res.Status)
	}

	return res.Status, nil
}
