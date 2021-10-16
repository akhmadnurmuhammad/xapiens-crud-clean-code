package credential

import (
	"time"
	"xapiens/pkg/entities"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Repository interface {
	CreateCredential(credentials *entities.Credential) (*entities.Credential, error)
	FindAllCredentials() (*[]entities.Credential, error)
	FindByIdCredential(id string) (*entities.Credential, error)
	FindByClientKey(clientKey string, notId string) (*entities.Credential, error)
	UpdateCredential(credentials *entities.Credential) (*entities.Credential, error)
	DeleteCredential(id string) error
	GetToken(clientKey string, secretKey string) (*entities.Credential, error)
}

type repository struct {
	db *gorm.DB
}

func NewCredentialRepository(db *gorm.DB) Repository {
	return &repository{
		db: db,
	}
}

func (r *repository) CreateCredential(data *entities.Credential) (*entities.Credential, error) {
	data.CredentialId = uuid.NewString()
	data.CreatedAt = time.Now().Format("2006-01-02 15:04:05")
	data.UpdatedAt = time.Now().Format("2006-01-02 15:04:05")
	create := r.db.Table("credentials").Create(data)
	if create.Error != nil {
		return nil, create.Error
	}

	return data, nil
}

func (r *repository) FindAllCredentials() (*[]entities.Credential, error) {
	var credentials []entities.Credential
	find := r.db.Table("credentials").Find(&credentials)
	if find.Error != nil {
		return nil, find.Error
	}

	return &credentials, nil
}

func (r *repository) FindByIdCredential(id string) (*entities.Credential, error) {
	var credentials entities.Credential

	find := r.db.Where("credential_id", id).Table("credentials").Find(&credentials)

	if find.Error != nil {
		return nil, find.Error
	}

	return &credentials, nil
}

func (r *repository) FindByClientKey(clientKey string, notId string) (*entities.Credential, error) {
	var credentials entities.Credential

	query := r.db.Where("client_key = ?", clientKey)
	if notId != "" {
		query = query.Where("credential_id != ?", notId)
	}
	find := query.Table("credentials").Find(&credentials)

	if find.Error != nil {
		return nil, find.Error
	}

	return &credentials, nil
}

func (r *repository) UpdateCredential(data *entities.Credential) (*entities.Credential, error) {
	data.UpdatedAt = time.Now().Format("2006-01-02 15:04:05")
	update := r.db.Where("credential_id = ?", data.CredentialId).Table("credentials").Updates(data)
	if update.Error != nil {
		return nil, update.Error
	}

	return r.FindByIdCredential(data.CredentialId)
}

func (r *repository) DeleteCredential(id string) error {
	delete := r.db.Table("credentials").Where("credential_id", id).Delete(&entities.Credential{})

	if delete.Error != nil {
		return delete.Error
	}

	return nil
}

func (r *repository) GetToken(clientKey string, secretKey string) (*entities.Credential, error) {
	var credential entities.Credential

	find := r.db.Table("credentials").Where("client_key = ?", clientKey).Where("secret_key = ?", secretKey).Find(&credential)

	if find.Error != nil {
		return nil, find.Error
	}

	return r.FindByIdCredential(credential.CredentialId)
}
