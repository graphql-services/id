package id

import (
	"context"
) // THIS CODE IS A STARTING POINT ONLY. IT WILL NOT BE UPDATED WITH SCHEMA CHANGES.

type Resolver struct{}

func (r *Resolver) Mutation() MutationResolver {
	return &mutationResolver{r}
}
func (r *Resolver) Query() QueryResolver {
	return &queryResolver{r}
}

type mutationResolver struct{ *Resolver }

func (r *mutationResolver) InviteUser(ctx context.Context, email string) (user *User, err error) {
	u, err := InviteUser(ctx, email)
	user = &u
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

func (r *queryResolver) User(ctx context.Context, id string) (*User, error) {
	panic("not implemented")
}
