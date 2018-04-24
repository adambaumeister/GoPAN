package pan

import (
	"crypto/tls"
	"github.com/beevik/etree"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"time"
    "errors"
    "strconv"
)

//Function to generate an API key
func Keygen(fqdn string, user string, pass string) string {
    
    apiKey := ""
    
    //Validating that all login flags are set
    if (fqdn == "" || user =="" || pass=="") {
        e := "Error: required flags \"ip-address\", \"password\" or \"user\" not set"
        println (e)
        Logerror(errors.New(e), true)
    } 
	//Defining secondary variables
	req, err := url.Parse("https://" + fqdn + "/api/?")
	if err != nil {
		Logerror(err, true)
	}
	q := url.Values{}
	q.Add("password", pass)
	q.Add("user", user)
	q.Add("type", "keygen")
	req.RawQuery = q.Encode()
    resp, err := HttpValidate(req.String(), false)
    if err != nil {
		Logerror(err, true)
	} 
    doc := etree.NewDocument()
    doc.ReadFromBytes(resp)
    for _, e := range doc.FindElements("./response/result/*") {
        apiKey = e.Text()
    }
    return apiKey
}

func Wlog(fileName string, text string, newline bool) {
	// If the file doesn't exist, create it, or append to the file
	if newline {
		text += "\n"
	}
	f, err := os.OpenFile(fileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		Logerror(err, true)
	}
	if _, err := f.Write([]byte(text)); err != nil {
		Logerror(err, true)
	}
	if err := f.Close(); err != nil {
		Logerror(err, true)
	}
}

//Prints the error and exit execution if fatal is set
func Logerror(err error, fatal bool) {
	if err != nil {
		Wlog("error.txt", err.Error(), true)
        if fatal{
		  os.Exit(1)
        }
	}
}

//function to validate succesful http request received a 200 code and a "success" on it. 
//It receives an http request and a debug flag in case you want to see the HTTP calls
//It returns the response and an error if something fails

func HttpValidate (req string, debug bool) ([]byte , error) {
    //Initialazing the error it'll return if anyone it's found.
    var problem error
    //HTTP requests are print in case debug flag is set
    if debug{
        println(req)
    }
    //Ignoring TLS certificate checking
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	//Setting HTTP timeout as 15 seconds.
	netClient := &http.Client{
		Timeout:   time.Second * 15,
		Transport: tr,
	}
    
    resp, err := netClient.Get(req)
	if err != nil {
		Logerror(err, true)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		Logerror(err, true)
	}
	//making sure the API responds with a 200 code and a success on it
    if resp.StatusCode == 200 {
		doc := etree.NewDocument()
		doc.ReadFromBytes(body)
        //extraccting the response status from the http response and comparing it with "success"
        status := doc.FindElement("./*").SelectAttrValue("status", "unknown")
        if status != "success"{
            problem = errors.New("error with HTTP request:\t" + req + "\nreceived status " + status +  " and response :\t" + string(body))
        }
	}else {
        problem = errors.New("error with HTTP request:\t" + req + "\nreceived status code:\t" + strconv.Itoa(resp.StatusCode))
	}

    return body,problem
}
