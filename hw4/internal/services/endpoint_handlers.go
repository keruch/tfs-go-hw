package services

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/keruch/tfs-go-hw/hw4/internal/domain"
	"github.com/keruch/tfs-go-hw/hw4/internal/domain/storage"
	"github.com/keruch/tfs-go-hw/hw4/pkg/token"
)

func (cs *ChatService) register(w http.ResponseWriter, r *http.Request) {
	var user domain.UserData
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	err = cs.Controller.SaveUserData(user)
	if err != nil {
		if err == storage.ErrorUserAlreadyExist {
			w.WriteHeader(http.StatusConflict)
			return
		}
		cs.Logger.Errorf("register: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (cs *ChatService) login(w http.ResponseWriter, r *http.Request) {
	var user domain.UserData
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	savedUser, err := cs.Repo.GetUser(user.Username)
	if err != nil {
		if err == storage.ErrorUserNotFound {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		cs.Logger.Errorf("login: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if savedUser.Password != user.Password {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	userToken, err := token.CreateUserToken(savedUser)
	if err != nil {
		cs.Logger.Errorf("login: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	n, err := w.Write([]byte(userToken))
	if err != nil || n != len(userToken) {
		cs.Logger.Errorf("login: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (cs *ChatService) getUsers(w http.ResponseWriter, r *http.Request) {
	users := cs.Repo.GetAllUsers()
	userData, err := json.MarshalIndent(users, "", " ")
	if err != nil || len(userData) == 0 {
		cs.Logger.Errorf("GetUsers: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	n, err := w.Write(userData)
	if err != nil || n != len(userData) {
		cs.Logger.Errorf("GetUsers: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (cs *ChatService) postMessage(w http.ResponseWriter, r *http.Request) {
	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	username := r.Context().Value(ctxTokenKey).(string)
	message := domain.NewMessage(username, string(data))
	cs.Controller.SaveMessage(message)

	w.WriteHeader(http.StatusCreated)
}

func (cs *ChatService) getMessages(w http.ResponseWriter, r *http.Request) {
	messages := cs.Repo.GetMessages()
	userData, err := json.MarshalIndent(messages, "", " ")
	if err != nil {
		cs.Logger.Errorf("GetMessages: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	n, err := w.Write(userData)
	if err != nil || n != len(userData) {
		cs.Logger.Errorf("GetMessages: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (cs *ChatService) getNumMessages(w http.ResponseWriter, r *http.Request) {
	numStr := chi.URLParam(r, "num")
	num, err := strconv.Atoi(numStr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	messages, err := cs.Repo.GetNumMessages(num)
	if err != nil {
		if err == storage.ErrorIncorrectNumber {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		cs.Logger.Errorf("GetPrivateMessages: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	userData, err := json.MarshalIndent(messages, "", " ")
	if err != nil {
		cs.Logger.Errorf("GetMessages: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	n, err := w.Write(userData)
	if err != nil || n != len(userData) {
		cs.Logger.Errorf("GetMessages: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (cs *ChatService) postPrivateMessage(w http.ResponseWriter, r *http.Request) {
	toUser := chi.URLParam(r, "to_user")
	message, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	username := r.Context().Value(ctxTokenKey).(string)
	privateMessage := domain.NewMessage(username, string(message))
	cs.Controller.SendPrivateMessage(toUser, privateMessage)

	w.WriteHeader(http.StatusCreated)
}

func (cs *ChatService) getPrivateMessages(w http.ResponseWriter, r *http.Request) {
	username := r.Context().Value(ctxTokenKey).(string)
	messages, err := cs.Repo.GetPrivateMessages(username)
	if err != nil {
		if err == storage.ErrorMessageboxEmpty {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		cs.Logger.Errorf("GetPrivateMessages: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	messagesData, err := json.MarshalIndent(messages, "", " ")
	if err != nil {
		cs.Logger.Errorf("GetPrivateMessages: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	n, err := w.Write(messagesData)
	if err != nil || n != len(messagesData) {
		cs.Logger.Errorf("GetPrivateMessages: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
