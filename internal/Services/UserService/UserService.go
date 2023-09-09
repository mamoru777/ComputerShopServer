package UserService

import (
	"ComputerShopServer/internal/DataBaseImplement/Config"
	"ComputerShopServer/internal/Repositories/Models"
	"ComputerShopServer/internal/Repositories/UserRepository"
	"crypto/rand"
	"encoding/json"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"log"
	"math/big"
	"net/http"
)

type UserService struct {
	userrep UserRepository.UserRepository
	config  Config.Config
}

func New(userrep UserRepository.UserRepository, config Config.Config) *UserService {
	return &UserService{
		config:  config,
		userrep: userrep,
	}
}

func (us *UserService) GetHandler() http.Handler {
	router := mux.NewRouter()

	router.HandleFunc("/user/registration", us.CreateUser).Methods(http.MethodPost)
	router.HandleFunc("/user/logincheck", us.GetUserByLogin).Methods(http.MethodGet)
	router.HandleFunc("/user/emailcheck", us.GetUserByEmail).Methods(http.MethodGet)
	router.HandleFunc("/user/autho", us.GetUserByLoginAndPassword).Methods(http.MethodGet)
	router.HandleFunc("/user/confirmemail", us.ConfirmEmail).Methods(http.MethodGet)
	router.HandleFunc("/user/sendconfirmcode", us.SendConfirmCode).Methods(http.MethodGet)

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
	log.Println("Использована функция создания пользователя")
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

type UserCheck struct {
	IsExist bool `json:"isExist"`
}

func (us *UserService) GetUserByEmail(w http.ResponseWriter, r *http.Request) {
	email := r.URL.Query().Get("email")
	e, err := us.userrep.GetByEmail(r.Context(), email)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(&UserCheck{
		IsExist: e,
	})
}

func (us *UserService) GetUserByLogin(w http.ResponseWriter, r *http.Request) {
	log.Println("Использована функция получения пользователя по логину")
	//vars := mux.Vars(r)
	login := r.URL.Query().Get("login") //vars["login"]
	e, err := us.userrep.GetByLogin(r.Context(), login)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(&UserCheck{
		IsExist: e,
	})
}

type UserLogin struct {
	IsExist bool      `json:"isExist"`
	Id      uuid.UUID `json:"id"`
}

func (us *UserService) GetUserByLoginAndPassword(w http.ResponseWriter, r *http.Request) {
	log.Println("Использована функция получения пользователя по логину и паролю")
	login := r.URL.Query().Get("login")
	password := r.URL.Query().Get("password")
	isExist, id, err := us.userrep.GetByLoginAndPassword(r.Context(), login, password)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(&UserLogin{
		Id:      id,
		IsExist: isExist,
	})
}

/*type UserEmail struct {
	email string `json:"email"`
}*/

func (us *UserService) SendConfirmCode(w http.ResponseWriter, r *http.Request) {
	//smtpServer := us.config.SmtpAdr
	//smtpPort := us.config.SmtpPort
	//senderEmail := us.config.SmtpSenderEmail
	//senderPassword := us.config.SmtpSenderPassword
	email := r.URL.Query().Get("email")
	/*req := &UserEmail{}
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}*/
	log.Println(email)
	recipientEmail := email
	log.Println(recipientEmail)
	//subject := "Подтверждение почты на ComputerShop"
	newCode := generateRandomString(5)
	//body := "Код подтверждения\n" + newCode
	//message := []byte("Тема:" + subject + "\n" + body)

	//auth := smtp.PlainAuth("", senderEmail, senderPassword, smtpServer)
	isExist, ecFromBd, err := us.userrep.GetCode(r.Context(), recipientEmail)
	if err != nil {
		log.Println("Не удалось получить запись почты и кода", err)
	}
	if isExist {
		ec := Models.EmailCode{
			ID:    ecFromBd.ID,
			Email: recipientEmail,
			Code:  newCode,
		}
		us.userrep.UpdateCode(r.Context(), &ec)
	} else {
		ec := Models.EmailCode{
			Email: recipientEmail,
			Code:  newCode,
		}
		us.userrep.CreateCode(r.Context(), &ec)
	}
	/*er := smtp.SendMail(smtpServer+":"+smtpPort, auth, senderEmail, []string{recipientEmail}, message)
	if er != nil {
		log.Println("Не удалось отправить код подтверждения", er)
		w.WriteHeader(http.StatusInternalServerError)
	} else {
		w.WriteHeader(http.StatusOK)
	}*/
}

/*type isMatch struct {
	isMatch bool `json:"isExist"`
}*/

func (us *UserService) ConfirmEmail(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")
	email := r.URL.Query().Get("email")
	_, ecFromBd, err := us.userrep.GetCode(r.Context(), email)
	if err != nil {
		log.Println("Не удалось получить запись почты и кода", err)
	}
	var same bool
	log.Println("Код из бд", ecFromBd.Code)
	log.Println("Код от пользователя", code)
	if code == ecFromBd.Code {
		log.Println("Коды совпадают")
		same = true
	} else {
		log.Println("Коды не совпадают")
		same = false
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(&UserCheck{
		IsExist: same,
	})
	/*if code == ecFromBd.Code {
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(&isMatch{
			isMatch: true,
		})
	} else {
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(&isMatch{
			isMatch: false,
		})
	}*/
}

func generateRandomString(length int) string {
	// доступные символы
	chars := "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

	// создаем генератор случайных чисел
	var result string
	for i := 0; i < length; i++ {
		n, _ := rand.Int(rand.Reader, big.NewInt(int64(len(chars))))
		result += string(chars[n.Int64()])
	}

	return result
}
