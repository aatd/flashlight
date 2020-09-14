/*
 * Flashlight - By Asef Alper Tunga Dündar
 *
 * This Appication is something like Instagram for the University of Applied Sciences
 *
 * API version: 1.0.0
 */

package controller

import (
	"encoding/json"
	"io/ioutil"
	"model"
	"net/http"
	"os"
	"runtime"
	"strconv"
	"utils"

	"github.com/gorilla/mux"
)

// UploadImage let's a user with a valid session uploads a images with the respective metadata
// such as it's description being uploaded to the DB. UploadeImage also creates it's respective mode.Image
// Object stored wihting the DB to get properly allocated for all futher methods.
// Returns HTTP-Status Code 409 Conflict when wrong data is submitted.
// Returns HTTP-Status Code 201 Created when everything went okey.
func UploadImage(w http.ResponseWriter, r *http.Request) {

	//Get current Session
	session, _ := store.Get(r, "session")
	name := session.Values["username"].(string)

	//Get User uploading Image
	user, err := model.GetUserByUsername(name)
	if err != nil {

		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte(err.Error()))
		return

	}

	//Get FormFile
	r.ParseMultipartForm(10 << 20)
	file, handler, err := r.FormFile("uploadfile")
	defer file.Close()
	if err != nil {

		w.WriteHeader(http.StatusConflict)
		w.Write([]byte(err.Error()))
		return

	}

	//Get All Other Data from Form
	filename := handler.Filename                                           //Imagename
	description := r.FormValue("description")                              //ImageDescription
	uploadtime, err := strconv.ParseInt(r.FormValue("uploadtime"), 10, 64) //Upload Time in Millisesconds
	if err != nil {

		w.WriteHeader(http.StatusConflict)
		w.Write([]byte(err.Error()))
		return

	}

	//Create Temp File
	workingDir, _ := os.Getwd()
	tempDir := workingDir + "/temp"
	if runtime.GOOS == "windows" {
		tempDir = workingDir + "\\temp"
	}

	tempFile, err := ioutil.TempFile(tempDir, "upload-*.png")
	if err != nil {

		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return

	}
	defer tempFile.Close()

	//Write to File to Upload
	imagebytes, err := ioutil.ReadAll(file)
	if err != nil {

		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return

	}
	tempFile.Write(imagebytes)

	//Try Uploading Image with all Metadata
	err = user.CreateImage(imagebytes, filename, description, uploadtime)
	if err != nil {

		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return

	}

	//Delete Tempfile
	utils.FlushTempFiles()

	w.WriteHeader(http.StatusCreated)
}

// Logout updates the servers internal session store to not authentificated
// and detstroies the respective client's cookie.
// Returns HTTP-Status Code 409 Conflict when cookie couldn't be read.
// Returns Set-Cookie: {Flashlight-SessionCookie} injected Response
func Logout(w http.ResponseWriter, r *http.Request) {

	//Empty Session
	session, err := store.Get(r, "session")
	session.Values["authenticated"] = false
	session.Values["username"] = ""
	session.Values["hashedusername"] = ""
	session.Options.MaxAge = -1
	session.Save(r, w)
	if err != nil {

		w.WriteHeader(http.StatusConflict)
		w.Write([]byte(err.Error()))

	}

}

// Userdata returns metadata over a loggedin user with a valid usersession.
// No security issued data is returned wihtin this statement. Only metadata providing
// the SPA served troughout "/" or "/users/{userid}" are to be filled correctly.
// returns HTTP-Status Code 200 Ok and the Metadata with user information as "application/json".
func Userdata(w http.ResponseWriter, r *http.Request) {

	session, err := store.Get(r, "session")
	if err != nil {

		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte(err.Error()))

	}

	//Make Response JSON
	responseModel := struct {
		Username       string
		HashedUsername string
	}{
		Username:       session.Values["username"].(string),
		HashedUsername: session.Values["hashedusername"].(string),
	}
	responseJSON, err := json.Marshal(responseModel)
	if err != nil {

		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))

	}

	//Write response
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	w.Write(responseJSON)

}

// ProfilePage returns the Mainpage when logged in. Due to the Front-End Application is held as
// a SPA actual Session management is also made within the pages logic. This Handlers purpose is
// mainly when users first navigation in over "/users/{userid}" instead of "/"
func ProfilePage(w http.ResponseWriter, r *http.Request) {

	//Http Header
	w.Header().Set("Content-Type", "text/html; charset=UTF-8")
	w.WriteHeader(http.StatusOK)

	//For Debug only so the server musnt bne started all over again when .html file is rewirtten again!
	//tmpl = template.Must(template.ParseGlob("src/views/*.html"))

	tmpl.ExecuteTemplate(w, "StartPage", nil)

}

// CommentImage comments an Image with respective id and authentificated usersession. Over these
// informations CommentImage extraces the user commenting the image and saves these up wihtin the
// used DB.
func CommentImage(w http.ResponseWriter, r *http.Request) {

	//Get current Session
	session, _ := store.Get(r, "session")
	name := session.Values["username"].(string)

	//Pathparameter
	vars := mux.Vars(r)
	imageID := vars["imageID"]

	//Get User Commenting Image
	user, err := model.GetUserByUsername(name)
	if err != nil {

		w.WriteHeader(http.StatusConflict)
		w.Write([]byte(err.Error()))
		return

	}

	//Get Formdata
	comment := r.FormValue("comment")

	//Create Comment to DB
	err = user.CreateComment(imageID, comment)
	if err != nil {

		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return

	}

	//Make Response JSON
	responseModel := struct {
		Message string
	}{Message: "Image was Commented Sussesfully!"}
	responseJSON, err := json.Marshal(responseModel)
	if err != nil {

		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	////Write response
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusCreated)
	w.Write(responseJSON)
}

// DeleteImage deletes an image with respective id and authentificated usersession. This method only
// works when the the image is wihting the set of the user. If not or anything else's not conforming
// withing the specs DeleteImage responses with the HTTP-Status Code 409 Conflict.
// DeleteImage returns HTTP-Status Code 202 Accepted when everything went okey.
func DeleteImage(w http.ResponseWriter, r *http.Request) {

	//Get current Session
	session, _ := store.Get(r, "session")
	name := session.Values["username"].(string)

	//Pathparameter
	vars := mux.Vars(r)
	imageID := vars["imageID"]

	//Get User Commenting Image
	user, err := model.GetUserByUsername(name)
	if err != nil {

		w.WriteHeader(http.StatusConflict)
		w.Write([]byte(err.Error()))
		return

	}

	//Delete Image from Database
	user.DeleteImage(imageID)

	//Make Response JSON
	responseModel := struct {
		Message string
	}{Message: "Image was deleted Sussesfully!"}
	responseJSON, err := json.Marshal(responseModel)
	if err != nil {

		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return

	}

	//Write response
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusAccepted)
	w.Write(responseJSON)
}

// LikeImage adds or removes a like to a image with respective id and binds it to a authentificated usersession.
// It may add a like when user has not liked the image before and removes it else. Returns HTTP-Status Code 202
// Accepted when the Like Object was deleted due to users unliking it else returns HTTP-Status Code 201 Created
// when liked the respecitve image.
func LikeImage(w http.ResponseWriter, r *http.Request) {

	type ResponseModel struct {
		Message string
		ImageID string
	}

	//Get current Session
	session, _ := store.Get(r, "session")
	name := session.Values["username"].(string)

	//Pathparameter
	vars := mux.Vars(r)
	imageID := vars["imageID"]

	//Get User Commenting Image
	user, err := model.GetUserByUsername(name)
	if err != nil {

		w.WriteHeader(http.StatusConflict)
		w.Write([]byte(err.Error()))
		return

	}

	//Creat or Delete Like depending on users status
	isAlreadyLiked, likeID, err := user.GetLike(imageID)
	if err != nil {

		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return

	}

	if isAlreadyLiked {

		err = model.DeleteLike(likeID)
		if err != nil {

			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return

		}

		//Make Response JSON
		responseModel := ResponseModel{
			Message: "Image was Unliked Sussesfully!",
			ImageID: imageID,
		}
		responseJSON, err := json.Marshal(responseModel)
		if err != nil {

			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return

		}

		//Write response
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusAccepted)
		w.Write(responseJSON)

	} else {

		//Create Like to DB
		err = user.CreateLike(imageID)
		if err != nil {

			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return

		}

		//Make Response JSON
		responseModel := ResponseModel{
			Message: "Image was Unliked Sussesfully!",
			ImageID: imageID,
		}
		responseJSON, err := json.Marshal(responseModel)
		if err != nil {

			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return

		}

		//Write response
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusCreated)
		w.Write(responseJSON)

	}
}

// GetUserImages returns all images a authentificated user with respective valid usersession has uploaded to
// Flashlight. Returns HTTP-Status Code 200 Ok.
func GetUserImages(w http.ResponseWriter, r *http.Request) {

	//Get current Session
	session, _ := store.Get(r, "session")
	name := session.Values["username"].(string)

	//Get User
	user, err := model.GetUserByUsername(name)
	if err != nil {

		w.WriteHeader(http.StatusConflict)
		w.Write([]byte(err.Error()))
		return

	}

	//Get Images
	images, err := user.GetImages()
	if err != nil {

		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return

	}

	//Get like and comment counts for each Image
	for i := 0; i < len(images); i++ {

		images[i].Likes, err = images[i].GetLikeCounts()
		if err != nil {

			http.Error(w, err.Error(), http.StatusInternalServerError)
			return

		}
		comments, err := images[i].GetComments()
		if err != nil {

			http.Error(w, err.Error(), http.StatusInternalServerError)
			return

		}
		images[i].Comments = len(comments)

	}

	//Make Response JSON
	responseModel := struct {
		Images []model.Image
	}{
		Images: images,
	}
	responseJSON, err := json.Marshal(responseModel)
	if err != nil {

		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return

	}

	//Write response
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	w.Write(responseJSON)

}

// GetLike gives a respective valid usersession validation for a specific image with id located wihtin the
// query variable "imageID" is liked or not.
// Returns HTTP-Status Code 200 Ok and a JSON-Object if usersession and imageID is valid.
// Returns HTTP-Status Code 409 Conflict if not.
func GetLike(w http.ResponseWriter, r *http.Request) {

	type ResponseModel struct {
		ImageID string
		IsLiked bool
	}

	//Get current Session
	session, _ := store.Get(r, "session")
	name := session.Values["username"].(string)

	//Pathparameter
	vars := mux.Vars(r)
	imageID := vars["imageID"]

	//Get User Commenting Image
	user, err := model.GetUserByUsername(name)
	if err != nil {

		w.WriteHeader(http.StatusConflict)
		w.Write([]byte(err.Error()))
		return

	}

	//Creat or Delete Like depedngin on users status
	isLiked, _, err := user.GetLike(imageID)
	if err != nil {

		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return

	}

	//Make Response JSON
	responseModel := ResponseModel{
		ImageID: imageID,
		IsLiked: isLiked,
	}
	responseJSON, err := json.Marshal(responseModel)
	if err != nil {

		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return

	}

	//Write response
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	w.Write(responseJSON)

}
