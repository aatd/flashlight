/*
 * Flashlight - By Asef Alper Tunga Dündar
 *
 * This Appication is something like Instagram for the University of Applied Sciences
 *
 * API version: 1.0.0
 */

package model

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

//////////////////////////////////////////////////////////////////////
//							Model									//
//////////////////////////////////////////////////////////////////////

// User
type User struct {
	ID       string `json:"_id"`
	Rev      string `json:"_rev"`
	Type     string `json:"type"`
	Name     string `json:"name"`
	Password string `json:"password"`
	Email    string `json:"email"`
}

//////////////////////////////////////////////////////////////////////
//							Metohds									//
//////////////////////////////////////////////////////////////////////

// CreateUser
func (user User) CreateUser() (err error) {

	// Check wether username already exists
	userInDBByName, err := GetUserByUsername(user.Name)
	if err == nil && userInDBByName.Name == user.Name {
		return errors.New("Benutzername exisitert bereits.")
	}

	// Check wether username already exists
	userInDBByEmail, err := GetUserByEmail(user.Email)
	if err == nil && userInDBByEmail.Email == user.Email {
		return errors.New("Emailadresse bereits vergeben.")
	}

	//Check if
	if userInDBByEmail.ID != userInDBByName.ID {
		return errors.New("Benutzername und Emailadresse bereits vergeben.")
	}

	// Hash password and set on User
	hashedPwd, err := bcrypt.GenerateFromPassword([]byte(user.Password), 14)
	b64HashedPwd := base64.StdEncoding.EncodeToString(hashedPwd)
	user.Password = b64HashedPwd
	if err != nil {
		return err
	}

	// Convert to Golang Map
	userModel, err := user2Map(user)
	if err != nil {
		return err
	}

	// Delete _id and _rev from map, otherwise DB access will be denied (unauthorized)
	delete(userModel, "_id")
	delete(userModel, "_rev")

	// Add User to DB
	_, _, err = DB.Save(userModel, nil)
	if err != nil {
		return err
	}

	return err
}

// GetUserByUsername
func GetUserByUsername(username string) (user User, err error) {

	// Check before annoying DB
	if username == "" {
		return User{}, errors.New("No username provided")
	}

	// Get User
	query := fmt.Sprintf(`{"selector": {"type": "User","name":"%s"}}`, username)
	u, err := DB.QueryJSON(query)
	if err != nil || len(u) != 1 {
		return User{}, errors.New("Benutzername existiert nicht.")
	}

	// Convert to Golang Map
	user, err = map2User(u[0])
	if err != nil {
		return User{}, err
	}

	return user, err
}

// GetUserByEmail
func GetUserByEmail(username string) (user User, err error) {

	// Check before annoying DB
	if username == "" {
		return User{}, errors.New("No email provided")
	}

	// Get User
	query := fmt.Sprintf(`{"selector": {"type": "User","email":"%s"}}`, username)
	u, err := DB.QueryJSON(query)
	if err != nil || len(u) != 1 {
		return User{}, err
	}

	// Convert to Golang Map
	user, err = map2User(u[0])
	if err != nil {
		return User{}, err
	}

	return user, nil
}

//////////////////////////////////////////////////////////////////////
//							Helpers									//
//////////////////////////////////////////////////////////////////////

func user2Map(u User) (user map[string]interface{}, err error) {
	uJSON, err := json.Marshal(u)
	json.Unmarshal(uJSON, &user)

	return user, err
}
func map2User(user map[string]interface{}) (u User, err error) {
	uJSON, err := json.Marshal(user)
	json.Unmarshal(uJSON, &u)

	return u, err
}
