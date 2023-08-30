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

func (us *UserService) GetHandler() http.Handler {
	router := mux.NewRouter()

	router.HandleFunc("/users", us.CreateUser).Methods(http.MethodPost)

	/*header := handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"})
	method := handlers.AllowedMethods([]string{"POST"})
	origins := handlers.AllowedOrigins([]string{"*"})*/

	//return handlers.CORS(header, method, origins)(router)
	return router
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
