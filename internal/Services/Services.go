package Services

import (
	"ComputerShopServer/internal/DataBaseImplement/Config"
	"ComputerShopServer/internal/Repositories/GoodRepository"
	"ComputerShopServer/internal/Repositories/Models"
	"ComputerShopServer/internal/Repositories/UserRepository"
	"context"
	"crypto/rand"
	"encoding/json"
	"errors"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"io"
	"log"
	"math/big"
	"net/http"
	"net/smtp"
	"strconv"
)

type Service struct {
	userrep UserRepository.UserRepository
	goodrep GoodRepository.GoodRepository
	config  Config.Config
}

func New(userrep UserRepository.UserRepository, goodrep GoodRepository.GoodRepository, config Config.Config) *Service {
	return &Service{
		config:  config,
		userrep: userrep,
		goodrep: goodrep,
	}
}

func (us *Service) GetHandler() http.Handler {
	router := mux.NewRouter()

	router.HandleFunc("/user/registration", us.CreateUser).Methods(http.MethodPost)
	router.HandleFunc("/user/logincheck", us.GetUserByLogin).Methods(http.MethodGet)
	router.HandleFunc("/user/emailcheck", us.GetUserByEmail).Methods(http.MethodGet)
	router.HandleFunc("/user/autho", us.GetUserByLoginAndPassword).Methods(http.MethodGet)
	router.HandleFunc("/user/confirmemail", us.ConfirmEmail).Methods(http.MethodGet)
	router.HandleFunc("/user/sendconfirmcode", us.SendConfirmCode).Methods(http.MethodGet)
	router.HandleFunc("/user/userinfo", us.GetUserInfo).Methods(http.MethodGet)
	router.HandleFunc("/user/getavatar", us.GetAvatar).Methods(http.MethodGet)
	router.HandleFunc("/user/patchavatar", us.PatchAvatar).Methods(http.MethodPatch)
	router.HandleFunc("/user/changepassword", us.ChangePassword).Methods(http.MethodPatch)
	router.HandleFunc("/user/changedata", us.ChangeData).Methods(http.MethodPatch)

	router.HandleFunc("/good/create", us.CreateGood).Methods(http.MethodPost)
	router.HandleFunc("/good/goodcheck", us.GetGoodByName).Methods(http.MethodGet)
	router.HandleFunc("/good/goodsbytype", us.GetGoodsByType).Methods(http.MethodGet)
	router.HandleFunc("/good/getgood", us.GetGood).Methods(http.MethodGet)
	router.HandleFunc("/good/goodsbyid", us.GetGoodsById).Methods(http.MethodGet)
	router.HandleFunc("/good/changegood", us.ChangeGood).Methods(http.MethodPatch)
	router.HandleFunc("/good/deletegood", us.DeleteGood).Methods(http.MethodPatch)
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
	Role     string `json:"role"`
}

func (us *Service) CreateUser(w http.ResponseWriter, r *http.Request) {
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
		Role:     req.Role,
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

func (us *Service) GetUserByEmail(w http.ResponseWriter, r *http.Request) {
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

func (us *Service) GetUserByLogin(w http.ResponseWriter, r *http.Request) {
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
	Role    string    `json:"role"`
}

func (us *Service) GetUserByLoginAndPassword(w http.ResponseWriter, r *http.Request) {
	log.Println("Использована функция получения пользователя по логину и паролю")
	login := r.URL.Query().Get("login")
	password := r.URL.Query().Get("password")
	isExist, id, role, err := us.userrep.GetByLoginAndPassword(r.Context(), login, password)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(&UserLogin{
		Id:      id,
		IsExist: isExist,
		Role:    role,
	})
}

func (us *Service) SendConfirmCode(w http.ResponseWriter, r *http.Request) {
	smtpServer := us.config.SmtpAdr
	smtpPort := us.config.SmtpPort
	senderEmail := us.config.SmtpSenderEmail
	senderPassword := us.config.SmtpSenderPassword
	email := r.URL.Query().Get("email")
	log.Println(email)
	recipientEmail := email
	log.Println(recipientEmail)
	subject := "Подтверждение почты на ComputerShop"
	newCode := generateRandomString(5)
	body := "Код подтверждения\n" + newCode
	message := []byte("Тема:" + subject + "\n" + body)

	auth := smtp.PlainAuth("", senderEmail, senderPassword, smtpServer)
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
	er := smtp.SendMail(smtpServer+":"+smtpPort, auth, senderEmail, []string{recipientEmail}, message)
	if er != nil {
		log.Println("Не удалось отправить код подтверждения", er)
		w.WriteHeader(http.StatusInternalServerError)
	} else {
		w.WriteHeader(http.StatusOK)
	}
}

func (us *Service) ConfirmEmail(w http.ResponseWriter, r *http.Request) {
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
}

type UserInfo struct {
	Login    string `json:"login"`
	Password string `json:"password"`
	Name     string `json:"name"`
	Lastname string `json:"lastname"`
	Surname  string `json:"surname"`
	Email    string `json:"email"`
}

func (us *Service) GetUserInfo(w http.ResponseWriter, r *http.Request) {
	log.Println("Использована функция получения пользователя")
	id := r.URL.Query().Get("id")
	uuid, err := uuid.Parse(id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("Не удалось конвертировать строку в uuid", err)
		return
	}
	user, err := us.userrep.Get(r.Context(), uuid)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("Не удалось получить запись пользователя из бд", err)
		return
	}
	//log.Println(user)
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(&UserInfo{
		Login:    user.Login,
		Password: user.Password,
		Name:     user.Name,
		Lastname: user.LastName,
		Surname:  user.SurName,
		Email:    user.Email,
	})
	log.Println(user.Login)
}

type Avatar struct {
	Avatar string `json:"avatar"`
	Id     string `json:"id"`
}

func (us *Service) PatchAvatar(w http.ResponseWriter, r *http.Request) {
	log.Println("Использована функция добавления аватара")
	id := r.FormValue("id")
	uuid, err := uuid.Parse(id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("Не удалось конвертировать строку в uuid", err)
		return
	}
	file, _, err := r.FormFile("avatar")
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	defer file.Close()
	fileData, err := io.ReadAll(file)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	user, err := us.userrep.Get(r.Context(), uuid)
	user.Avatar = fileData
	//log.Println("Массив байт - ", user.Avatar)
	err = us.userrep.Update(r.Context(), user)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	} else {
		w.WriteHeader(http.StatusOK)
	}
}

func (us *Service) GetAvatar(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	uuid, err := uuid.Parse(id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("Не удалось конвертировать строку в uuid", err)
		return
	}
	user, err := us.userrep.Get(r.Context(), uuid)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("Не удалось получить запись пользователя из бд", err)
		return
	}
	avatarBytes := user.Avatar
	w.Header().Set("Content-Type", "image/jpeg")
	w.Header().Set("Content-Length", strconv.Itoa(len(avatarBytes)))
	if _, err := w.Write(avatarBytes); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

type ChangePass struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (us *Service) ChangePassword(w http.ResponseWriter, r *http.Request) {
	//email := r.URL.Query().Get("email")
	//pass := r.URL.Query().Get("pass")
	req := &ChangePass{}
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	user, err := us.userrep.GetByEmailUser(r.Context(), req.Email)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("Не удалось получить запись пользователя из бд", err)
		return
	}
	user.Password = req.Password
	err = us.userrep.Update(r.Context(), user)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	} else {
		w.WriteHeader(http.StatusOK)
	}
}

type ChangeData struct {
	Id       string `json:"id"`
	Login    string `json:"login"`
	Name     string `json:"name"`
	LastName string `json:"lastname"`
	SurName  string `json:"surname"`
}

func (us *Service) ChangeData(w http.ResponseWriter, r *http.Request) {
	log.Println("Использована функция изменения даты пользователя")
	req := &ChangeData{}
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	log.Println(req)
	uuid, err := uuid.Parse(req.Id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("Не удалось конвертировать строку в uuid", err)
		return
	}
	user, err := us.userrep.Get(r.Context(), uuid)
	if req.Login != "" {
		user.Login = req.Login
	}
	if req.Name != "" {
		user.Name = req.Name
	}
	if req.LastName != "" {
		user.LastName = req.LastName
	}
	if req.SurName != "" {
		user.SurName = req.SurName
	}
	err = us.userrep.Update(r.Context(), user)
	if err != nil {
		log.Println("Не удалось обновить данные о пользователе", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (us *Service) CreateAdmin() error {
	admin := Models.Usr{
		Login:    "admin",
		Password: "admin",
		Email:    "mamoru90000@gmail.com",
		Role:     "admin",
	}
	isExist, err := us.userrep.GetByLogin(context.Background(), admin.Login)
	if err != nil {
		log.Println("Не удалось получить пользователя по логину")
	}
	if isExist {
		log.Println("Пользователь админ уже существует")
		err = errors.New("Пользователь админ уже существует")
		return err
	}
	us.userrep.Create(context.Background(), &admin)
	return nil
}

type Good struct {
	Type        string  `json:"goodtype"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	Avatar      []byte  `json:"avatar"`
}

func (us *Service) CreateGood(w http.ResponseWriter, r *http.Request) {
	log.Println("Использована функция создания товара")
	file, _, err := r.FormFile("avatar")
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	defer file.Close()
	fileData, err := io.ReadAll(file)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	gtype := r.FormValue("goodtype")
	name := r.FormValue("name")
	descr := r.FormValue("description")
	price, err := strconv.ParseFloat(r.FormValue("price"), 64)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	good := &Models.Good{
		Name:        name,
		Description: descr,
		GoodType:    gtype,
		Price:       price,
		Avatar:      fileData,
	}
	//log.Println(good)
	err = us.goodrep.Create(r.Context(), good)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	} else {
		w.WriteHeader(http.StatusOK)
	}
}

type GoodCheck struct {
	IsExist bool `json:"isExist"`
}

func (us *Service) GetGoodByName(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")
	e, err := us.goodrep.GetByName(r.Context(), name)
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

func (us *Service) GetGoodsById(w http.ResponseWriter, r *http.Request) {
	var goods []Models.Good
	ides := r.URL.Query()["ides"]
	for _, id := range ides {
		uuid, err := uuid.Parse(id)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		Good, err := us.goodrep.Get(r.Context(), uuid)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		goods = append(goods, *Good)
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(goods)
}

func (us *Service) GetGoodsByType(w http.ResponseWriter, r *http.Request) {
	gtype := r.URL.Query().Get("good_type")
	goods, err := us.goodrep.GetByType(r.Context(), gtype)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	} else {
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(goods)
	}

}

func (us *Service) GetGood(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("good_id")
	uuid, err := uuid.Parse(id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("Не удалось конвертировать строку в uuid", err)
		return
	}
	good, err := us.goodrep.Get(r.Context(), uuid)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	} else {
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(good)
	}
}

func (us *Service) ChangeGood(w http.ResponseWriter, r *http.Request) {
	log.Println("Использована функция изменения товара")
	file, _, err := r.FormFile("avatar")
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	defer file.Close()
	fileData, err := io.ReadAll(file)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	gtype := r.FormValue("goodtype")
	name := r.FormValue("name")
	descr := r.FormValue("description")
	price, err := strconv.ParseFloat(r.FormValue("price"), 64)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	id := r.FormValue("id")
	uuid, err := uuid.Parse(id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("Не удалось конвертировать строку в uuid", err)
		return
	}
	good, err := us.goodrep.Get(r.Context(), uuid)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	good.GoodType = gtype
	good.Name = name
	good.Description = descr
	good.Price = price
	good.Avatar = fileData
	err = us.goodrep.Update(r.Context(), good)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

type GoodId struct {
	Id string `json:"id"`
}

func (us *Service) DeleteGood(w http.ResponseWriter, r *http.Request) {
	req := &GoodId{}
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	id := req.Id
	uuid, err := uuid.Parse(id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("Не удалось конвертировать строку в uuid", err)
		return
	}
	err = us.goodrep.Delete(r.Context(), uuid)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("Не удалось удалить товар из бд")
		return
	}
	w.WriteHeader(http.StatusOK)
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
