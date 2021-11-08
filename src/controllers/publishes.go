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

	"github.com/gorilla/mux"
)

func CreatePublish(w http.ResponseWriter, r *http.Request) {
	usuarioID, erro := authentication.ExtractUserID(r)
	if erro != nil {
		responses.Err(w, http.StatusUnauthorized, erro)
		return
	}

	requestBody, erro := ioutil.ReadAll(r.Body)
	if erro != nil {
		responses.Err(w, http.StatusUnprocessableEntity, erro)
		return
	}

	var publish models.Publish
	if erro = json.Unmarshal(requestBody, &publish); erro != nil {
		responses.Err(w, http.StatusBadRequest, erro)
		return
	}

	publish.AuthorID = usuarioID

	if erro = publish.Prepare(); erro != nil {
		responses.Err(w, http.StatusBadRequest, erro)
		return
	}

	db, erro := database.ToConnect()
	if erro != nil {
		responses.Err(w, http.StatusInternalServerError, erro)
		return
	}
	defer db.Close()

	repository := repositories.NewPublishRepository(db)
	publish.ID, erro = repository.Create(publish)
	if erro != nil {
		responses.Err(w, http.StatusInternalServerError, erro)
		return
	}

	responses.JSON(w, http.StatusCreated, publish)
}

func GetPublishes(w http.ResponseWriter, r *http.Request) {
	usuarioID, erro := authentication.ExtractUserID(r)
	if erro != nil {
		responses.Err(w, http.StatusUnauthorized, erro)
		return
	}

	db, erro := database.ToConnect()
	if erro != nil {
		responses.Err(w, http.StatusInternalServerError, erro)
		return
	}
	defer db.Close()

	repository := repositories.NewPublishRepository(db)
	publicacoes, erro := repository.Get(usuarioID)
	if erro != nil {
		responses.Err(w, http.StatusInternalServerError, erro)
		return
	}

	responses.JSON(w, http.StatusOK, publicacoes)
}

func GetPublish(w http.ResponseWriter, r *http.Request) {
	parameters := mux.Vars(r)
	publishID, erro := strconv.ParseUint(parameters["publishId"], 10, 64)
	if erro != nil {
		responses.Err(w, http.StatusBadRequest, erro)
		return
	}

	db, erro := database.ToConnect()
	if erro != nil {
		responses.Err(w, http.StatusInternalServerError, erro)
		return
	}
	defer db.Close()

	repository := repositories.NewPublishRepository(db)
	publish, erro := repository.GetByID(publishID)
	if erro != nil {
		responses.Err(w, http.StatusInternalServerError, erro)
		return
	}

	responses.JSON(w, http.StatusOK, publish)
}

func UpdatePublish(w http.ResponseWriter, r *http.Request) {
	usuarioID, erro := authentication.ExtractUserID(r)
	if erro != nil {
		responses.Err(w, http.StatusUnauthorized, erro)
		return
	}

	parameters := mux.Vars(r)
	publishID, erro := strconv.ParseUint(parameters["publishId"], 10, 64)
	if erro != nil {
		responses.Err(w, http.StatusBadRequest, erro)
		return
	}

	db, erro := database.ToConnect()
	if erro != nil {
		responses.Err(w, http.StatusInternalServerError, erro)
		return
	}
	defer db.Close()

	repository := repositories.NewPublishRepository(db)
	publishSalvaNoBanco, erro := repository.GetByID(publishID)
	if erro != nil {
		responses.Err(w, http.StatusInternalServerError, erro)
		return
	}

	if publishSalvaNoBanco.AuthorID != usuarioID {
		responses.Err(w, http.StatusForbidden, errors.New("Não é possível atualizar uma publicação que não seja sua"))
		return
	}

	requestBody, erro := ioutil.ReadAll(r.Body)
	if erro != nil {
		responses.Err(w, http.StatusUnprocessableEntity, erro)
		return
	}

	var publish models.Publish
	if erro = json.Unmarshal(requestBody, &publish); erro != nil {
		responses.Err(w, http.StatusBadRequest, erro)
		return
	}

	if erro = publish.Prepare(); erro != nil {
		responses.Err(w, http.StatusBadRequest, erro)
		return
	}

	if erro = repository.Update(publishID, publish); erro != nil {
		responses.Err(w, http.StatusInternalServerError, erro)
		return
	}

	responses.JSON(w, http.StatusNoContent, nil)
}

func DeletePublish(w http.ResponseWriter, r *http.Request) {
	usuarioID, erro := authentication.ExtractUserID(r)
	if erro != nil {
		responses.Err(w, http.StatusUnauthorized, erro)
		return
	}

	parameters := mux.Vars(r)
	publishID, erro := strconv.ParseUint(parameters["publishId"], 10, 64)
	if erro != nil {
		responses.Err(w, http.StatusBadRequest, erro)
		return
	}

	db, erro := database.ToConnect()
	if erro != nil {
		responses.Err(w, http.StatusInternalServerError, erro)
		return
	}
	defer db.Close()

	repository := repositories.NewPublishRepository(db)
	publishSalvaNoBanco, erro := repository.GetByID(publishID)
	if erro != nil {
		responses.Err(w, http.StatusInternalServerError, erro)
		return
	}

	if publishSalvaNoBanco.AuthorID != usuarioID {
		responses.Err(w, http.StatusForbidden, errors.New("Não é possível deletar uma publicação que não seja sua"))
		return
	}

	if erro = repository.Delete(publishID); erro != nil {
		responses.Err(w, http.StatusInternalServerError, erro)
		return
	}

	responses.JSON(w, http.StatusNoContent, nil)
}

func GetPublishByUser(w http.ResponseWriter, r *http.Request) {
	parameters := mux.Vars(r)
	usuarioID, erro := strconv.ParseUint(parameters["usuarioId"], 10, 64)
	if erro != nil {
		responses.Err(w, http.StatusBadRequest, erro)
		return
	}

	db, erro := database.ToConnect()
	if erro != nil {
		responses.Err(w, http.StatusInternalServerError, erro)
		return
	}
	defer db.Close()

	repository := repositories.NewPublishRepository(db)
	publicacoes, erro := repository.GetByUser(usuarioID)
	if erro != nil {
		responses.Err(w, http.StatusInternalServerError, erro)
		return
	}

	responses.JSON(w, http.StatusOK, publicacoes)
}

func LikePublish(w http.ResponseWriter, r *http.Request) {
	parameters := mux.Vars(r)
	publishID, erro := strconv.ParseUint(parameters["publishId"], 10, 64)
	if erro != nil {
		responses.Err(w, http.StatusBadRequest, erro)
		return
	}

	db, erro := database.ToConnect()
	if erro != nil {
		responses.Err(w, http.StatusInternalServerError, erro)
		return
	}
	defer db.Close()

	repository := repositories.NewPublishRepository(db)
	if erro = repository.Like(publishID); erro != nil {
		responses.Err(w, http.StatusInternalServerError, erro)
		return
	}

	responses.JSON(w, http.StatusNoContent, nil)
}

func DeslikePublish(w http.ResponseWriter, r *http.Request) {
	parameters := mux.Vars(r)
	publishID, erro := strconv.ParseUint(parameters["publishId"], 10, 64)
	if erro != nil {
		responses.Err(w, http.StatusBadRequest, erro)
		return
	}

	db, erro := database.ToConnect()
	if erro != nil {
		responses.Err(w, http.StatusInternalServerError, erro)
		return
	}
	defer db.Close()

	repository := repositories.NewPublishRepository(db)
	if erro = repository.Deslike(publishID); erro != nil {
		responses.Err(w, http.StatusInternalServerError, erro)
		return
	}

	responses.JSON(w, http.StatusNoContent, nil)
}
