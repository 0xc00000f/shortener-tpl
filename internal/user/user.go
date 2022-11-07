package user

import "github.com/google/uuid"

type User struct {
	userID uuid.UUID
}

func New() User {
	return User{userID: uuid.New()}
}
