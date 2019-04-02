package main
import (
	"fmt"
	"os"
	"log"
	"net/url"
	"net/http"
	"io/ioutil"
	"regexp"
	"io"
	"flag"
	//"time"
)
func main() {

	//fmt.Println("Please Provide the Arguments in the sequence apikey=xxxxxx-xxxx-xxxx-xxxx-xxxxxx&statuses=4,5&from=2019-02-25T00:00:00&to=2019-03-03T23:59:59")
	//argsWithProg := os.Args
	//fmt.Println(argsWithProg)
    //argsWithoutProg := os.Args[1:]
    // API_KEY := os.Args[1]
    // STATUSES := os.Args[2]
    // D_FROM := os.Args[3]
    // d_TO := os.Args[4]
    API_KEY_Ptr := flag.String("apikey", "xxx-xxx-xxxx-xxxx-xxxx", "a string")
    STATUSES_Ptr := flag.String("statuses", "x,y", "a string")
    FROM_Ptr := flag.String("from", "2018-01-01:00:00:00", "a string")
    TO_Ptr := flag.String("to", "2019-01-01:00:00:00", "a string")
    flag.Parse()

    // fmt.Println("apikey:", *API_KEY_Ptr)
    // fmt.Println("statuses:", *STATUSES_Ptr)
    // fmt.Println("from:", *FROM_Ptr)
    // fmt.Println("to", *TO_Ptr)

    // fmt.Println(API_KEY, STATUSES, D_FROM, d_TO)
    
    //fmt.Println(argsWithoutProg)
    //fmt.Println(arg)

    u, err := url.Parse("http://bing.com/search")
	if err != nil {
		log.Fatal(err)
	}
	u.Scheme = "https"
	u.Host = "api.elasticemail.com"
	u.Path = "v2/log/exportevents"
	q := u.Query()
	q.Set("apikey", *API_KEY_Ptr)
	q.Set("statuses", *STATUSES_Ptr)
	q.Set("from", *FROM_Ptr)
	q.Set("to", *TO_Ptr)
	u.RawQuery = q.Encode()
	fmt.Println(u)

	//url := "https://api.elasticemail.com/v2/log/exportevents?apikey=3da0a60f-37a8-4363-bfd6-1cb0dbc525ee&statuses=4,5&from=2019-02-25T00:00:00&to=2019-03-03T23:59:59"

	url := u.String()

	req, _ := http.NewRequest("GET", url, nil)

	res, _ := http.DefaultClient.Do(req)

	defer res.Body.Close()

	body, _ := ioutil.ReadAll(res.Body)

	fmt.Println("The Status of the query is:")
	
	fmt.Println(string(body))

    myString := string(body)
    
    pat := regexp.MustCompile(`https?://.*\.csv`)
	s := pat.FindString(myString)
	fmt.Println("The extracted url is:")
	fmt.Println(s)

	//fileUrl := "https://api.elasticemail.com/userfile/55fa99ce-d2c7-41ad-9cad-e9d6476ba136/export/546dc7d5-bdf8-4842-917a-06b65092e347-eventslog.csv"
	fileUrl := s
	DOWNLOADFILE_NAME := "eventslog-" + string(*FROM_Ptr) + "-" + string(*TO_Ptr)
    // if err := DownloadFile("email-events", fileUrl); err != nil {
       if err := DownloadFile(DOWNLOADFILE_NAME, fileUrl); err != nil {
        panic(err)
    }
}

// DownloadFile will download a url to a local file. It's efficient because it will
// write as it downloads and not load the whole file into memory.
	func DownloadFile(filepath string, url string) error {

	    // Get the data
	    resp, err := http.Get(url)
	    if err != nil {
	        return err
	    }
	    defer resp.Body.Close()

	    // Create the file
	    out, err := os.Create(filepath)
	    if err != nil {
	        return err
	    }
	    defer out.Close()

	    // Write the body to file
	    _, err = io.Copy(out, resp.Body)
	    return err

}