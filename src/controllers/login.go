package controllers

import (
	"api/src/authentication"
	"api/src/database"
	"api/src/models"
	"api/src/repositories"
	"api/src/responses"
	"api/src/security"
	"encoding/json"
	"io/ioutil"
	"net/http"
)

func Login(w http.ResponseWriter, r *http.Request) {
	requestBody, err := ioutil.ReadAll(r.Body)

	if err != nil {
		responses.Err(w, http.StatusUnprocessableEntity, err)
		return
	}

	var user models.User

	if err = json.Unmarshal(requestBody, &user); err != nil {
		responses.Err(w, http.StatusBadRequest, err)
		return
	}

	db, err := database.ToConnect()
	if err != nil {
		responses.Err(w, http.StatusInternalServerError, err)
		return
	}

	defer db.Close()

	repository := repositories.NewUserRepository(db)
	userSavedDatabase, err := repository.GetByEmail(user.Email)
	if err != nil {
		responses.Err(w, http.StatusInternalServerError, err)
		return
	}

	if err = security.VerifyPassoword(userSavedDatabase.Password, user.Password); err != nil {
		responses.Err(w, http.StatusUnauthorized, err)
		return
	}

	token, err := authentication.CreateToken(user.ID)

	if err != nil {
		responses.Err(w, http.StatusInternalServerError, err)
		return
	}
	w.Write([]byte(token))
}
