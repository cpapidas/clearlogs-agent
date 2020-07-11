package main

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

func main() {
	http.HandleFunc("/hello", func(writer http.ResponseWriter, request *http.Request) {
		mockLogMessage := `{"request_id":"b1c564ca-a239-47af-9ce9-ad1b043ed280","severity":"ANY","time":"2020-07-11T13:57:28.686+03:00","program_name":null,"message":{"prop":{"status":200,"time_spent":0.023198,"headers":{"X-Frame-Options":"SAMEORIGIN","X-XSS-Protection":"1; mode=block","X-Content-Type-Options":"nosniff","X-Download-Options":"noopen","X-Permitted-Cross-Domain-Policies":"none","Referrer-Policy":"strict-origin-when-cross-origin","Content-Type":"text/html; charset=utf-8"}}}}`
		log.Println(mockLogMessage)
		fmt.Fprintf(writer, "hello\n")
	})
	go func() {
		for {
			time.Sleep(1 * time.Second)
			mockLogMessage := `{"request_id":"b1c564ca-a239-47af-9ce9-ad1b043ed280","severity":"ANY","time":"2020-07-11T13:57:28.686+03:00","program_name":null,"message":{"prop":{"status":200,"time_spent":0.023198,"headers":{"X-Frame-Options":"SAMEORIGIN","X-XSS-Protection":"1; mode=block","X-Content-Type-Options":"nosniff","X-Download-Options":"noopen","X-Permitted-Cross-Domain-Policies":"none","Referrer-Policy":"strict-origin-when-cross-origin","Content-Type":"text/html; charset=utf-8"}}}}`
			log.Println(mockLogMessage)
		}
	}()
	err := http.ListenAndServe(":8090", nil)
	if err != nil {
		log.Fatalf("failed to start http server with error: %v", err)
	}
}
