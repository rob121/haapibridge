package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/spf13/viper"
	"io/ioutil"
	"net/http"
	"time"
	"github.com/rob121/vhelp"
	"log"
)

var conf *viper.Viper

func main(){

    vhelp.Load("config")
	conf,_ = vhelp.Get("config")

    startHttpServer()

}


func startHttpServer() {

	r := mux.NewRouter()
	r.HandleFunc("/api/states/{entity_id}/{state}", statesHandler)

	srv := &http.Server{
		Handler:      r,
		Addr:         fmt.Sprintf(":%s",conf.GetString("port")),
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Fatal(srv.ListenAndServe())
}


func statesHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	entity_id := vars["entity_id"]

	postBody, _ := json.Marshal(map[string]string{
		"state":  vars["state"],
	})

	url := fmt.Sprintf("%s/%s/%s",conf.GetString("haurl"),"api/states",entity_id)

	responseBody := bytes.NewBuffer(postBody)

	log.Println(url,responseBody,entity_id,vars["state"])

	req, err := http.NewRequest("POST", url, responseBody)

	if err != nil {
		fmt.Fprint(w,"ERROR")
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s",conf.GetString("key")))
	req.Header.Set("Content-Type", "application/json")


	//Leverage Go's HTTP Post function to make request
	resp, err := http.DefaultClient.Do(req)
	//Handle Error
	if err != nil {
		log.Fatalf("An Error Occured %v", err)
	}

	defer resp.Body.Close()
	//Read the response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}
	sb := string(body)

	fmt.Fprint(w,sb)

}
