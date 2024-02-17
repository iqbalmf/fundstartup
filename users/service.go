package users

import (
	"errors"
	"golang.org/x/crypto/bcrypt"
	"strconv"
)

type Service interface {
	RegisterUser(input RegisterUserInput) (User, error)
	LoginUser(input LoginInput) (User, error)
	IsEmailAvailable(input CheckEmailInput) (bool, error)
	SaveAvatar(ID int, fileLocation string) (User, error)
	GetUserById(ID int) (User, error)
}

type service struct {
	repository Repository
}

func NewService(repository Repository) *service {
	return &service{repository: repository}
}
func (s *service) RegisterUser(input RegisterUserInput) (User, error) {
	user := User{}
	user.Name = input.Name
	user.Email = input.Email
	user.Occupation = input.Occupation
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.MinCost)
	if err != nil {
		return user, err
	}
	user.PasswordHash = string(passwordHash)
	user.Role = "user"

	newUser, err := s.repository.Save(user)
	if err != nil {
		return newUser, err
	}
	return newUser, nil
}
func (s *service) LoginUser(input LoginInput) (User, error) {
	email := input.Email
	password := input.Password

	user, err := s.repository.FindByEmail(email)
	if err != nil {
		return user, err
	}
	if user.ID == 0 {
		return user, errors.New("No user found")
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		return user, err
	}
	return user, nil
}
func (s *service) IsEmailAvailable(input CheckEmailInput) (bool, error) {
	email := input.Email
	user, err := s.repository.FindByEmail(email)
	if err != nil {
		return false, err
	}
	if user.ID != 0 {
		return false, err
	}
	return true, nil
}
func (s *service) SaveAvatar(ID int, fileLocation string) (User, error) {
	//getting user by ID
	//update attribute avatar file name
	//save update avatar file name
	user, err := s.repository.FindById(strconv.Itoa(ID))
	if err != nil {
		return user, err
	}
	user.AvatarFileName = fileLocation
	updatedUser, err := s.repository.UpdateUser(user)
	if err != nil {
		return user, err
	}
	return updatedUser, nil
}
func (s *service) GetUserById(ID int) (User, error) {
	user, err := s.repository.FindById(strconv.Itoa(ID))
	if err != nil {
		return user, err
	}
	if user.ID == 0 {
		return user, errors.New("No user found on with that ID")
	}
	return user, nil
}
