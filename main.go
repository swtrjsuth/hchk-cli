package main

import (
	"fmt"
	"time"
	"flag"
	"encoding/json"
	"bytes"
	"net/http"
	"io"
	"os"
)

func main() {
	var flagConfigFile string
	flag.StringVar(&flagConfigFile, "conf", "", "Full filesystem path to config file (JSON)")

	flag.Parse()

	var confdata map[string]interface{}

	conftext, err := os.ReadFile(flagConfigFile)
	if err != nil {
		panic(err)
	}

	if err := json.Unmarshal(conftext, &confdata); err != nil {
		panic(err)
	}

	firebaseUrl := confdata["fb_url"].(string)
	loginUrl := confdata["login_url"].(string)
	email := confdata["email"].(string)
	password := confdata["password"].(string)
	expectEvery := confdata["every"].(string)

	mapD := map[string]string{"email": email, "password": password, "returnSecureToken": "true"}
	mapB, _ := json.Marshal(mapD)

	reader := bytes.NewReader(mapB)
	resp, err := http.Post(loginUrl, "application/json", reader)

	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	body, err2 := io.ReadAll(resp.Body)
	if err2 != nil {
		panic(err2)
	} else {

		var authdata map[string]interface{}

		if err := json.Unmarshal(body, &authdata); err != nil {
			panic(err)
		}

		uid := authdata["localId"]
		idtoken := authdata["idToken"]
		timestamp := fmt.Sprintf("%d",time.Now().Unix())

		payloadD := map[string]string{"t": timestamp, "e": expectEvery}
		payloadB, _ := json.Marshal(payloadD)

		updateUrl := fmt.Sprintf("%s/workloads/%s.json?auth=%s", firebaseUrl, uid, idtoken)

		client := &http.Client{}

		req, err := http.NewRequest(http.MethodPatch, updateUrl, bytes.NewBuffer(payloadB))
		req.Header.Set("Content-Type", "application/json")
		if err != nil {
			panic(err)
		}

		resp, err := client.Do(req)
		if err != nil {
			panic(err)
		}

		defer resp.Body.Close()

		_, err2 := io.ReadAll(resp.Body)
		if err2 != nil {
			panic(err2)
		}

	}
}

