package UserService

import (
	"ComputerShopServer/internal/Repositories/UserRepository"
	userapi "ComputerShopServer/pkg"
	"context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type UserService struct {
	userapi.UnimplementedUserServiceServer
	userrep UserRepository.UserRepository
}

func New(userrep UserRepository.UserRepository) *UserService {
	return &UserService{

		userrep: userrep,
	}
}

//type Service struct {
//	userapi.UnimplementedUserServiceServer
//	userservice *UserService
//}
//
//func NewServ() *Service {
//	usrep := UserRepository.UserRepository()
//	return &Service{
//		userservice: New(),
//	}
//}

func (us *UserService) CreateUser(ctx context.Context, request *userapi.CreateUserRequest) (*userapi.CreateUserResponse, error) {
	model := UserRepository.Usr{
		Login:    request.Login,
		Password: request.Password,
		Name:     request.Name,
		LastName: request.Lastname,
		SurName:  request.Surname,
		Email:    request.Email,
		Avatar:   request.Avatar,
	}
	if err := model.Validate(); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	us.userrep.Create(ctx, &model)
	return &userapi.CreateUserResponse{}, nil
}
