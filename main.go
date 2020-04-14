/*
 * Flashlight - By Asef Alper Tunga Dündar
 *
 * This Appication is something like Instagram for the University of Applied Sciences
 *
 * API version: 1.0.0
 */

package main

import (
	"log"
	"net/http"

	"controller"
)

func main() {
	log.Printf("Server started")
	router := controller.NewRouter()
	log.Fatal(http.ListenAndServe(":8080", router))
}
