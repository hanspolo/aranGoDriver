package aranGoDriver

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"encoding/json"

	"github.com/TobiEiss/aranGoDriver/models"
)

// AranGoSession represent to Session
type AranGoSession struct {
	urlRoot   string
	jwtString string
}

const urlAuth = "/_open/auth"

// NewAranGoDriverSession creates a new instance of a AranGoDriver-Session.
// Need a host (e.g. "http://localhost:8529/")
func NewAranGoDriverSession(host string) *AranGoSession {
	return &AranGoSession{host, ""}
}

// Connect to arangoDB
func (session AranGoSession) Connect(username string, password string) {
	credentials := models.Credentials{}
	credentials.Username = username
	credentials.Password = password

	resp := post(&session, urlAuth, credentials)
	session.jwtString = resp["jwt"].(string)
}

func post(session *AranGoSession, url string, object interface{}) map[string]interface{} {
	// marshal body
	jsonBody, err := json.Marshal(object)
	failOnError(err, "Cant marshal object")

	// build url
	url = session.urlRoot + url
	fmt.Println("URL:>", url)

	// build request
	var jsonString = []byte(jsonBody)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonString))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)

	failOnError(err, "Cant do post-request to "+url)

	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)

	// unmarshal to map
	var responseMap map[string]interface{}
	err = json.Unmarshal(body, &responseMap)
	return responseMap
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}
