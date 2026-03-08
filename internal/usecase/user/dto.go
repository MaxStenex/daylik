package user

type RegisterInput struct {
	Email    string
	Password string
}

type RegisterOutput struct {
	ID    string
	Email string
}
