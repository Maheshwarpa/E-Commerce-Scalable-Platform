package UserService

type LoginCred struct {
	UserName string `json:"username"`
	Password string `json:"password"`
}

var LoginCredList []LoginCred
