package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

type WorldClockApi struct {
	CurrentFileTime int
}

func main() {
	p := message.NewPrinter(language.English)
	totalCount := 0
	successCount := 0

	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		failCount := float64(totalCount - successCount)
		failRate := (failCount / float64(totalCount)) * 100

		fmt.Println("---------- Process Terminated ----------")
		fmt.Printf("%d out of %d request has been completed successfully\n", successCount, totalCount)
		fmt.Printf("%0.2f %% packet loss\n", failRate)
		fmt.Println("---------- Bye Bye ----------")
		os.Exit(1)
	}()

	urlPtr := flag.String("url", "https://www.google.com", "URL to connect")
	flag.Parse()

	fmt.Println("---------- HTTP GET Request Test ----------")
	fmt.Printf("URL: %s\n", *urlPtr)
	fmt.Println("---------- oXo ----------")

	for {
		start := time.Now()
		// resp, err := http.Get("http://worldclockapi.com/api/json/utc/now")
		resp, err := http.Get(*urlPtr)
		elapsed := time.Since(start)

		if err != nil {
			fmt.Printf("Unable to complete request in %d ms\n", elapsed.Milliseconds())
		} else if resp.StatusCode != 200 {
			fmt.Printf("Non 200 status code in %d ms\n", elapsed.Milliseconds())
		} else {
			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				fmt.Printf("Unable to read response in %d ms\n", elapsed.Milliseconds())
			} else {
				sb := string(body)
				// var res WorldClockApi
				// json.Unmarshal([]byte(sb), &res)

				successCount++
				p.Printf("%v bytes of request completed in %d ms\n", len(sb), elapsed.Milliseconds())
			}
		}

		totalCount++
		time.Sleep(1 * time.Second)
	}
}
