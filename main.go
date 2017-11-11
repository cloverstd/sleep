package main

import (
	"flag"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

func main() {
	port := flag.Int("port", 1234, "port to listen")
	maxSleep := flag.Int("max-rand", 120, "max second to sleep")
	homePage := flag.String("home", "https://github.com/cloverstd/sleep", "root to redirect")
	flag.Parse()
	if _, err := url.Parse(*homePage); err != nil {
		log.Fatalf("invalid home page: %s", *homePage)
	}

	rand.Seed(time.Now().Unix())

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		accept := r.Header.Get("Accept")
		if strings.Contains(accept, "text/html") {
			http.Redirect(w, r, *homePage, http.StatusMovedPermanently)
			return
		}
		var second int
		var err error
		if second, err = strconv.Atoi(r.URL.Path[1:]); err != nil {
			if r.Header.Get("X-Sleep") == "" {
				second = rand.Intn(*maxSleep)
			} else {
				second, err = strconv.Atoi(r.Header.Get("X-Sleep"))
				if err != nil {
					w.WriteHeader(http.StatusBadRequest)
					fmt.Fprintf(w, "%s is not a invalid sleep second", r.Header.Get("X-Sleep"))
					return
				}
			}
		}
		if second > *maxSleep || second < 0 {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, "bad sleep time. sleep time should between in (1, %d)", *maxSleep)
			return
		}
		notify := w.(http.CloseNotifier).CloseNotify()
		t := time.NewTimer(time.Duration(second) * time.Second)
		defer t.Stop()
		select {
		case <-t.C:
			w.Header().Set("X-Sleep", fmt.Sprint(second))
			fmt.Fprintf(w, "I sleep %ds.", second)
			return
		case <-notify:
			log.Printf("Close by client <%s>", r.RemoteAddr)
			return
		}
	})

	log.Printf("Sleep server is running on :%d", *port)
	http.ListenAndServe(fmt.Sprintf(":%d", *port), nil)
}
