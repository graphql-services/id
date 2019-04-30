package id

import (
	"context"
	"fmt"

	"github.com/badoux/checkmail"
	"github.com/graphql-services/id/database"
	uuid "github.com/satori/go.uuid"
) // THIS CODE IS A STARTING POINT ONLY. IT WILL NOT BE UPDATED WITH SCHEMA CHANGES.

type Resolver struct {
	DB *database.DB
}

func (r *Resolver) Mutation() MutationResolver {
	return &mutationResolver{r}
}
func (r *Resolver) Query() QueryResolver {
	return &queryResolver{r}
}

type mutationResolver struct{ *Resolver }

func (r *mutationResolver) InviteUser(ctx context.Context, email string) (u *User, err error) {
	err = checkmail.ValidateFormat(email)
	if err != nil {
		err = fmt.Errorf("invalid email format")
		return
	}

	u = &User{}
	// TODO: search user by account emails, not just primary user email
	res := r.DB.Client().First(u, "email = ?", email)
	err = res.Error
	if err != nil && !res.RecordNotFound() {
		return
	}

	if res.RecordNotFound() {
		u = &User{
			ID:    uuid.Must(uuid.NewV4()).String(),
			Email: email,
		}
		err = r.DB.Client().Save(u).Error
	}

	return
}
func (r *mutationResolver) ForgotPassword(ctx context.Context, email string) (*ForgotPasswordRequest, error) {
	panic("not implemented")
}
func (r *mutationResolver) RegisterUser(ctx context.Context, email string, password string, info UserInfo) (*User, error) {
	panic("not implemented")
}
func (r *mutationResolver) ActivateUser(ctx context.Context, userActivationRequestID string) (bool, error) {
	panic("not implemented")
}
func (r *mutationResolver) UpdateUser(ctx context.Context, info UserInfo) (*User, error) {
	panic("not implemented")
}
func (r *mutationResolver) ResetPassword(ctx context.Context, forgotPasswordRequestID string, newPassword string) (bool, error) {
	panic("not implemented")
}
func (r *mutationResolver) UpdatePassword(ctx context.Context, oldPassword string, newPassword string) (bool, error) {
	panic("not implemented")
}

type queryResolver struct{ *Resolver }

func (r *queryResolver) User(ctx context.Context, id string) (u *User, err error) {
	u = &User{}
	res := r.DB.Client().First(u, "id = ?", id)

	if res.RecordNotFound() {
		u = nil
		return
	}
	err = res.Error

	return
}
