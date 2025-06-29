package user

type CreateUserInput struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (i *CreateUserInput) validate() error {
	if i.Name == "" {
		return ErrNameIsRequired
	}

	if i.Email == "" {
		return ErrEmailIsRequired
	}

	if i.Password == "" {
		return ErrPasswordIsRequired
	}

	if len(i.Password) < 8 {
		return ErrPasswordTooShort
	}

	return nil
}
