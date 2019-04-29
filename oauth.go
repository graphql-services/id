// https://gist.github.com/sdorra/1c95de8cb80da31610d2ad767cd6f251
package id

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/machinebox/graphql"
)

const (
	inviteUserMutation = `
mutation($email: String!) {
	result: inviteUser(email: $email) {
		id
		email
	}
}  
`
)

type AuthInviteResponse struct {
	Result User
}

func InviteUser(ctx context.Context, email string) (user User, err error) {
	var res AuthInviteResponse

	req := graphql.NewRequest(inviteUserMutation)
	req.Var("email", email)
	err = sendRequest(ctx, req, &res)
	if err != nil {
		return
	}

	user = res.Result

	return
}

func sendRequest(ctx context.Context, req *graphql.Request, data interface{}) error {
	URL := os.Getenv("OAUTH_URL")

	if URL == "" {
		return fmt.Errorf("Missing required environment variable OAUTH_URL")
	}

	client := graphql.NewClient(URL)
	client.Log = func(s string) {
		log.Println(s)
	}

	return client.Run(ctx, req, data)
}
