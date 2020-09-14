/*
 * Flashlight - By Asef Alper Tunga Dündar
 *
 * This Appication is something like Instagram for the University of Applied Sciences
 *
 * API version: 1.0.0
 */

package controller

import (
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"html/template"
	"net/http"
	"net/mail"
	"strings"

	"model"

	"github.com/gorilla/mux"
	"golang.org/x/crypto/bcrypt"
)

// GetImages
func GetImages(w http.ResponseWriter, r *http.Request) {

	//Get Formdata
	queryVar := mux.Vars(r)
	lastrecordedtime := queryVar["[0-9]+"]

	//Get ImageIDs
	imageIDs, err := model.GetImageIDs(lastrecordedtime)
	if err != nil {

		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return

	}

	/*
		//Get Images
		images := make([]model.Image, 0)
		for i := 0; i < len(imagesIDs); i++ {

			image, err := model.GetImageMetaData(imagesIDs[i])
			if err != nil {

				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(err.Error()))
				return

			}
			images = append(images, image)

		}*/

	//Make Response JSON
	responseModel := struct {
		ImagesIDs []string
	}{
		ImagesIDs: imageIDs,
	}

	//Make Response JSON
	//responseModel := struct {
	//	LastRecordedTime int64
	//	Images           []model.Image
	//}{
	//	LastRecordedTime: images[len(images)-1].UploadTime,
	//	Images:           images,
	//}
	responseJSON, err := json.Marshal(responseModel)
	if err != nil {
		panic(err)
	}

	//Write response
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	w.Write(responseJSON)

}

// GetImage
func GetImage(w http.ResponseWriter, r *http.Request) {

	//Response Parameter
	vars := mux.Vars(r)
	imageID := vars["imageID"]

	//Get Data and make Response
	image, mimeType, err := model.GetImage(imageID)
	if err != nil {

		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(err.Error()))
		return

	}

	//Write Response
	w.Header().Set("Content-Type", mimeType)
	w.WriteHeader(http.StatusOK)
	w.Write(image)

}

// GetImageMetaData
func GetImageMetaData(w http.ResponseWriter, r *http.Request) {

	type ResponseModel struct {
		ImageMetaData model.Image
		Comments      []model.Comment
	}

	//Response Parameter
	vars := mux.Vars(r)
	imageID := vars["imageID"]

	//Get Data and make Response
	image, err := model.GetImageMetaData(imageID)
	if err != nil {

		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(err.Error()))
		return

	}

	image.Likes, err = image.GetLikeCounts()
	if err != nil {

		http.Error(w, err.Error(), http.StatusInternalServerError)
		return

	}

	comments, err := image.GetComments()
	if err != nil {

		http.Error(w, err.Error(), http.StatusInternalServerError)
		return

	}

	repsoneModel := ResponseModel{
		ImageMetaData: image,
		Comments:      comments,
	}

	//Create JSON
	responseJSON, err := json.Marshal(repsoneModel)
	if err != nil {

		http.Error(w, err.Error(), http.StatusInternalServerError)
		return

	}

	//Write Response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(responseJSON)

}

// GetMainPage
func GetMainPage(w http.ResponseWriter, r *http.Request) {

	//Http Header
	w.Header().Set("Content-Type", "text/html; charset=UTF-8")
	w.WriteHeader(http.StatusOK)

	//For Debug only so the server musnt bne started all over again when .html file is rewirtten again!
	tmpl = template.Must(template.ParseGlob("src/views/*.html"))

	//Return Webpage
	tmpl.ExecuteTemplate(w, "StartPage", nil)

}

// Login
func Login(w http.ResponseWriter, r *http.Request) {

	//Create Error messages
	messages := make([]string, 0)
	type MultiErrorMessages struct {
		Messages []string
	}

	// Get Formdata
	var user = model.User{}
	username := r.FormValue("username")
	password := r.FormValue("password")

	// Try Authentification
	user, err := model.GetUserByUsername(username)
	if err != nil {

		//Add error Message
		messages = append(messages, "Benutzer existiert nicht.")

	}

	//Encode Password to base64 byte array
	passwordDB, _ := base64.StdEncoding.DecodeString(user.Password)
	err = bcrypt.CompareHashAndPassword(passwordDB, []byte(password))
	if err != nil {

		//Add error Message
		messages = append(messages, "Password falsch.")

	}

	//Check if any Error Message was assembled
	if len(messages) != 0 {

		responseModel := MultiErrorMessages{
			Messages: messages,
		}

		responseJSON, err := json.Marshal(responseModel)
		if err != nil {

			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))

		}

		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusConflict)
		w.Write(responseJSON)
		return

	}

	//Hash Username
	md5HashInBytes := md5.Sum([]byte(user.Name))
	md5HashedUsername := hex.EncodeToString(md5HashInBytes[:])

	//Create Session
	session, err := store.Get(r, "session")
	session.Values["authenticated"] = true
	session.Values["username"] = username
	session.Values["hashedusername"] = md5HashedUsername
	if err != nil {

		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))

	}

	//Save Session
	err = session.Save(r, w)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
	}

	//Redirect to ProfilePage
	http.Redirect(w, r, "/users?action=userdata", http.StatusFound)
}

// Register
func Register(w http.ResponseWriter, r *http.Request) {

	messages := make([]string, 0)
	type MultiErrorMessages struct {
		Messages []string
	}

	//Get Formdata
	username := r.FormValue("username")
	password := r.FormValue("password")
	emailadress := r.FormValue("email")
	repeatPassword := r.FormValue("repeatpassword")

	//Check Password
	if password != repeatPassword {

		//Add error Message
		messages = append(messages, "Passwort ist nicht richtig wiedeholt worden.")

	}

	//Check Email
	email, err := mail.ParseAddress(emailadress)
	if err != nil || !strings.Contains(email.Address, ".") {

		//Add error Message
		messages = append(messages, "Dies ist keine gültige Emailadresse.")

	}

	//Fill Model
	user := model.User{}
	user.Name = username
	user.Password = password
	user.Email = emailadress
	user.Type = "User"

	//Try and check Creating User
	err = user.CreateUser()
	if err != nil {

		//Write Data
		messages = append(messages, err.Error())

	}

	//Check if any Error Message was assembled
	if len(messages) != 0 {

		responseModel := MultiErrorMessages{
			Messages: messages,
		}

		responseJSON, err := json.Marshal(responseModel)
		if err != nil {

			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))

		}

		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusConflict)
		w.Write(responseJSON)
		return
	}

	//Hash Username
	md5HashInBytes := md5.Sum([]byte(user.Name))
	md5HashedUsername := hex.EncodeToString(md5HashInBytes[:])

	//Create Session
	session, _ := store.Get(r, "session")
	session.Values["authenticated"] = true
	session.Values["username"] = username
	session.Values["hashedusername"] = md5HashedUsername
	session.Save(r, w)

	//Write Respone
	http.Redirect(w, r, "/users?action=userdata", http.StatusFound)
}
