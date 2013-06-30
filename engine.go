package main

import (
	"encoding/json"
	"fmt"
	"github.com/cactus/go-statsd-client/statsd"
	"github.com/gorilla/mux"
	"io/ioutil"
	"log"
	"net/http"

	"os"
	"time"
)

// put it in hdfs
// put it in a file
// put it in redis
// put it in a rdbms
// put it in a message queue
// put it in elasticsearch
// put it on another handler within runny
// write in avro format
// write into irc
// write to null
// write to hbase

// read from restful api
// read from http
// read from log4j style
// read from port
// read from directory
// read from ftp
// read from message queue

// buffer to file
// buffer to redis
// buffer to queue
// buffer to memory

// interceptors

// api to get status from

// statsd/collectd/metricsd/??/ support

// security

// support different encodings such as json, avro, etc

// support fan out via a component

// support queuing via channels

var VERSION float32 = 1.1

func main() {
	startupBegan := time.Now()
	// first create a client
	client, err := statsd.Dial("127.0.0.1:8125", "test-client")
	// handle any errors
	if err != nil {
		log.Fatal(err)
	}
	// make sure to clean up
	defer client.Close()

	// Send a stat
	err = client.Inc("stat1", 42, 1.0)
	// handle any errors
	if err != nil {
		log.Printf("Error sending metric: %+v", err)
	}
	fmt.Printf("Starting Runny v%v\n", VERSION)
	//	makeConfig()
	var settings Config = readConfig()
	var router = mux.NewRouter()
	router.StrictSlash(true) // trailing slashes resolve to no slash
	router.HandleFunc("/", homeHandler)
	router.HandleFunc("/hello", helloHandler)

	for _, process := range settings.Processes {
		// index is the index where we are
		// element is the element from someSlice for where we are
		router.HandleFunc(process.MyHandler.Url, getFunctionPointer(process.MyHandler.Name)).Methods(process.MyHandler.Methods...)
	}

	http.Handle("/", router)

	client.Timing("runny.startup", 1, time.Since(startupBegan))
	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}

func getFunctionPointer(name string) (f func(http.ResponseWriter, *http.Request)) {
	if "webhookCallbackHandler" == name {
		return helloHandler
	}
	if "webhookCallbackHandler2" == name {
		return homeHandler
	} else {
		return nilHandler
	}
}

func printError(msg string, params ...string) {
	if false {
		log.Printf(msg, params)
	}
}

func makeConfig() {
	h1 := Handler{Name: "webhookCallbackHandler", Url: "/url/webhook"}
	h2 := Handler{Name: "webhookCallbackHandler", Url: "/url/webhook/2"}
	h1.Methods = append(h1.Methods, "get")
	h1.Methods = append(h1.Methods, "put")
	h2.Methods = append(h1.Methods, "get")

	p1 := Process{Id: 1, Name: "test1", MyHandler: h1}
	p2 := Process{Id: 2, Name: "test2", MyHandler: h2}
	var c Config = New()
	c.Name = "test"
	c.Processes = append(c.Processes, p1, p2)
}

func writeFile(contents string) {
	// write whole the body
	err := ioutil.WriteFile("sample.config.json", []byte(contents), 0644)
	if err != nil {
		panic(err)
	}
}

func readConfig() Config {
	// then config file settings

	configFile, err := os.Open("config.json")
	if err != nil {
		printError("opening config file", err.Error())
	}
	var settings Config
	jsonParser := json.NewDecoder(configFile)
	if err = jsonParser.Decode(&settings); err != nil {
		printError("parsing config file", err.Error())
	}

	fmt.Printf("read config: %v\n", settings.String())
	return settings
}

func nilHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "text/html")
	fmt.Fprint(w, "Nil Handler")
}
func homeHandler(w http.ResponseWriter, r *http.Request) {
	printError("home handler")
	w.Header().Add("Content-Type", "text/html")
	fmt.Fprint(w, "Try a <a href='/Hello/world'>hello</a>.")
}

func helloHandler(w http.ResponseWriter, r *http.Request) {
	printError("hello handler")
	w.Header().Add("Content-Type", "text/html")
	phrase := r.FormValue("salutation") + ", " + r.FormValue("name") + "!"
	fmt.Fprint(w, phrase)
}

type Config struct {
	Name      string
	Processes []Process `json:"processes"`
}

func New() Config {
	return Config{Name: "testing"}
}

func (c *Config) String() string {
	b, _ := json.Marshal(c)
	return string(b)
}

type Process struct {
	Id        int32   `json:"id"`
	Name      string  `json:"name"`
	MyHandler Handler `json:"myHandler"`
}

type Handler struct {
	Name    string   `json:"name"`
	Url     string   `json:"url"`
	Methods []string `json:"method"`
}
