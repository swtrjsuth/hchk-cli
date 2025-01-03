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
        "strconv"
)

func main() {
	var flagConfigFile string
        var flagQueryMode string

	flag.StringVar(&flagConfigFile, "conf", "", "Full filesystem path to config file (JSON)")
        flag.StringVar(&flagQueryMode, "q", "", "Query health status of all hosts")

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

                var updateUrl string
                httpMethod := http.MethodPatch
                var payloadB []byte

                if flagQueryMode == "" {
                    updateUrl = fmt.Sprintf("%s/workloads/%s.json?auth=%s", firebaseUrl, uid, idtoken)
		    payloadD := map[string]string{"t": timestamp, "e": expectEvery}
		    payloadB, _ = json.Marshal(payloadD)
                } else {
                    httpMethod = http.MethodGet
                    updateUrl = fmt.Sprintf("%s/workloads.json?auth=%s", firebaseUrl, idtoken)
                }

		client := &http.Client{}

		req, err := http.NewRequest(httpMethod, updateUrl, bytes.NewBuffer(payloadB))
		req.Header.Set("Content-Type", "application/json")
		if err != nil {
			panic(err)
		}

		resp, err := client.Do(req)
		if err != nil {
			panic(err)
		}

		defer resp.Body.Close()

		rb2, err2 := io.ReadAll(resp.Body)
		if err2 != nil {
			panic(err2)
		} else if flagQueryMode != "" {
                    var hcdata map[string]interface{}

                    if err := json.Unmarshal(rb2, &hcdata); err != nil {
                        panic(err)
                    }

                    for key, val := range hcdata {
                        tstmp := val.(map[string]interface{})["t"]
                        i, err := strconv.ParseInt(tstmp.(string), 10, 64)
                        if err != nil {
                            panic(err)
                        }
                        tm := time.Unix(i, 0)
                        fmt.Printf("%s -> %s\n",key,tm)
                    }
                }

	}
}

