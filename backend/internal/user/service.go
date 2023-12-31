package user

import (
	"fmt"
	"github/wry-0313/exchange/internal/models"
	"github/wry-0313/exchange/pkg/security"
	"github/wry-0313/exchange/pkg/validator"
	"regexp"
	"strings"

	"github.com/oklog/ulid/v2"
)

type service struct {
	userRepo  Repository
	validator validator.Validate
}

type Service interface {
	CreateUser(input CreateUserInput) (models.User, error)
	UpdateUserName(userID, name string) error
	GetUser(userID string) (models.User, error)
	GetUserPrivateInfo(userID string) (UserPrivateInfo, error)
}

func NewService(userRepo Repository, validator validator.Validate) Service {
	return &service{
		userRepo:  userRepo,
		validator: validator,
	}
}

func (s *service) CreateUser(input CreateUserInput) (models.User, error) {
	if err := s.validator.Struct(input); err != nil {
		return models.User{}, fmt.Errorf("service: validation error: %w", err)
	}

	// Prepare user input
	// id := ulid.Make()
	name := toNameCase(input.Name)
	user := models.User{
		ID:        ulid.Make().String(),
		Name:      name,
		Email:     input.Email,
	}

	// Hash the password
	if input.Password != nil {
		hashedPassword, err := security.HashPassword(*input.Password)
		if err != nil {
			return models.User{}, fmt.Errorf("service: hashing password: %w", err)
		}
		user.Password = &hashedPassword
	}

	if err := s.userRepo.CreateUser(user); err != nil {
		return models.User{}, fmt.Errorf("service: failed creating user: %w", err)
	}

	// Hide password
	user.Password = nil
	return user, nil
}

func (s *service) UpdateUserName(userID, name string) error {
	if err := s.validator.Var(name, "required"); err != nil {
		return fmt.Errorf("service: validation error: %w", err)
	}

	name = toNameCase(name)

	if err := s.userRepo.UpdateUserName(userID, name); err != nil {
		return fmt.Errorf("service: failed updating user name: %w", err)
	}
	return nil
}

func (s *service) GetUser(userID string) (models.User, error) {
	if err := s.validator.Var(userID, "required"); err != nil {
		return models.User{}, fmt.Errorf("service: validation error: %w", err)
	}

	user, err := s.userRepo.GetUser(userID)
	if err != nil {
		return models.User{}, fmt.Errorf("service: failed getting user: %w", err)
	}

	// Hide password
	user.Password = nil
	return user, nil
}

// toNameCase creates a regular expression to match word boundaries and convert them to name case
func toNameCase(word string) string {
	re := regexp.MustCompile(`\b\w`)
	nameCase := re.ReplaceAllStringFunc(word, strings.ToUpper)

	return nameCase
}

func (s *service) GetUserPrivateInfo(userID string) (UserPrivateInfo, error) {
	if err := s.validator.Var(userID, "required"); err != nil {
		return UserPrivateInfo{}, fmt.Errorf("service: validation error: %w", err)
	}

	userPrivateInfo, err := s.userRepo.GetUserPrivateInfo(userID)
	if err != nil {
		return UserPrivateInfo{}, fmt.Errorf("service: failed getting user private info: %w", err)
	}

	return userPrivateInfo, nil
}