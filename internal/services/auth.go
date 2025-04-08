package services

import (
	"errors"
	"fmt"

	"github.com/aakritigkmit/payment-gateway/internal/dto"
	"github.com/aakritigkmit/payment-gateway/internal/model"
	"github.com/aakritigkmit/payment-gateway/internal/repository"
	"github.com/aakritigkmit/payment-gateway/internal/utils"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	AuthRepo *repository.AuthRepo
	Repo     *repository.UserRepo
}

func NewAuthService(authRepo *repository.AuthRepo, repo *repository.UserRepo) *AuthService {
	return &AuthService{AuthRepo: authRepo, Repo: repo}
}

func (s *AuthService) RegisterUser(userReq dto.UserRequest) error {

	existingUser, _ := s.AuthRepo.FindUserByEmail(userReq.EmailID)
	if existingUser != nil {
		return fmt.Errorf("user with email %s already exists", userReq.EmailID)
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(userReq.Password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("error hashing password: %v", err)
	}

	user := model.User{

		FirstName:    userReq.FirstName,
		LastName:     userReq.LastName,
		EmailID:      userReq.EmailID,
		Password:     string(hashedPassword),
		MobileNumber: userReq.MobileNumber,
	}

	userID, err := s.Repo.CreateUser(user) // Call the repository method
	if err != nil {
		fmt.Println("Error while register user in DB:", err)
		return err
	}

	fmt.Println("User successfully Registered ", userID)
	return nil

}

func (s *AuthService) Login(email, password string) (string, error) {
	user, err := s.AuthRepo.FindUserByEmail(email)
	if err != nil {
		fmt.Println("Error finding user:", err) // Print error
		return "", errors.New("invalid credentials")
	}

	// Compare hashed password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		fmt.Println("Password mismatch:", err) // Print error
		return "", errors.New("invalid credentials")
	}

	// Generate JWT token
	tokenString, err := utils.GenerateJWT(user.ID.Hex(), user.EmailID)
	if err != nil {
		fmt.Println("Error generating JWT:", err) // Print error
		return "", errors.New("failed to generate token")
	}

	return tokenString, nil
}
