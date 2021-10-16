package entities

type Credential struct {
	CredentialId string `json:"CredentialId"`
	ClientKey    string `json:"ClientKey" required:"true" max:"64"`
	SecretKey    string `json:"SecretKey" required:"true" max:"64"`
	Platform     string `json:"Platform" required:"true" max:"20"`
	Scope        string `json:"Scope" required:"true" max:"20"`
	CreatedAt    string
	UpdatedAt    string
	DeletedAt    interface{}
}

type JwtToken struct {
	Credential
	Token string `json:"Token"`
}

type Login struct {
	ClientKey string `json:"ClientKey" required:"true"`
	SecretKey string `json:"SecretKey" required:"true"`
}
