package domain

import (
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type Password string

func (p Password) Hash() (HashedPassword, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(p), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	return HashedPassword(hashedPassword), nil
}

type HashedPassword []byte

type User struct {
	id       uuid.UUID
	name     string
	email    string
	password HashedPassword
}

func NewUser(name string, email string, password Password) (*User, error) {
	hashedPassword, err := password.Hash()
	if err != nil {
		return nil, err
	}

	return &User{id: uuid.New(), name: name, email: email, password: hashedPassword}, nil
}

func (u *User) ID() uuid.UUID {
	return u.id
}

func (u *User) Name() string {
	return u.name
}

func (u *User) Email() string {
	return u.email
}

func (u *User) Password() HashedPassword {
	return u.password
}

func UserFromSource(id uuid.UUID, name string, email string) *User {
	return &User{
		id:    id,
		name:  name,
		email: email,
	}
}
