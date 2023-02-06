package parseableclient

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

// Send builds logs batches and send to Parseable
func Send(functionName string, records []interface{}) {
	url := GetEnv("PARSEABLE_LOG_URL", "")
	username := GetEnv("PARSEABLE_USERNAME", "")
	password := GetEnv("PARSEABLE_PASSWORD", "")

	applicationName := GetEnv("PARSEABLE_APP_NAME", functionName)
	subsystemName := GetEnv("PARSEABLE_SUB_SYSTEM", "logs")
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
			var text string
			record := record.(map[string]interface{})
			timestamp, _ := time.Parse("2006-01-02T15:04:05.000Z", record["time"].(string))

			switch v := record["record"].(type) {
			case string:
				text = string(v)
			default:
				jsonText, _ := json.Marshal(v)
				text = string(jsonText)
			}

			logEntries = append(logEntries, map[string]interface{}{
				"timestamp": timestamp.UnixNano() / 1000000,
				"severity":  GetSeverityLevel(text),
				"text":      text,
				"category":  record["type"],
			})
		}

		logsBulk, _ := json.Marshal(map[string]interface{}{
			"privateKey":      privateKey,
			"applicationName": applicationName,
			"subsystemName":   subsystemName,
			"logEntries":      logEntries,
		})

		client := &http.Client{}
		request, _ := http.NewRequest("POST", url, bytes.NewBuffer(logsBulk))
		request.SetBasicAuth(username, password)
		request.Close = true
		request.Header.Set("Content-Type", "application/json")
		response, err := client.Do(request)
		if err != nil {
			log.Println("Cannot send logs to Parseable:", err)
		} else {
			defer response.Body.Close()
			if response.StatusCode != 200 {
				log.Println("Parseable API failed with code:", response.StatusCode)
			}
		}
	}
}

// GetEnv extract environment variable or default value
func GetEnv(key string, fallback string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		return fallback
	}
	return value
}

// GetSeverityLevel extract serverity from log message
func GetSeverityLevel(message string) int {
	var severity int

	message = strings.ToLower(message)

	switch {
	case strings.Contains(message, "debug"):
		severity = 1
	case strings.Contains(message, "verbose"), strings.Contains(message, "trace"):
		severity = 2
	case strings.Contains(message, "warning"), strings.Contains(message, "warn"):
		severity = 4
	case strings.Contains(message, "error"), strings.Contains(message, "exception"):
		severity = 5
	case strings.Contains(message, "fatal"), strings.Contains(message, "critical"):
		severity = 6
	default:
		severity = 3
	}

	return severity
}
