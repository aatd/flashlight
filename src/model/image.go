/*
 * Flashlight - By Asef Alper Tunga Dündar
 *
 * This Appication is something like Instagram for the University of Applied Sciences
 *
 * API version: 1.0.0
 */

package model

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"mime"
	"path/filepath"
)

//////////////////////////////////////////////////////////////////////
//							Model									//
//////////////////////////////////////////////////////////////////////

type Image struct {
	ID            string `json:"_id"`
	Rev           string `json:"_rev"`
	Type          string `json:"type"`
	Owner         string `json:"owner"`
	ImageLocation string `json:"imagelocation"`
	Description   string `json:"description"`
	Likes         int    `json:"likes,omitempty"`
	Comments      int    `json:"comments,omitempty"`
	UploadTime    int64  `json:"uploadtime"`
	MimeType      string `json:"mimetype"`
}

//////////////////////////////////////////////////////////////////////
//							Methods									//
//////////////////////////////////////////////////////////////////////

// DeleteImage
func (user User) DeleteImage(imageID string) (err error) {

	//Delete all Concering Likes
	likeQuery := fmt.Sprintf(`{
		"selector": {
		   "type": "Like",
		   "imageid":"%s"
		}
	}`, imageID)
	likes, err := DB.QueryJSON(likeQuery)
	if err != nil {
		return err
	}
	for i := 0; i < len(likes); i++ {
		id := likes[i]["_id"].(string)
		err = DeleteLike(id)
		if err != nil {
			return err
		}
	}

	//Delete all Concering Comments
	commentQuery := fmt.Sprintf(`{
		"selector": {
		   "type": "Comment",
		   "imageid":"%s"
		}
	}`, imageID)
	comments, err := DB.QueryJSON(commentQuery)
	if err != nil {
		return err
	}
	for i := 0; i < len(comments); i++ {
		id := comments[i]["_id"].(string)
		err = DeleteComment(id)
		if err != nil {
			return err
		}
	}

	err = DB.Delete(imageID)

	return err

}

// GetImage
func GetImage(imageID string) (data []byte, mimeType string, err error) {

	//Get Imag from DB
	imageModel, err := DB.Get(imageID, nil)
	if err != nil {
		return nil, "", err
	}

	data, err = DB.GetAttachment(imageModel, imageModel["imagelocation"].(string))
	return data, imageModel["mimetype"].(string), err
}

// GetImage
func GetImageMetaData(imageID string) (image Image, err error) {

	//Get Imag from DB
	imageModel, err := DB.Get(imageID, nil)
	if err != nil {
		return Image{}, err
	}

	//Map Image to Golang Struct
	image, err = map2Image(imageModel)
	likeCounter := 5
	image.Likes = likeCounter

	return image, err
}

// GetImages
func GetImageIDs(recordTimeInMilliseconds string) (images []string, err error) {

	//DB Query
	query := fmt.Sprintf(`{
		"selector": {
		   "type": "Image",
		   "uploadtime": {
			  "$lt": %s
		   }
		},
		"fields": [
		   "_id",
		   "uploadtime"
		],
		"sort": [
		   {
			  "uploadtime": "desc"
		   }
		],
		"limit": 2,
		"skip": 0,
		"execution_stats": true
	}`, recordTimeInMilliseconds)

	//Try get the Images list
	imageList, err := DB.QueryJSON(query)
	if err != nil {

		return nil, err

	}

	//Map raw model to Golang Struct
	imageListModel := make([]string, 0)
	for i := 0; i < len(imageList); i++ {
		imageListModel = append(imageListModel, imageList[i]["_id"].(string))
	}

	return imageListModel, err
}

// CreateImage
func (user User) CreateImage(bytes []byte, filename string, description string, uploadtime int64) (err error) {

	//Assemble Image Data
	md5HashInBytes := md5.Sum([]byte(user.Name))
	md5HashedFilename := hex.EncodeToString(md5HashInBytes[:])
	md5HashedFilename += GenerateUUID()                      //Hash filename with UUID so we don't have problems with same filenames
	mimeType := mime.TypeByExtension(filepath.Ext(filename)) //MimeType

	//Save new ImageMetadata
	newImageModel := Image{
		Type:          "Image",
		Owner:         user.Name,
		ImageLocation: md5HashedFilename,
		Description:   description,
		UploadTime:    uploadtime,
		MimeType:      mimeType,
	}

	//Map to Golang Interface
	newImageModelDB := image2Map(newImageModel)
	delete(newImageModelDB, "_id")
	delete(newImageModelDB, "_rev")

	//Save Metadata
	_, _, err = DB.Save(newImageModelDB, nil)
	if err != nil {

		return err

	}

	//Try Put image Data into document
	err = DB.PutAttachment(newImageModelDB, bytes, md5HashedFilename, mimeType)
	if err != nil {

		return err

	}

	//Update Users and save
	//imageIDs := append(user.ImageIDs, newImageId)
	//user.ImageIDs = imageIDs
	//userModel, _ := user2Map(user)
	//_, _, err = DB.Save(userModel, nil)
	//if err != nil {
	//
	//		panic(err)
	//
	//	}

	return err
}

// GetImages
func (user User) GetImages() (images []Image, err error) {

	//DB Query
	query := fmt.Sprintf(`{
			"selector": {
			   "type": "Image",
			   "owner": "%s"
			},
			"sort": [
			   {
				  "uploadtime": "desc"
			   }
			]
		 }`, user.Name)

	imagesMap, err := DB.QueryJSON(query)

	//Append all Images to output
	for i := 0; i < len(imagesMap); i++ {

		image, err := map2Image(imagesMap[i])
		if err != nil {
			return nil, err
		}

		images = append(images, image)
	}

	return images, err

}

// UpdateImage
func (user User) UpdateImage() (err error) {
	return nil
}

//////////////////////////////////////////////////////////////////////
//							Helpers									//
//////////////////////////////////////////////////////////////////////

func image2Map(image Image) map[string]interface{} {
	var doc map[string]interface{}
	tJSON, _ := json.Marshal(image)
	json.Unmarshal(tJSON, &doc)

	return doc
}
func map2Image(image map[string]interface{}) (img Image, err error) {
	uJSON, err := json.Marshal(image)
	json.Unmarshal(uJSON, &img)

	return img, err
}
