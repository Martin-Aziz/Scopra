package utils

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

type Logger struct {
	service string
}

func NewLogger(service string) Logger {
	return Logger{service: service}
}

func (l Logger) Info(message string, fields map[string]any) {
	l.log("info", message, fields)
}

func (l Logger) Warn(message string, fields map[string]any) {
	l.log("warn", message, fields)
}

func (l Logger) Error(message string, fields map[string]any) {
	l.log("error", message, fields)
}

func (l Logger) log(level, message string, fields map[string]any) {
	payload := map[string]any{
		"timestamp": time.Now().UTC().Format(time.RFC3339Nano),
		"level":     level,
		"service":   l.service,
		"message":   message,
	}

	for key, value := range fields {
		payload[key] = value
	}

	encoded, err := json.Marshal(payload)
	if err != nil {
		fmt.Fprintf(os.Stderr, "{\"level\":\"error\",\"message\":\"failed to encode log\",\"error\":%q}\n", err.Error())
		return
	}

	fmt.Fprintln(os.Stdout, string(encoded))
}
