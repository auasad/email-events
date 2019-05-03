package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"time"
)

const (
	// The statuses API parameter specifies event types: 0 for all, 1 for ReadyToSend, 2 for InProgress, 4 for Bounced, 5 for Sent, 6 for Opened, 7 for Clicked, 8 for Unsubscribed, 9 for Abuse Report.
	defaultStatuses = "4,5"
	timeFormat      = "2006-01-02T15:04:05"
)

func main() {

	// Command line parameters must be available to specify api-key, statuses and date range.
	// There should be flags to specify yesterday, last hour, and last 5 minute interval (0-5, 5-10, ...) in addition to custom range.
	var (
		apiKey   string
		statuses string
		from     string
		to       string
	)

	// Yesterday in local time zone converted to UTC
	t := time.Now()
	defaultTo := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, time.Local).UTC()
	defaultFrom := defaultTo.AddDate(0, 0, -1)

	// to take the input as a command parameters.
	flag.StringVar(&apiKey, "apikey", "", "Elastic Email API key")
	flag.StringVar(&statuses, "statuses", defaultStatuses, "Event types to include")
	flag.StringVar(&from, "from", defaultFrom.Format(timeFormat), "Start time for events")
	flag.StringVar(&to, "to", defaultTo.Format(timeFormat), "End time for events")
	flag.Parse()

	if apiKey == "" {
		fmt.Fprintf(os.Stderr, "Please specify -apikey, see %s -help\n", os.Args[0])
		os.Exit(1)
	}

	// creates an http client and hits the elastic email api and stores the JASON output.
	link, err := exportEvents(apiKey, statuses, from, to)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
	fmt.Printf("Downloading from %s\n", link)

	//randers the downloadable file name with from and to dates.
	filename := "eventslog-" + from + "-" + to + ".csv"

	//calling a function that downloads the file, takes url as input and file name to be given
	err = downloadFile(filename, link)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
}

type exportInfo struct {
	Success bool `json:"success"`
	Data    struct {
		Link           string `json:"link"`
		PublicExportID string `json:"publicexportid"`
	} `json:"data"`
}

func exportEvents(key, statuses, from, to string) (link string, err error) {
	// The "from" and "to" API parameters are in UTC. The "to" parameter means "up to and including".
	q := url.Values{}
	q.Set("apikey", key)
	q.Set("statuses", statuses)
	q.Set("from", from)
	q.Set("to", to)
	url := "https://api.elasticemail.com/v2/log/exportevents?" + q.Encode()
	resp, err := http.Get(url)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	// decode json response body
	var info exportInfo
	err = json.NewDecoder(resp.Body).Decode(&info)
	if err != nil {
		return
	}
	if !info.Success {
		err = fmt.Errorf("exportevents not successful") // TODO: check response body
		return
	}
	link = info.Data.Link
	return
}

// function that takes file path and url and downloads the file
func downloadFile(filepath string, url string) error {

	// Poll and wait until export is available
	var body io.Reader
	wait := 1
	for {
		resp, err := http.Get(url)
		if err != nil {
			return err
		}

		if resp.StatusCode == 200 {
			defer resp.Body.Close()
			body = resp.Body
			break
		} else {
			if resp.StatusCode == 404 && wait <= 32 {
				fmt.Printf("No document yet, waiting %d seconds...\n", wait)
				time.Sleep(time.Second * time.Duration(wait))
				wait *= 2
			} else {
				return fmt.Errorf("Status %d", resp.StatusCode)
			}
		}
	}

	// Create the file
	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	// Write the body to file
	_, err = io.Copy(out, body)
	return err
}
