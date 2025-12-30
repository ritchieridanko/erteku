package models

type RequestMeta struct {
	UserAgent string
	IPAddress string
}

type SignUpRequest struct {
	Email    string
	Password string
}

type SignInRequest struct {
	Email    string
	Password string
}
