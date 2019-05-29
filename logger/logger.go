package logger

import (
	"context"
	"fmt"
	"log"
	"path/filepath"
	"runtime"
	"time"

	"sanservices.git.beanstalkapp.com/goproposal.git/settings"
)

// LogError Prints error and line where error occurred
func LogError(ctx context.Context, msg string) {

	var errToDisplay string
	_, file, _, _ := runtime.Caller(1)
	log.Printf("\033[31m Error on %s: \033[39m", append([]interface{}{filepath.Base(file)})...)

	// validate if context was passed
	if ctx != nil {

		// get variables to track request
		val := ctx.Value(settings.RequestTracking)

		//assert data is a map of strings
		data, ok := val.(map[string]string)

		if ok {
			apikey := data["api-key"]
			requestID := data["Kong-Request-ID"]
			errToDisplay = "\trequestId=\"%s\", apiKey=\"%s', %s\n"
			fmt.Printf(errToDisplay, requestID, apikey, msg)
		} else {
			fmt.Printf("\t %s\n", msg)
		}

	} else {
		fmt.Printf("\t %s\n", msg)
	}

}

// TrackExecutionTime logs time elapsed for caller function
func TrackExecutionTime(ctx context.Context, caller string, t time.Time) {
	elapsed := time.Since(t)

	if caller == "" {
		fmt.Println("TrackExecutionTime Caller is nil...")
	}

	if ctx == nil {
		log.Println(fmt.Sprintf("Function=\"%s\" ms=\"%f\"", caller, float64(elapsed.Nanoseconds())/float64(time.Millisecond)))
	} else {
		msg := fmt.Sprintf("Function=\"%s\" ms=\"%f\"", caller, float64(elapsed.Nanoseconds())/float64(time.Millisecond))
		Log(ctx, msg)
	}
}

//Log logs standard message
func Log(ctx context.Context, msg string) {

	if ctx == nil {
		log.Println(msg + "\n")
	}

	// get variables to track request
	val := ctx.Value(settings.RequestTracking)

	//assert data is a map of strings
	data, ok := val.(map[string]string)

	if ok {
		apikey := data["api-key"]
		requestID := data["Kong-Request-ID"]
		log.Printf("requestId=\"%s\", apiKey=\"%s\", %s\n", requestID, apikey, msg)
	} else {
		LogError(nil, "Invalid request-tracking value")
		log.Println(msg + "\n")
	}
}
