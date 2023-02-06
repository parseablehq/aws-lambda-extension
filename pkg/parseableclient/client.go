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

package parseableclient

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

// Send builds logs batches and send to Parseable
func Send(functionName string, records []interface{}) {
	url := getEnv("PARSEABLE_LOG_URL", "")
	stream := getEnv("PARSEABLE_LOG_STREAM", functionName)
	username := getEnv("PARSEABLE_USERNAME", "")
	password := getEnv("PARSEABLE_PASSWORD", "")

	applicationName := getEnv("PARSEABLE_APP_NAME", functionName)
	logEntries := []map[string]interface{}{}

	if url == "" {
		log.Fatalln("PARSEABLE_LOG_URL is not set")
	}
	if username == "" {
		log.Fatalln("PARSEABLE_USERNAME is not set")
	}
	if password == "" {
		log.Fatalln("PARSEABLE_PASSWORD is not set")
	}

	if len(records) > 0 {
		for _, record := range records {
			record := record.(map[string]interface{})
			logEntries = append(logEntries, record)
		}
		bulkLogs, err := json.Marshal(logEntries)
		if err != nil {
			log.Println("Cannot marshal log entry:", err)
		}

		client := &http.Client{}
		request, _ := http.NewRequest("POST", url, bytes.NewBuffer(bulkLogs))
		request.SetBasicAuth(username, password)
		request.Header.Set("Content-Type", "application/json")
		request.Header.Set("x-p-meta-application", applicationName)
		request.Header.Set("x-p-stream", stream)
		request.Close = true

		response, err := client.Do(request)
		if err != nil {
			log.Println("Cannot send logs to Parseable:", err)
		} else {
			defer response.Body.Close()
			if response.StatusCode != 200 {
				body, err := ioutil.ReadAll(response.Body)
				if err != nil {
					log.Println("Failed to read Parseable API response:", err)
				}
				log.Printf("Parseable API failed with code: %d and message: %s", response.StatusCode, string(body))
			}
		}
	}
}

// getEnv extract environment variable or default value
func getEnv(key string, fallback string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		return fallback
	}
	return value
}

// getSeverityLevel extract severity from log message
// func getSeverityLevel(record map[string]interface{}) int {
// 	var message string
// 	switch v := record["record"].(type) {
// 	case string:
// 		message = string(v)
// 	default:
// 		jsonText, _ := json.Marshal(v)
// 		message = string(jsonText)
// 	}

// 	var severity int
// 	message = strings.ToLower(message)

// 	switch {
// 	case strings.Contains(message, "debug"):
// 		severity = 1
// 	case strings.Contains(message, "verbose"), strings.Contains(message, "trace"):
// 		severity = 2
// 	case strings.Contains(message, "warning"), strings.Contains(message, "warn"):
// 		severity = 4
// 	case strings.Contains(message, "error"), strings.Contains(message, "exception"):
// 		severity = 5
// 	case strings.Contains(message, "fatal"), strings.Contains(message, "critical"):
// 		severity = 6
// 	default:
// 		severity = 3
// 	}

// 	return severity
// }
