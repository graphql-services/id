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

func (r *mutationResolver) InviteUser(ctx context.Context, email string, userInfo *UserInfo) (u *User, err error) {
	err = checkmail.ValidateFormat(email)
	if err != nil {
		err = fmt.Errorf("invalid email format")
		return
	}

	u, isnew, err := r.UserStore.InviteUser(ctx, email, userInfo)
	if err != nil {
		return
	}

	if isnew {
		request, requestErr := r.RequestStore.CreateInvitationRequest(u.ID)
		err = requestErr
		if err != nil {
			return
		}

		err = r.EventController.SendUserInvitationRequest(ctx, request, u)
	}

	return
}
func (r *mutationResolver) ForgotPassword(ctx context.Context, email string) (res bool, err error) {
	err = checkmail.ValidateFormat(email)
	if err != nil {
		err = fmt.Errorf("invalid email format")
		return
	}

	u, err := r.UserStore.FindUserByEmail(ctx, email)
	if err != nil {
		return
	}

	request, requestErr := r.RequestStore.CreateForgotPasswordRequest(u.ID)
	err = requestErr
	if err != nil {
		return
	}

	err = r.EventController.SendForgotPasswordRequest(ctx, request, u)

	return
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
func (r *mutationResolver) ResetPassword(ctx context.Context, requestID string, newPassword string) (res bool, err error) {
	req, err := r.RequestStore.GetForgotPasswordRequest(requestID)
	if err != nil {
		return
	}

	u, err := r.UserStore.GetUser(ctx, req.UserID)
	if err != nil {
		return
	}

	_, err = r.IDPClient.ChangePassword(ctx, u.Email, newPassword)
	if err != nil {
		return
	}

	_, err = r.RequestStore.DeleteForgotPasswordRequest(requestID)

	if err == nil {
		res = true
	}
	return
}
func (r *mutationResolver) UpdateUser(ctx context.Context, info UserInfo) (*User, error) {
	panic("not implemented")
}
func (r *mutationResolver) UpdatePassword(ctx context.Context, oldPassword string, newPassword string) (bool, error) {
	panic("not implemented")
}

type queryResolver struct{ *Resolver }

func (r *queryResolver) _service(ctx context.Context) (*_Service, error) {
	sdl := SchemaSDL
	return &_Service{
		Sdl: &sdl,
	}, nil
}
func (r *queryResolver) _entities(ctx context.Context, representations []interface{}) (res []_Entity, err error) {
	res = []_Entity{}

	keys := []string{}

	for _, repr := range representations {
		values, ok := repr.(map[string]interface{})
		if !ok {
			err = fmt.Errorf("The _entities resolver received invalid representation type")
			break
		}

		typename, ok := values["__typename"].(string)
		if !ok || typename != "User" {
			err = fmt.Errorf("Unexpected typename '%s'", typename)
			continue
		}

		identifier, ok := values["id"].(string)
		if !ok {
			res = append(res, &User{})
			continue
		}

		keys = append(keys, identifier)
	}

	users, err := r.UserStore.GetUsers(ctx, keys)
	if err != nil {
		return
	}
	usersMap := map[string]User{}
	for _, u := range users {
		usersMap[u.ID] = u
	}

	for _, key := range keys {
		user, ok := usersMap[key]
		if ok {
			res = append(res, &user)
		} else {
			res = append(res, nil)
		}
	}

	return res, nil
}
func (r *queryResolver) User(ctx context.Context, id string) (u *User, err error) {
	u, err = r.UserStore.GetUser(ctx, id)
	return
}
