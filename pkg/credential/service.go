package credential

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"reflect"
	"strconv"
	"strings"
	"time"
	"xapiens/pkg/entities"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

type Service interface {
	InsertCredential(c *gin.Context) (*entities.Credential, error)
	FetchCredentials() (*[]entities.Credential, error)
	DetailCredential(id string) (*entities.Credential, error)
	UpdateCredential(c *gin.Context) (*entities.Credential, error)
	DeleteCredential(id string) error
	Validate(interface{}) (bool, error)
	Login(c *gin.Context) (*entities.JwtToken, error)
	IsValidToken(c *gin.Context) error
	TokenValid(r *http.Request) error
	VerifyToken(r *http.Request) (*jwt.Token, error)
	ExtractToken(r *http.Request) string
}

type service struct {
	repository Repository
}

func NewCredentialService(r Repository) Service {
	return &service{
		repository: r,
	}
}

func (s *service) InsertCredential(c *gin.Context) (*entities.Credential, error) {
	var credential entities.Credential
	bind := c.ShouldBindJSON(&credential)
	if bind != nil {
		return &credential, bind
	}

	validate, errValidate := s.Validate(credential)
	if !validate {
		return &credential, errValidate
	}

	exists, errExists := s.repository.FindByClientKey(credential.ClientKey, "")
	if errExists != nil {
		return &credential, errExists
	}
	if exists.CredentialId != "" {
		return &credential, errors.New("clientKey already taken")
	}

	return s.repository.CreateCredential(&credential)
}

func (s *service) FetchCredentials() (*[]entities.Credential, error) {
	return s.repository.FindAllCredentials()
}

func (s *service) DetailCredential(id string) (*entities.Credential, error) {
	var credential entities.Credential
	if id == "" {
		return &credential, errors.New("credentialId required")
	}
	exists, err := s.repository.FindByIdCredential(id)
	if err != nil {
		return &credential, err
	}

	if exists.CredentialId == "" {
		return &credential, errors.New("data not found")
	}
	return exists, nil
}

func (s *service) UpdateCredential(c *gin.Context) (*entities.Credential, error) {
	var credential entities.Credential
	bind := c.ShouldBindJSON(&credential)
	if bind != nil {
		return &credential, bind
	}

	if credential.CredentialId == "" {
		return &credential, errors.New("credentialId required")
	}

	// check data is exists

	exists, err := s.repository.FindByIdCredential(credential.CredentialId)
	if err != nil {
		return &credential, err
	}

	if exists.CredentialId == "" {
		return &credential, errors.New("data not found")
	}

	validate, errValidate := s.Validate(credential)
	if !validate {
		return &credential, errValidate
	}

	exists, errExists := s.repository.FindByClientKey(credential.ClientKey, credential.CredentialId)
	if errExists != nil {
		return &credential, errExists
	}
	if exists.CredentialId != "" {
		return &credential, errors.New("clientKey already taken")
	}

	return s.repository.UpdateCredential(&credential)
}

func (s *service) DeleteCredential(id string) error {

	if id == "" {
		return errors.New("credentialId required")
	}

	// check data is exists

	exists, err := s.repository.FindByIdCredential(id)
	if err != nil {
		return err
	}

	if exists.CredentialId == "" {
		return errors.New("data not found")
	}

	return s.repository.DeleteCredential(id)
}

func (s *service) Validate(data interface{}) (bool, error) {
	t := reflect.TypeOf(data)
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		if field.Tag.Get("required") == "true" {
			value := reflect.ValueOf(data).Field(i).Interface()
			if value == "" {
				return false, errors.New(reflect.TypeOf(data).Field(i).Name + " should not be empty")
			}
		}

		if field.Tag.Get("max") != "" {
			max := reflect.ValueOf(data).Field(i).Interface()
			maxLength, _ := strconv.Atoi(field.Tag.Get("max"))
			if len(max.(string)) > maxLength {
				return false, errors.New(reflect.TypeOf(data).Field(i).Name + " should not more than " + field.Tag.Get("max") + " character")
			}
		}
	}

	return true, nil
}

func (s *service) Login(c *gin.Context) (*entities.JwtToken, error) {
	var login entities.Login

	bind := c.ShouldBindJSON(&login)
	if bind != nil {
		return &entities.JwtToken{}, bind
	}

	validate, errValidate := s.Validate(login)
	if !validate {
		return &entities.JwtToken{}, errValidate
	}

	exists, err := s.repository.GetToken(login.ClientKey, login.SecretKey)
	if err != nil {
		return &entities.JwtToken{}, err
	}

	if exists.CredentialId == "" {
		return &entities.JwtToken{}, errors.New("data not found")
	}

	atClaims := jwt.MapClaims{}
	atClaims["scope"] = exists.Scope
	atClaims["platform"] = exists.Platform
	atClaims["exp"] = time.Now().Add(time.Minute * 60 * 60 * 24).Unix()
	at := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)
	token, err := at.SignedString([]byte(os.Getenv("JWT_SECRET")))

	if err != nil {
		return &entities.JwtToken{}, err
	}
	res := entities.JwtToken{
		Credential: *exists,
		Token:      token,
	}
	return &res, nil

}

func (s *service) IsValidToken(c *gin.Context) error {

	return s.TokenValid(c.Request)
}

func (s *service) TokenValid(r *http.Request) error {
	token, err := s.VerifyToken(r)
	if err != nil {
		return err
	}
	if _, ok := token.Claims.(jwt.Claims); !ok && !token.Valid {
		return err
	}
	return nil
}

func (s *service) VerifyToken(r *http.Request) (*jwt.Token, error) {
	tokenString := s.ExtractToken(r)
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("JWT_SECRET")), nil
	})
	if err != nil {
		return nil, err
	}
	return token, nil
}

func (s *service) ExtractToken(r *http.Request) string {
	bearToken := r.Header.Get("Authorization")
	strArr := strings.Split(bearToken, " ")
	if len(strArr) == 2 {
		return strArr[1]
	}
	return ""
}
