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
	"net/http"
	"net/mail"
	"strings"

	"model"

	"github.com/gorilla/mux"
	"golang.org/x/crypto/bcrypt"
)

// GetImages returns a list of id's based on query variable "lastrecordtime".
// These id's containing references to images. These id's are sorted from
// recent to older. The SPA uses these id's to get the model.Image Object through
// model.GetImage.
// Returns HTTP-Status Code 200 Ok and a JSON-Object with imageID's.
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

	//Make Response Model
	responseModel := struct {
		ImagesIDs []string
	}{
		ImagesIDs: imageIDs,
	}

	//Make ResponeJSON
	responseJSON, err := json.Marshal(responseModel)
	if err != nil {
		panic(err)
	}

	//Write response
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	w.Write(responseJSON)

}

// GetImage returns a image file based on a valid id. The respective
// image file is retrieved through the query variable "imageID".
// Returns HTTP-Status Code 404 NotFound when imageID hasn't a
// respective file assoicated with it in the DB.
// Returns HTTP-Status Code 200 Ok and the image file with respecitve
// MIME-Type.
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

// GetImageMetaData returns a model.Image Object based on a valid id. The respective
// model.Image Object is retrieved through the query variable "imageID".
// Returns HTTP-Status Code 404 NotFound when imageID hasn't a
// respective model.Image Obejct assoicated with it in the DB.
// Returns HTTP-Status Code 200 Ok and the JSON-Object with respective model.Image Object
// and associated model.Comments Objects.
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

// GetMainPage returns the MainPage. After called most of the Tasks (e.g. Servercommunication)
// is done by this template with it's respective static files.
// Returns HTTP-Status Code 200 Ok and the HTML-Document with respective template parsed by
// controller.tmpl at it's package intitializaion step at controller.init()
func GetMainPage(w http.ResponseWriter, r *http.Request) {

	//Http Header
	w.Header().Set("Content-Type", "text/html; charset=UTF-8")
	w.WriteHeader(http.StatusOK)

	//For Debug only so the server musnt bne started all over again when .html file is rewirtten again!
	//tmpl = template.Must(template.ParseGlob("src/views/*.html"))

	//Return Webpage
	tmpl.ExecuteTemplate(w, "StartPage", nil)

}

// Login makes an Login attempt to Flashlight. All input are only parsed from posted Form.
// Returns HTTP-Status Code 202 Found and Set-Cookie with respective Authentification
// Session Cookie accepted in the controller.Auth Middleware.
// Returns HTTP-Status Code 409 Conflict and error messages describing derived probelm
// with the Login-Attempt
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

// Register makes an attempt to create a new User. All input is checked wether it already exits
// in DB or credentials are inputted double to verify orignallity of new User.
// Returns HTTP-Status Code 201 Created creates a new model.User-Obejct within DB and Set-Cookie with respective
// authentification Session Cookie accepted in the controller.Auth Middleware.
// Returns HTTP-Status Code 409 Conflict and error messages describing derived problem
// with the Register-Attempt
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
