package credential

import (
	"xapiens/pkg/entities"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Repository interface {
	CreateCredential(credentials *entities.Credential) (*entities.Credential, error)
	FindAllCredentials() (*[]entities.Credential, error)
	FindByIdCredential(id string) (*entities.Credential, error)
	UpdateCredential(credentials *entities.Credential) (*entities.Credential, error)
	DeleteCredential(id string) error
}

type repository struct {
	db *gorm.DB
}

func NewCredentialRepository(db *gorm.DB) Repository {
	return &repository{
		db: db,
	}
}

func (r repository) CreateCredential(data *entities.Credential) (*entities.Credential, error) {
	data.CredentialId = uuid.NewString()

	create := r.db.Create(data)
	if create.Error != nil {
		return nil, create.Error
	}

	return data, nil
}

func (r repository) FindAllCredentials() (*[]entities.Credential, error) {
	var credentials []entities.Credential
	find := r.db.Find(&credentials)
	if find.Error != nil {
		return nil, find.Error
	}

	return &credentials, nil
}

func (r repository) FindByIdCredential(id string) (*entities.Credential, error) {
	var credentials entities.Credential

	find := r.db.Where(&entities.Credential{
		CredentialId: id,
	}).Find(&credentials)

	if find.Error != nil {
		return nil, find.Error
	}

	return &credentials, nil
}

func (r repository) UpdateCredential(data *entities.Credential) (*entities.Credential, error) {
	update := r.db.Where(&entities.Credential{
		CredentialId: data.CredentialId,
	}).Updates(data)
	if update.Error != nil {
		return nil, update.Error
	}

	return data, nil
}

func (r repository) DeleteCredential(id string) error {
	delete := r.db.Model(&entities.Credential{}).Where(&entities.Credential{
		CredentialId: id,
	}).Delete(&entities.Credential{})

	if delete.Error != nil {
		return delete.Error
	}

	return nil
}
