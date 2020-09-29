package internalhttphandlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"

	storage "github.com/ezhk/golang-learning/hw12_13_14_15_calendar/internal/storage"
	"github.com/gorilla/mux"
)

type UserBody struct {
	Status string       `json:"status"`
	User   storage.User `json:"user"`
}

func generateUser(w io.Writer, u storage.User) error {
	return json.NewEncoder(w).Encode(UserBody{
		Status: StatusOK,
		User:   u,
	})
}

func (sh *ServeHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var user storage.User
	w.Header().Set("Content-Type", "application/json")

	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		_ = generateError(w, err)

		return
	}

	user, err = sh.db.CreateUser(user.Email, user.FirstName, user.LastName)
	if err != nil {
		_ = generateError(w, err)

		return
	}

	_ = generateUser(w, user)
}

func (sh *ServeHandler) ReadUser(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	w.Header().Set("Content-Type", "application/json")

	fmt.Println(params)
	user, err := sh.db.GetUserByEmail(params["email"])
	if err != nil {
		_ = generateError(w, err)

		return
	}

	_ = generateUser(w, user)
}

func (sh *ServeHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	var user storage.User

	params := mux.Vars(r)
	userID, err := strconv.ParseInt(params["id"], 10, 64)
	if err != nil {
		_ = generateError(w, err)

		return
	}

	err = json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		_ = generateError(w, err)

		return
	}

	user.ID = userID
	err = sh.db.UpdateUser(user)
	if err != nil {
		_ = generateError(w, err)

		return
	}

	_ = generateUser(w, user)
}

func (sh *ServeHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	userID, err := strconv.ParseInt(params["id"], 10, 64)
	if err != nil {
		_ = generateError(w, err)

		return
	}

	err = sh.db.DeleteUser(storage.User{ID: userID})
	if err != nil {
		_ = generateError(w, err)

		return
	}

	_ = generateUser(w, storage.User{ID: userID})
}
