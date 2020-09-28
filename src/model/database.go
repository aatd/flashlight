/*
 * Flashlight - By Asef Alper Tunga Dündar
 *
 * This Appication is something like Instagram for the University of Applied Sciences
 *
 * API version: 1.0.0
 */

package model

import (
	"fmt"

	couchdb "github.com/leesper/couchdb-golang"
)

// DB is the reference to the intilized central Database within Flashlight
var DB *couchdb.Database
var imageDocID string

func init() {

	//Error to panic
	var err error

	//Init DB
	DB, err = couchdb.NewDatabase("http://localhost:5984/flashlight")
	if err != nil {
		panic(err)
	}
	fmt.Printf("DB: \"%s\" is initialized.\n", "flashlight")

	//Init Images Document
	//allImageRef, err := DB.QueryJSON(`{"selector": {"type": "allImages"},"fields": ["_id","_rev"]}`)
	//if err != nil {
	//	panic(err)
	//}
	//imageDocID = allImageRef[0]["_id"].(string)
	//fmt.Printf("DB: ID of where all Images are is: %s.\n", imageDocID)

}
