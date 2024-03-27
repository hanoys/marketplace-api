package dto

type SignUpDTO struct {
	Login    string `json:"login" binding:"required,max=64"`
	Password string `json:"password" binding:"required,max=64"`
}

type LogInDTO struct {
	Login    string `json:"login" binding:"required,max=64"`
	Password string `json:"password" binding:"required,max=64"`
}
