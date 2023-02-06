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

package logsclient

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

// Subscribe creates subscription to Logs API
func Subscribe(agentID string, subscriptionBody interface{}) error {
	runtimeAPIAddress, exists := os.LookupEnv("AWS_LAMBDA_RUNTIME_API")
	if !exists {
		return errors.New("AWS_LAMBDA_RUNTIME_API is not set")
	}

	subscriptionBodyJSON, err := json.Marshal(subscriptionBody)
	if err != nil {
		return err
	}

	log.Println("Subscribing to Logs API")

	client := &http.Client{}
	request, _ := http.NewRequest(
		"PUT",
		fmt.Sprintf("http://%s/2020-08-15/logs", runtimeAPIAddress),
		bytes.NewBuffer(subscriptionBodyJSON),
	)
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Lambda-Extension-Identifier", agentID)
	response, err := client.Do(request)
	if err != nil {
		return err
	}

	defer response.Body.Close()
	responseBytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return err
	}
	responseBody := string(responseBytes)

	if response.StatusCode != 200 {
		log.Fatalln("Could not subscribe to Logs API:", response.StatusCode, responseBody)
	}

	log.Println("Successfully subscribed to Logs API:", responseBody)
	return nil
}
