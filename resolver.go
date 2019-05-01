package id

import (
	"context"
	"fmt"

	"github.com/badoux/checkmail"
) // THIS CODE IS A STARTING POINT ONLY. IT WILL NOT BE UPDATED WITH SCHEMA CHANGES.

type Resolver struct {
	UserStore       *UserStore
	RequestStore    *RequestStore
	IDPClient       *IDPClient
	EventController *EventController
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

	u, isnew, err := r.UserStore.InviteUser(ctx, email)
	if err != nil {
		return
	}

	if isnew {
		request, requestErr := r.RequestStore.CreateInvitationRequest(u.ID)
		err = requestErr
		if err != nil {
			return
		}

		err = r.EventController.SendUserInvitationRequest(ctx, request)

		// TODO: send email to user with invitation and instructions
	}

	return
}
func (r *mutationResolver) ForgotPassword(ctx context.Context, email string) (bool, error) {
	panic("not implemented")
}
func (r *mutationResolver) RegisterUser(ctx context.Context, email string, password string, info *UserInfo) (*User, error) {
	panic("not implemented")
}
func (r *mutationResolver) ConfirmInvitation(ctx context.Context, requestID string, password string, info *UserInfo) (u *User, err error) {
	req, err := r.RequestStore.GetInvitationRequest(requestID)
	if err != nil {
		return
	}

	u, err = r.UserStore.UpdateUser(ctx, req.UserID, info)
	if err != nil {
		return
	}

	_, err = r.IDPClient.CreateUser(ctx, u.Email, password)
	if err != nil {
		return
	}

	_, err = r.RequestStore.DeleteInvitationRequest(requestID)

	return
}

func (r *mutationResolver) ActivateUser(ctx context.Context, requestID string, info *UserInfo) (bool, error) {
	panic("not implemented")
}
func (r *mutationResolver) ResetPassword(ctx context.Context, requestID string, newPassword string) (bool, error) {
	panic("not implemented")
}
func (r *mutationResolver) UpdateUser(ctx context.Context, info UserInfo) (*User, error) {
	panic("not implemented")
}
func (r *mutationResolver) UpdatePassword(ctx context.Context, oldPassword string, newPassword string) (bool, error) {
	panic("not implemented")
}

type queryResolver struct{ *Resolver }

func (r *queryResolver) User(ctx context.Context, id string) (u *User, err error) {
	u, err = r.UserStore.GetUser(ctx, id)
	return
}
