package UserService

import (
	"ComputerShopServer/internal/Repositories/Models"
	"ComputerShopServer/internal/Repositories/UserRepository"
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
)

type UserService struct {
	userrep UserRepository.UserRepository
}

func New(userrep UserRepository.UserRepository) *UserService {
	return &UserService{

		userrep: userrep,
	}
}

type CreateUserRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
	Email    string `json:"email"`
}

func (us *UserService) CreateUser(w http.ResponseWriter, r *http.Request) {
	req := &CreateUserRequest{}
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	u := Models.Usr{
		Login:    req.Login,
		Password: req.Password,
		Email:    req.Email,
	}
	if err := us.userrep.Create(r.Context(), &u); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)

}

func (us *UserService) GetHandler() http.Handler {
	router := mux.NewRouter()

	router.HandleFunc("/user", us.CreateUser).Methods(http.MethodPost)
	return router
}
