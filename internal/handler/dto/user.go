package dto

type SignUpDTO struct {
	Login    string `json:"login" binding:"required,max=64,min=4"`
	Password string `json:"password" binding:"required,max=64,min=8"`
}

type LogInDTO struct {
	Login    string `json:"login" binding:"required,max=64,min=4"`
	Password string `json:"password" binding:"required,max=64,min=8"`
}
