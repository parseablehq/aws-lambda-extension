// Copyright 2023 Cloudnatively Pvt. Ltd. All rights reserved.

// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at

//     http://www.apache.org/licenses/LICENSE-2.0

// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"
	"strconv"
	"time"

	"github.com/parseablehq/aws-lambda-extension/pkg/extensionsclient"
	"github.com/parseablehq/aws-lambda-extension/pkg/logsclient"
	"github.com/parseablehq/aws-lambda-extension/pkg/parseableclient"
)

var (
	// ReceiverHost describes logs receiver hostname
	ReceiverHost = "0.0.0.0"
	// ReceiverPort describes logs receiver port
	ReceiverPort = 4342
)

func main() {
	agentName := path.Base(os.Args[0])
	listenerState := make(chan bool)
	queue := make(chan interface{})

	log.Println("Initializing Lambda Extension", agentName)
	agentID, functionName, err := extensionsclient.Register(agentName)
	if err != nil {
		log.Fatalln("Failed to register Lambda Extension", agentName)
	}

	go serve(queue, listenerState)
	select {
	case <-listenerState:
		logsclient.Subscribe(agentID.(string), map[string]interface{}{
			"destination": map[string]interface{}{
				"protocol": "HTTP",
				"URI":      fmt.Sprintf("http://sandbox:%d", ReceiverPort),
			},
			"types": []string{"platform", "function"},
			"buffering": map[string]uint{
				"timeoutMs": 1000,
				"maxBytes":  1048576,
				"maxItems":  10000,
			},
		})

		for {
			extensionsclient.Next(agentID.(string))
			parseableclient.Send(functionName.(string), (<-queue).([]interface{}))
		}
	case <-time.After(9 * time.Second):
		log.Fatalln("HTTP Server time out before starting")
	}
}

// serve start HTTP server to accept incoming events
func serve(queue chan<- interface{}, state chan bool) {
	address := fmt.Sprintf("%s:%d", ReceiverHost, ReceiverPort)

	log.Println("Initializing HTTP Server on", address)
	http.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		if request.Method == http.MethodPost {
			var contentLength int
			var requestBody interface{}

			contentLength, err := strconv.Atoi(request.Header.Get("Content-Length"))
			if err != nil {
				contentLength = 0
			}

			if contentLength > 0 {
				defer request.Body.Close()
				requestBytes, err := ioutil.ReadAll(request.Body)
				if err != nil {
					panic(err)
				}

				err = json.Unmarshal(requestBytes, &requestBody)
				if err != nil {
					panic(err)
				}

				queue <- requestBody
			}

			writer.WriteHeader(http.StatusOK)
			return
		}
	})

	log.Println("Serving HTTP Server on", address)
	state <- true
	http.ListenAndServe(address, nil)
}
