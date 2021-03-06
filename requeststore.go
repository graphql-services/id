package id

import (
	"fmt"

	uuid "github.com/gofrs/uuid"
	"github.com/graphql-services/id/database"
)

type RequestStore struct {
	DB *database.DB
}

func (s *RequestStore) CreateInvitationRequest(userId string) (r *UserInvitationRequest, err error) {
	r = &UserInvitationRequest{
		ID:     uuid.Must(uuid.NewV4()).String(),
		UserID: userId,
	}
	fmt.Println("invitation request", r)

	err = s.DB.Client().Save(r).Error

	return
}
func (s *RequestStore) GetInvitationRequest(id string) (r *UserInvitationRequest, err error) {
	r = &UserInvitationRequest{}
	err = s.DB.Client().First(r, "id = ?", id).Error
	return
}
func (s *RequestStore) DeleteInvitationRequest(id string) (r *UserInvitationRequest, err error) {
	r = &UserInvitationRequest{}
	err = s.DB.Client().Delete(r, "id = ?", id).Error
	return
}

func (s *RequestStore) CreateActivationRequest(userId string) (r *UserActivationRequest, err error) {
	r = &UserActivationRequest{
		ID:     uuid.Must(uuid.NewV4()).String(),
		UserID: userId,
	}

	err = s.DB.Client().Save(r).Error

	return
}

func (s *RequestStore) CreateForgotPasswordRequest(userId string) (r *ForgotPasswordRequest, err error) {
	r = &ForgotPasswordRequest{
		ID:     uuid.Must(uuid.NewV4()).String(),
		UserID: userId,
	}

	err = s.DB.Client().Save(r).Error

	return
}

func (s *RequestStore) GetForgotPasswordRequest(id string) (r *ForgotPasswordRequest, err error) {
	r = &ForgotPasswordRequest{}
	err = s.DB.Client().First(r, "id = ?", id).Error
	return
}
func (s *RequestStore) DeleteForgotPasswordRequest(id string) (r *ForgotPasswordRequest, err error) {
	r = &ForgotPasswordRequest{}
	err = s.DB.Client().Delete(r, "id = ?", id).Error
	return
}
