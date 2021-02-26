package main

import (
    	"fmt"
    	"os"
	"time"

	"log"
	"net/http"
	"io/ioutil"
	"encoding/json"
	"strconv"

"github.com/valyala/fasthttp"

)



func main() {

vcap_application := os.Getenv("VCAP_APPLICATION")
	if vcap_application == "" {
		fmt.Printf("No VCAP_APPLICATION")
		os.Exit(1)
	}



TIMEOUT_SEC, err := strconv.ParseInt(os.Getenv("TIMEOUT_SEC"), 0, 64)
if err != nil {
    panic(err)
}
fmt.Println(TIMEOUT_SEC)
var INDEX = os.Getenv("CF_INSTANCE_INDEX")


SLEEP_MICROSEC, err := strconv.ParseInt(os.Getenv("SLEEP_MICROSEC"), 0, 64)
if err != nil {
    panic(err)
}
fmt.Println("SLEEP_MICROSEC: ",SLEEP_MICROSEC)

var URL = os.Getenv("URL")

var application_data map[string]interface{}

var i int

err_2 := json.Unmarshal([]byte(vcap_application), &application_data)

	if err_2 != nil {
		log.Printf("Failed to unmarshal (via JSON) message (%s): %s", string(vcap_application), err_2)
		return
	}

application_name := application_data["application_name"].(string)
organization_name := application_data["organization_name"].(string)
space_name := application_data["space_name"].(string)
fmt.Println("TEST NEW JSON :  -> application_name:", application_name)
fmt.Println("TEST NEW JSON :  -> organization_name:", organization_name)
fmt.Println("TEST NEW JSON :  -> space_name:", space_name)



go func() {


/// Client fasthttp

var w HTTPSWriter

w.client = httpClient()


/////// NEW PART

    fmt.Println("Logging - START TEST STRESS")

	now := time.Now()
	batchid := now.Format("20060102150405")
	
// LOOP FOR SOME SECONDS

loop:
for timeout := time.After(time.Second * time.Duration(TIMEOUT_SEC)); ; {

     select {
        case <-timeout:
		fmt.Printf("Logging - Stress test completed Timeout expired\n")
		break loop

        default:

	now := time.Now()
	//currentTime := now.Format("20060102150405")

currentTime := now.Format("2006-01-02T15:04:05.999999999Z07:00")


//var msg = "MM_TRACE " + "batchid" + " " + "SLEEP_MICROSEC" + " " + i + " " + INDEX + " " + application_name  + " " + currentTime + " Generated_log"

msg := fmt.Sprintf("MM_TRACE %s %d %d %s %s %s Generated_log\n ",batchid,SLEEP_MICROSEC,i,INDEX,application_name,currentTime)



//var msg = "stress test fasthtp body"

		req := fasthttp.AcquireRequest()
		req.SetRequestURI(URL)
		req.Header.SetMethod("POST")
		req.Header.SetContentType("text/plain")
		req.SetBody([]byte(msg))

		resp := fasthttp.AcquireResponse()

		err := w.client.Do(req, resp)


		if err != nil {
			 fmt.Errorf("error %d ", err)
		}

		if resp.StatusCode() < 200 || resp.StatusCode() > 299 {
			fmt.Errorf("syslog Writer: Post responded with %d status code", resp.StatusCode())
		}

i++

time.Sleep(time.Duration(SLEEP_MICROSEC)* time.Microsecond)
//time.Sleep(SLEEP_MICROSEC* time.Microsecond)
        } //end select  

//fmt.Println("BUILDBATCH - Exited SELECT")

    } //end fortimeout/loop

fmt.Printf("%d raw data generated - Logging - FISNISHEDTEST STRESS\n",i)

////////////////
}()



handler:= NewSyslog(10)

http.ListenAndServe(fmt.Sprintf(":%s", os.Getenv("PORT")), handler)

fmt.Println("MAIN - main -> HTTP step")



}

type Handler struct {
	messages chan string
	
}

func NewSyslog(BUFFER_RAW int64) *Handler {
	return &Handler{
		messages: make(chan string,BUFFER_RAW),
	}
}

// FUNCTION FASTHTTP


type HTTPSWriter struct {

	client          *fasthttp.Client

}

func httpClient() *fasthttp.Client {
	return &fasthttp.Client{
		MaxConnsPerHost:     5,
		MaxIdleConnDuration: 90 * time.Second,
		
		ReadTimeout:         20 * time.Second,
		WriteTimeout:        20 * time.Second,
	}
}

// FUNCTION -----------------------------------


func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Printf("ServeHTTP - Failed to read body: %s", err)
		return
	}

	if len(body) < 1 {
		w.WriteHeader(http.StatusBadRequest)
		log.Print("ServeHTTP - Empty body")
		return
	}


	return
}











