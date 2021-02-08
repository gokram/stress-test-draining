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



fmt.Printf("MM_TRACE %s %d %d %s %s %s Generated_log\n ",batchid,TIMEOUT_SEC,i,INDEX,application_name,currentTime)

i++

//time.Sleep(1* time.Second)
        } //end select  

//fmt.Println("BUILDBATCH - Exited SELECT")

    } //end fortimeout/loop

fmt.Printf("%d raw data generated - Logging - FISNISHEDTEST STRESS\n",i)

////////////////
}()

//time.Sleep(60* time.Second)

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











