package id

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/machinebox/graphql"
)

const (
	createIDPUserMutation = `
mutation($email: String!, $password: String!) {
	result: createUser(input: {email:$email, password: $password}) {
		id
		email
		email_verified
	}
}  
`
	changeIDPUserPasswordMutation = `
mutation($email: String!, $newPassword: String!) {
	result: changePassword(email: $email, newPassword: $password) {
		id
		email
		email_verified
	}
}  
`
)

type IDPClient struct {
	URL string
}

func NewIDPClient() *IDPClient {
	URL := os.Getenv("IDP_URL")

	if URL == "" {
		panic(fmt.Errorf("Missing required environment variable IDP_URL"))
	}
	return &IDPClient{URL}
}

type IDPUser struct {
	ID            string
	Email         string
	EmailVerified bool `json:"email_verified"`
}
type IDPUserResponse struct {
	Result IDPUser
}

func (c *IDPClient) CreateUser(ctx context.Context, email, password string) (user IDPUser, err error) {
	var res IDPUserResponse

	req := graphql.NewRequest(createIDPUserMutation)
	req.Var("email", email)
	req.Var("password", password)
	err = c.sendRequest(ctx, req, &res)

	user = res.Result

	return
}

func (c *IDPClient) ChangePassword(ctx context.Context, email, newPassword string) (user IDPUser, err error) {
	var res IDPUserResponse

	req := graphql.NewRequest(changeIDPUserPasswordMutation)
	req.Var("email", email)
	req.Var("newPassword", newPassword)
	err = c.sendRequest(ctx, req, &res)

	user = res.Result

	return
}

func (c *IDPClient) sendRequest(ctx context.Context, req *graphql.Request, data interface{}) error {
	client := graphql.NewClient(c.URL)
	client.Log = func(s string) {
		log.Println(s)
	}

	return client.Run(ctx, req, data)
}
