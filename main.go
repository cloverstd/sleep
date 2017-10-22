package main

import (
	"flag"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"time"
)

func main() {
	port := flag.Int("port", 1234, "port to listen")
	maxSleep := flag.Int("max-rand", 120, "max second to sleep")
	flag.Parse()

	rand.Seed(time.Now().Unix())

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
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
		time.Sleep(time.Duration(second) * time.Second)
		w.Header().Set("X-Sleep", fmt.Sprint(second))
		fmt.Fprintf(w, "I sleep %ds.", second)
	})

	log.Printf("Sleep server is running on :%d", *port)
	http.ListenAndServe(fmt.Sprintf(":%d", *port), nil)
}
