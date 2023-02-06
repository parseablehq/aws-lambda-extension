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

package extensionsclient

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

// Based on the AWS Lambda Extension API docs here
// https://docs.aws.amazon.com/lambda/latest/dg/runtimes-extensions-api.html

// Register creates registration of lambda extension
func Register(agentName string) (interface{}, interface{}, error) {
	runtimeAPIAddress, exists := os.LookupEnv("AWS_LAMBDA_RUNTIME_API")
	if !exists {
		return nil, nil, errors.New("AWS_LAMBDA_RUNTIME_API is not set")
	}

	registrationBodyJSON, err := json.Marshal(map[string]interface{}{
		"events": []string{"INVOKE", "SHUTDOWN"},
	})
	if err != nil {
		return nil, nil, err
	}

	log.Println("Registering to Extensions API")

	client := &http.Client{}
	request, _ := http.NewRequest(
		"POST",
		fmt.Sprintf("http://%s/2020-01-01/extension/register", runtimeAPIAddress),
		bytes.NewBuffer(registrationBodyJSON),
	)
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Lambda-Extension-Name", agentName)
	response, err := client.Do(request)
	if err != nil {
		return nil, nil, err
	}

	defer response.Body.Close()
	responseBytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, nil, err
	}
	responseBody := string(responseBytes)

	if response.StatusCode != 200 {
		log.Fatalln("Could not register to Extensions API:", response.StatusCode, responseBody)
	}

	var responseJSON interface{}
	err = json.Unmarshal(responseBytes, &responseJSON)
	if err != nil {
		return nil, nil, err
	}
	functionName := responseJSON.(map[string]interface{})["functionName"]

	return response.Header.Get("Lambda-Extension-Identifier"), functionName, nil
}

// Next pulls lambda extension next event
func Next(agentID string) (interface{}, error) {
	runtimeAPIAddress, exists := os.LookupEnv("AWS_LAMBDA_RUNTIME_API")
	if !exists {
		return nil, errors.New("AWS_LAMBDA_RUNTIME_API is not set")
	}

	client := &http.Client{}
	request, _ := http.NewRequest(
		"GET",
		fmt.Sprintf("http://%s/2020-01-01/extension/event/next", runtimeAPIAddress),
		nil,
	)
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Lambda-Extension-Identifier", agentID)
	response, err := client.Do(request)
	if err != nil {
		return nil, err
	}

	defer response.Body.Close()
	responseBytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	responseBody := string(responseBytes)

	if response.StatusCode != 200 {
		log.Fatalln("Request to Extensions API failed:", response.StatusCode, responseBody)
	}

	//log.Println("Received response from Extensions API:", responseBody)

	return responseBody, nil
}
