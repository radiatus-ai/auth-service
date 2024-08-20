package auth

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"github.com/radiatus-ai/auth-service/internal/model"
	"github.com/radiatus-ai/auth-service/internal/repository"
	"google.golang.org/api/idtoken"
)

type Service interface {
	LoginGoogle(token string) (*UserData, error)
	VerifyToken(token string) (string, error)
	GetUserByID(userID string) (*model.User, error)
}

type service struct {
	userRepo        repository.UserRepository
	orgRepo         repository.OrganizationRepository
	jwtSecret       string
	googleClientIDs []string
	emailWhitelist  []string
}

func NewService(userRepo repository.UserRepository, orgRepo repository.OrganizationRepository, jwtSecret string, googleClientIDs []string, emailWhitelist []string) Service {
	return &service{
		userRepo:        userRepo,
		orgRepo:         orgRepo,
		jwtSecret:       jwtSecret,
		googleClientIDs: googleClientIDs,
		emailWhitelist:  emailWhitelist,
	}
}

type UserData struct {
	Token          string     `json:"token"`
	User           model.User `json:"user"`
	OrganizationID uuid.UUID  `json:"organization_id"`
}

func (s *service) LoginGoogle(token string) (*UserData, error) {
	log.Println("Starting Google login process")

	var payload *idtoken.Payload
	var err error

	for _, clientID := range s.googleClientIDs {
		payload, err = idtoken.Validate(context.Background(), token, clientID)
		if err == nil {
			break
		}
	}

	if err != nil {
		log.Printf("Failed to validate Google token: %v", err)
		return nil, err
	}

	googleID := payload.Subject
	email := payload.Claims["email"].(string)
	log.Printf("Google login attempt for email: %s", email)

	if !s.isEmailAllowed(email) {
		log.Printf("Email %s is not in the whitelist", email)
		return nil, ErrUnauthorizedEmail
	}

	var organizationID uuid.UUID
	user, err := s.userRepo.GetByGoogleID(googleID)
	if err != nil {
		if errors.Is(err, repository.ErrUserNotFound) {
			// Create new user
			user = &model.User{
				Email:    email,
				GoogleID: googleID,
			}
			if err := s.userRepo.Create(user); err != nil {
				log.Printf("Failed to create user: %v", err)
				return nil, err
			}

			// Create new organization
			org := &model.Organization{
				Name: email, // todo: need to change this to a different naming convention
			}
			if err := s.orgRepo.Create(org); err != nil {
				log.Printf("Failed to create organization: %v", err)
				return nil, err
			}

			if err := s.orgRepo.AddUser(org.ID, user.ID); err != nil {
				log.Printf("Failed to add user to organization: %v", err)
				return nil, err
			}

			// Use the newly created organization's ID
			organizationID = org.ID
		} else {
			log.Printf("Error retrieving user: %v", err)
			return nil, err
		}
	} else {
		log.Printf("Existing user found for email: %s", email)
		// Get the user's organization
		org, err := s.orgRepo.GetUserOrganization(user.ID)
		if err != nil {
			log.Printf("Failed to get user organization: %v", err)
			return nil, err
		}
		organizationID = org.ID
	}

	token, err = s.generateToken(user.ID)
	if err != nil {
		log.Printf("Failed to generate token: %v", err)
		return nil, err
	}
	log.Println("Successfully generated token")

	return &UserData{
		Token:          token,
		User:           *user,
		OrganizationID: organizationID,
	}, nil
}

func (s *service) isEmailAllowed(email string) bool {
	log.Printf("Checking if email is allowed: %s", email)
	for _, allowed := range s.emailWhitelist {
		if strings.HasSuffix(email, allowed) {
			log.Printf("Email %s is allowed", email)
			return true
		}
	}
	log.Printf("Email %s is not allowed", email)
	return false
}

func (s *service) VerifyToken(tokenString string) (string, error) {
	log.Printf("Received token for verification: %s", tokenString)

	parts := strings.Split(tokenString, ".")
	log.Printf("Token parts: %d", len(parts))

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			log.Printf("Unexpected signing method: %v", token.Header["alg"])
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(s.jwtSecret), nil
	})

	if err != nil {
		log.Printf("Error parsing token: %v", err)
		return "", err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		userID, ok := claims["sub"].(string)
		if !ok {
			log.Println("Invalid user ID in token")
			return "", errors.New("invalid user ID in token")
		}
		log.Printf("Token verified for user ID: %s", userID)
		return userID, nil
	}

	log.Println("Invalid token")
	return "", errors.New("invalid token")
}

func (s *service) GetUserByID(userID string) (*model.User, error) {
	log.Printf("Getting user by ID: %s", userID)
	id, err := uuid.Parse(userID)
	if err != nil {
		log.Printf("Failed to parse user ID: %v", err)
		return nil, errors.New("invalid user ID")
	}

	user, err := s.userRepo.GetByID(id)
	if err != nil {
		log.Printf("Failed to get user by ID: %v", err)
		return nil, err
	}

	return user, nil
}

func (s *service) generateToken(userID uuid.UUID) (string, error) {
	log.Printf("Generating token for user ID: %s", userID)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": userID.String(),
		"exp": time.Now().Add(time.Hour * 24).Unix(),
	})

	return token.SignedString([]byte(s.jwtSecret))
}
