package id

import (
	"context"

	uuid "github.com/gofrs/uuid"
	"github.com/graphql-services/id/database"
)

type UserStore struct {
	DB *database.DB
}

func (s *UserStore) InviteUser(ctx context.Context, email string) (u *User, new bool, err error) {
	u = &User{}
	// TODO: search user by account emails, not just primary user email
	res := s.DB.Client().First(u, "email = ?", email)
	err = res.Error
	if err != nil && !res.RecordNotFound() {
		return
	}

	if res.RecordNotFound() {
		u = &User{
			ID:    uuid.Must(uuid.NewV4()).String(),
			Email: email,
		}
		new = true
		err = s.DB.Client().Save(u).Error
	}

	return
}

func (s *UserStore) FindUserByEmail(ctx context.Context, email string) (u *User, err error) {
	u = &User{}
	res := s.DB.Client().First(u, "email = ?", email)

	if res.RecordNotFound() {
		u = nil
		return
	}
	err = res.Error

	return
}

func (s *UserStore) GetUser(ctx context.Context, id string) (u *User, err error) {
	u = &User{}
	res := s.DB.Client().First(u, "id = ?", id)

	if res.RecordNotFound() {
		u = nil
		return
	}
	err = res.Error

	return
}

func (s *UserStore) UpdateUser(ctx context.Context, id string, info *UserInfo) (u *User, err error) {
	u = &User{}
	res := s.DB.Client().First(u, "id = ?", id)

	if res.RecordNotFound() {
		u = nil
		return
	}
	err = res.Error
	if err != nil {
		return
	}

	if info != nil {
		err = s.DB.Client().Model(u).Updates(*info).Error
	}

	return
}
