/*
 * Flashlight - By Asef Alper Tunga Dündar
 *
 * This Appication is something like Instagram for the University of Applied Sciences Cologne
 *
 * API version: 1.0.0
 * Generated by: Swagger Codegen (https://github.com/swagger-api/swagger-codegen.git)
 */

package utils

import (
	"log"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"time"

	"github.com/gorilla/sessions"
)

// Logger
func Logger(inner http.Handler, name string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		inner.ServeHTTP(w, r)

		log.Printf(
			"%s %s %s %s",
			r.Method,
			name,
			r.RequestURI,
			time.Since(start),
		)
	})
}

// LoggerCookie
func LoggerCookie(session *sessions.Session) {
	log.Printf(
		"SessionCookieEvaluation: Auth: %t, Username: %s, UsernameMD5: %s",
		session.Values["authenticated"],
		session.Values["username"],
		session.Values["hashedusername"],
	)
}

// FlushTempFiles
func FlushTempFiles() {

	workingDir, _ := os.Getwd()

	tempDir := workingDir + "/temp"
	if runtime.GOOS == "windows" {
		tempDir = workingDir + "\\temp"
	}
	allFilenames, _ := filePathWalkDir(tempDir)
	for i := 0; i < len(allFilenames); i++ {
		os.Remove(allFilenames[i])
	}

}

// MakeTimeStamp
func MakeTimeStamp() int64 {
	return time.Now().UnixNano() / int64(time.Millisecond)
}

func filePathWalkDir(root string) ([]string, error) {
	var files []string
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			files = append(files, path)
		}
		return nil
	})
	return files, err
}
