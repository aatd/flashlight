/*
 * Flashlight - By Asef Alper Tunga Dündar
 *
 * This Appication is something like Instagram for the University of Applied Sciences
 *
 * API version: 1.0.0
 */

package model

import (
	"encoding/json"
	"errors"
	"fmt"
)

//////////////////////////////////////////////////////////////////////
//							Model									//
//////////////////////////////////////////////////////////////////////

type Like struct {
	ID      string `json:"_id"`
	Rev     string `json:"_rev"`
	ImageID string `json:"imageid"`
	Type    string `json:"type"`
	Owner   string `json:"owner"`
}

//////////////////////////////////////////////////////////////////////
//							Methods									//
//////////////////////////////////////////////////////////////////////

// CreateLike    Done
func (user User) CreateLike(imageID string) (err error) {

	//Create Like Model
	likeModel := like2Map(Like{
		ImageID: imageID,
		Owner:   user.Name,
		Type:    "Like",
	})
	delete(likeModel, "_id")
	delete(likeModel, "_rev")

	//Save Like
	_, _, err = DB.Save(likeModel, nil)

	return err

}

// GetLike       Done
func (user User) GetLike(imageID string) (liked bool, likeID string, err error) {

	//init
	liked = true

	//Query for DB
	query := fmt.Sprintf(`{
		"selector": {
		   "type": "Like",
		   "imageid": "%s",
		   "owner": "%s"
		}
	}`, imageID, user.Name)

	//Get from DB to check existince
	likeMaps, err := DB.QueryJSON(query)
	if err != nil || len(likeMaps) == 0 {

		return false, "", err

	}

	//Make all Checks to verify Validaty
	likeModel, err := map2Like(likeMaps[0])
	if err != nil {

		return false, "", err

	}
	if likeModel.Owner != user.Name {

		return false, "", errors.New("Wrong Person likes the Image, are you real YOU, you f**** bot!?")

	}

	return liked, likeModel.ID, err

}

// GetLikeCounts Done
func (image Image) GetLikeCounts() (amount int, err error) {

	//Query for DB
	query := fmt.Sprintf(`{
		"selector": {
		   "type": "Like",
		   "imageid": "%s"
		},
		"fields": [
		   "_id"
		]
	}`, image.ID)

	//Get Like Objects
	likeMaps, err := DB.QueryJSON(query)
	if err != nil {

		return 0, err

	}

	//Just length required
	amount = len(likeMaps)

	return amount, err
}

// DeleteLike    Done
func DeleteLike(likeID string) (err error) {

	err = DB.Delete(likeID)
	return err

}

//////////////////////////////////////////////////////////////////////
//							Helpers									//
//////////////////////////////////////////////////////////////////////

func like2Map(like Like) map[string]interface{} {

	var doc map[string]interface{}
	tJSON, _ := json.Marshal(like)
	json.Unmarshal(tJSON, &doc)

	return doc
}
func map2Like(like map[string]interface{}) (l Like, err error) {

	uJSON, err := json.Marshal(like)
	json.Unmarshal(uJSON, &l)

	return l, err
}
