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
	"fmt"
	"utils"
)

//////////////////////////////////////////////////////////////////////
//							Model									//
//////////////////////////////////////////////////////////////////////

// Comment is the Model which stores all required information to
// contain a literal comment for a model.Image. For simplicity all cross
// references to other Flashlight Models (e.g. models.Image) are only stored
// to their respective "_id"
type Comment struct {
	ID         string `json:"_id"`
	Rev        string `json:"_rev"`
	ImageID    string `json:"imageid"`
	Type       string `json:"type"`
	UploadTime int64  `json:"uploadtime"`
	Owner      string `json:"owner"`
	Comment    string `json:"comment"`
}

//////////////////////////////////////////////////////////////////////
//							Methods									//
//////////////////////////////////////////////////////////////////////

// GetComments
func (image Image) GetComments() (comments []Comment, err error) {

	//Query for getting all Comments of an Image
	query := fmt.Sprintf(`{"selector":{"type":"Comment","imageid":"%s","uploadtime":{"$gt":%d}},"sort": [{"uploadtime":"desc"}]}`, image.ID, 0)

	//Get all Comments
	dbComments, err := DB.QueryJSON(query)
	if err != nil {

		return nil, err

	}

	//Put them into Golang Struct
	for i := 0; i < len(dbComments); i++ {

		dbComment, err := map2Comment(dbComments[i])
		if err != nil {

			return nil, err

		}
		comments = append(comments, dbComment)

	}

	return comments, err
}

// CreateComment
func (user User) CreateComment(imageID string, comment string) (err error) {

	//Create Comment and map to Golang Interface
	commentModel := Comment{
		Type:       "Comment",
		Owner:      user.Name,
		Comment:    comment,
		UploadTime: utils.MakeTimeStamp(),
		ImageID:    imageID,
	}

	//Save comment to DB
	commentMap := comment2Map(commentModel)
	delete(commentMap, "_id")
	delete(commentMap, "_rev")

	_, _, err = DB.Save(commentMap, nil)

	return err
}

// DeleteComment
func DeleteComment(commentID string) (err error) {

	err = DB.Delete(commentID)
	return err
}

// UpdateComment
func (user User) UpdateComment(comment Comment) (err error) {
	return nil
}

//////////////////////////////////////////////////////////////////////
//							Helpers									//
//////////////////////////////////////////////////////////////////////

func comment2Map(comment Comment) map[string]interface{} {
	var doc map[string]interface{}
	tJSON, _ := json.Marshal(comment)
	json.Unmarshal(tJSON, &doc)

	return doc
}
func map2Comment(comment map[string]interface{}) (c Comment, err error) {
	uJSON, err := json.Marshal(comment)
	json.Unmarshal(uJSON, &c)

	return c, err
}
