package user

import (
	"context"
	"errors"
	"log"

	apiV1 "github.com/turao/topics/api/users/v1"
	eventsV1 "github.com/turao/topics/events/users/v1"
	"github.com/turao/topics/metadata"
	"github.com/turao/topics/users/entity/user"
)

type UserRepository interface {
	Save(ctx context.Context, user user.User) error
	FindByID(ctx context.Context, userID user.ID) (user.User, error)
}

type service struct {
	userRepository UserRepository
}

var _ apiV1.Users = (*service)(nil)

func NewService(
	userRepository UserRepository,
) (*service, error) {
	return &service{
		userRepository: userRepository,
	}, nil
}

// RegisterUser implements apiV1.Users
func (svc *service) RegisterUser(ctx context.Context, req apiV1.RegisteUserRequest) (apiV1.RegisterUserResponse, error) {
	log.Println("registering user", req)
	usercfg, errs := user.NewConfig(
		user.WithEmail(req.Email),
		user.WithFirstName(req.FirstName),
		user.WithLastName(req.LastName),
		user.WithTenancy(metadata.Tenancy(req.Tenancy)),
	)
	if len(errs) > 0 {
		return apiV1.RegisterUserResponse{}, errors.Join(errs...)
	}

	user := user.NewUser(usercfg)
	err := svc.userRepository.Save(ctx, user)
	if err != nil {
		return apiV1.RegisterUserResponse{}, err
	}

	log.Println(eventsV1.UserRegistered{
		ID:        user.ID().String(),
		Email:     user.Email(),
		FirstName: user.FirstName(),
		LastName:  user.LastName(),
		Tenancy:   user.Tenancy().String(),
	})

	log.Println("user registered succesfully")
	return apiV1.RegisterUserResponse{
		ID: user.ID().String(),
	}, nil
}

func (svc *service) DeleteUser(ctx context.Context, req apiV1.DeleteUserRequest) (apiV1.DeleteUserResponse, error) {
	log.Println("deleting user", req)
	user, err := svc.userRepository.FindByID(ctx, user.ID(req.ID))
	if err != nil {
		return apiV1.DeleteUserResponse{}, err
	}

	user.Delete()
	err = svc.userRepository.Save(ctx, user)
	if err != nil {
		return apiV1.DeleteUserResponse{}, err
	}

	log.Println("user deleted succesfully")
	return apiV1.DeleteUserResponse{}, nil
}

func (svc *service) GetUserInfo(ctx context.Context, req apiV1.GetUserInfoRequest) (apiV1.GetUserInfoResponse, error) {
	user, err := svc.userRepository.FindByID(ctx, user.ID(req.ID))
	if err != nil {
		return apiV1.GetUserInfoResponse{}, err
	}

	return apiV1.GetUserInfoResponse{
		ID:        user.ID().String(),
		Email:     user.Email(),
		FirstName: user.FirstName(),
		LastName:  user.LastName(),
		Tenancy:   user.Tenancy().String(),
	}, nil
}
