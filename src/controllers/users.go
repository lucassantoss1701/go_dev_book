package controllers

import (
	"api/src/authentication"
	"api/src/database"
	"api/src/models"
	"api/src/repositories"
	"api/src/responses"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
)

func CreateUser(w http.ResponseWriter, r *http.Request) {
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

	if err = user.Prepare("registration"); err != nil {
		responses.Err(w, http.StatusBadRequest, err)
		return
	}

	db, err := database.ToConnect()
	if err != nil {
		responses.Err(w, http.StatusInternalServerError, err)
	}

	defer db.Close()

	repository := repositories.NewUserRepository(db)
	user.ID, err = repository.Create(user)
	if err != nil {
		responses.Err(w, http.StatusInternalServerError, err)
	}

	responses.JSON(w, http.StatusCreated, user)
}

func GetUsers(w http.ResponseWriter, r *http.Request) {
	nameOrNick := strings.ToLower(r.URL.Query().Get("user"))

	db, err := database.ToConnect()
	if err != nil {
		responses.Err(w, http.StatusInternalServerError, err)
	}
	defer db.Close()

	repository := repositories.NewUserRepository(db)

	users, err := repository.GetByNameOrNick(nameOrNick)
	if err != nil {
		responses.Err(w, http.StatusInternalServerError, err)
	}

	responses.JSON(w, http.StatusOK, users)
}

func GetUserById(w http.ResponseWriter, r *http.Request) {
	parameters := mux.Vars(r)

	userID, err := strconv.ParseUint(parameters["userId"], 10, 64)
	if err != nil {
		responses.Err(w, http.StatusBadRequest, err)
		return
	}

	db, err := database.ToConnect()
	if err != nil {
		responses.Err(w, http.StatusInternalServerError, err)
	}
	defer db.Close()

	repository := repositories.NewUserRepository(db)
	user, err := repository.GetUserById(userID)
	if err != nil {
		responses.Err(w, http.StatusInternalServerError, err)
	}

	responses.JSON(w, http.StatusOK, user)

}

func UpdateUser(w http.ResponseWriter, r *http.Request) {
	parameters := mux.Vars(r)
	userID, err := strconv.ParseUint(parameters["userId"], 10, 64)
	if err != nil {
		responses.Err(w, http.StatusBadRequest, err)
		return
	}

	userIDToken, err := authentication.ExtractUserID(r)
	if err != nil {
		responses.Err(w, http.StatusUnauthorized, err)
		return
	}

	if userID != userIDToken {
		responses.Err(w, http.StatusForbidden, errors.New("Não é possível atualizar um usuário que não é o seu!"))
		return
	}

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

	if err = user.Prepare("update"); err != nil {
		responses.Err(w, http.StatusInternalServerError, err)
		return
	}

	db, err := database.ToConnect()
	if err != nil {
		responses.Err(w, http.StatusInternalServerError, err)
	}

	defer db.Close()

	repository := repositories.NewUserRepository(db)
	if err = repository.Update(userID, user); err != nil {
		responses.Err(w, http.StatusInternalServerError, err)
		return
	}

	responses.JSON(w, http.StatusNoContent, nil)
}
func DeleteUser(w http.ResponseWriter, r *http.Request) {
	parameters := mux.Vars(r)
	userID, err := strconv.ParseUint(parameters["userId"], 10, 64)
	if err != nil {
		responses.Err(w, http.StatusBadRequest, err)
		return
	}

	userIDToken, err := authentication.ExtractUserID(r)
	if err != nil {
		responses.Err(w, http.StatusUnauthorized, err)
		return
	}

	if userID != userIDToken {
		responses.Err(w, http.StatusForbidden, errors.New("Não é possível deletar um usuário que não é o seu!"))
		return
	}

	db, err := database.ToConnect()
	if err != nil {
		responses.Err(w, http.StatusInternalServerError, err)
	}

	defer db.Close()

	repository := repositories.NewUserRepository(db)
	if err = repository.Delete(userID); err != nil {
		responses.Err(w, http.StatusInternalServerError, err)
		return
	}

	responses.JSON(w, http.StatusNoContent, nil)
}

func FollowUser(w http.ResponseWriter, r *http.Request) {
	followerID, err := authentication.ExtractUserID(r)
	if err != nil {
		responses.Err(w, http.StatusUnauthorized, err)
		return
	}

	parameters := mux.Vars(r)
	userID, err := strconv.ParseUint(parameters["userId"], 10, 64)
	if err != nil {
		responses.Err(w, http.StatusBadRequest, err)
		return
	}

	if userID == followerID {
		responses.Err(w, http.StatusForbidden, errors.New("Não é possível seguir você mesmo."))
		return
	}

	db, err := database.ToConnect()
	if err != nil {
		responses.Err(w, http.StatusInternalServerError, err)
	}

	defer db.Close()

	repository := repositories.NewUserRepository(db)
	if err = repository.FollowUser(userID, followerID); err != nil {
		responses.Err(w, http.StatusInternalServerError, err)
		return
	}

	responses.JSON(w, http.StatusNoContent, nil)
}
