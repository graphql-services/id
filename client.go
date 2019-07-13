package id

import (
	"context"

	"github.com/machinebox/graphql"
)

// Client ...
type Client struct {
	gc *graphql.Client
}

// NewClient ...
func NewClient(URL string) *Client {
	client := graphql.NewClient(URL)
	return &Client{client}
}

func (c *Client) run(ctx context.Context, req *graphql.Request, data interface{}) error {
	return c.gc.Run(ctx, req, data)
}

const (
	inviteUserQuery = `
mutation($email: String!) {
	result: inviteUser(email: $email) {
		id
		email
		given_name
		family_name
		middle_name
	}
}  
`
)

// InvitedUser ...
type InvitedUser struct {
	ID         string `json:"id" gorm:"primary_key"`
	Name       string `json:"name"`
	Email      string `json:"email"`
	GivenName  string `json:"given_name"`
	FamilyName string `json:"family_name"`
	MiddleName string `json:"middle_name"`
}

// UserProviderInviteResponse ...
type UserProviderInviteResponse struct {
	Result InvitedUser
}

// InviteUser ...
func (c *Client) InviteUser(ctx context.Context, email string) (user InvitedUser, err error) {
	var res UserProviderInviteResponse

	req := graphql.NewRequest(inviteUserQuery)
	req.Var("email", email)
	err = c.run(ctx, req, &res)

	user = res.Result

	return
}
