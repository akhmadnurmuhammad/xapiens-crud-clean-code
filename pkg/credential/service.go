package credential

import (
	"errors"
	"reflect"
	"xapiens/pkg/entities"

	"github.com/gin-gonic/gin"
)

type Service interface {
	InsertCredential(c *gin.Context) (*entities.Credential, error)
	FetchCredentials() (*[]entities.Credential, error)
	DetailCredential(id string) (*entities.Credential, error)
	UpdateCredential(c *gin.Context) (*entities.Credential, error)
	DeleteCredential(id string) error
	Validate(credentials *entities.Credential) bool,error
}

type service struct {
	repository Repository
}

func NewCredentialService(r Repository) Service {
	return &service{
		repository: r,
	}
}

func (s *service).InsertCredential(c *gin.Context) (*entities.Credential, error) {
	var credential entities.Credential
	bind := c.ShouldBindJSON(&credential)
	if bind != nil {
		return &credential,bind
	}

	validate,errValidate := s.Validate(credential)
	if !validate {
		return &credential,errValidate
	}

	return s.repository.CreateCredential(&credential)
}

func (s *service).FetchCredentials() (*[]entities.Credential, error) {
	return s.repository.FindAllCredentials()
}

func (s *service).DetailCredential(id string) (*entities.Credential, error) {
	if id == "" {
		return &credential,errors.New("credentialId required.")
	}
	return s.repository.FindByIdCredential()
}

func (s *service).UpdateCredential(c *gin.Context) (*entities.Credential, error) {
	var credential entities.Credential
	bind := c.ShouldBindJSON(&credential)
	if bind != nil {
		return &credential,bind
	}

	if credential.CredentialId == "" {
		return &credential,errors.New("credentialId required.")
	}

	// check data is exists

	exists,err := s.repository.FindByIdCredential(credential.CredentialId)
	if err != nil {
		return &credential,err
	}

	if exists.CredentialId == "" {
		return &credential,errors.New("data not found.")
	}

	validate,errValidate := s.Validate(credential)
	if !validate {
		return &credential,errValidate
	}

	return s.repository.UpdateCredential(&credential)
}

func (s *service).DeleteCredential(id string) error {

	if id == "" {
		return errors.New("credentialId required.")
	}

	// check data is exists

	exists,err := s.repository.FindByIdCredential(id)
	if err != nil {
		return err
	}

	if exists.CredentialId == "" {
		return errors.New("data not found.")
	}

	return s.repository.DeleteCredential(id)
}

func (s *service) Validate(data interface{}) (bool,error) {
	t := reflect.TypeOf(data)
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		if field.Tag.Get("required") == "true" {
			value := reflect.ValueOf(data).Field(i).Interface()
			if value == "" {
				return false,errors.New(reflect.TypeOf(data).Field(i).Name + " should not be empty.")
			}
		}

		if field.Tag.Get("max") != "" {
			max := reflect.ValueOf(data).Field(i).Interface()
			maxLength, _ := strconv.Atoi(field.Tag.Get("max"))
			if len(max.(string)) > maxLength {
				return false,errors.New(reflect.TypeOf(data).Field(i).Name + " should not more than " + field.Tag.Get("max") + " character")
			}
		}
	}

	return true
}
