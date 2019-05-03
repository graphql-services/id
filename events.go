package id

import (
	"context"
	"fmt"
	"os"

	cloudevents "github.com/cloudevents/sdk-go"
)

const (
	UserInvitedEvent    = "com.graphql.id.user.invited"
	ForgotPasswordEvent = "com.graphql.id.user.forgot-password"
)

type EventController struct {
	client *cloudevents.Client
}

func NewEventController() (ec EventController, err error) {
	URL := os.Getenv("EVENT_TRANSPORT_URL")
	var _client *cloudevents.Client
	if URL != "" {
		t, tErr := cloudevents.NewHTTPTransport(
			cloudevents.WithTarget(URL),
			cloudevents.WithBinaryEncoding(),
		)
		err = tErr
		if err != nil {
			return
		}

		client, cErr := cloudevents.NewClient(t)
		err = cErr
		if err != nil {
			return
		}
		_client = &client
	}
	ec = EventController{_client}
	return
}

func (c *EventController) send(ctx context.Context, e cloudevents.Event) error {
	if c.client == nil {
		return nil
	}
	_, err := (*c.client).Send(ctx, e)
	return err
}

type UserInvitationRequestEvent struct {
	RequestID string `json:"requestID"`
	User      User   `json:"user"`
}

func (c *EventController) SendUserInvitationRequest(ctx context.Context, r *UserInvitationRequest, u *User) (err error) {
	event := cloudevents.NewEvent()
	event.SetID(r.ID)
	event.SetType(UserInvitedEvent)
	event.SetSource(fmt.Sprintf("http://graphql-id/invited-user/%s", r.UserID))
	event.SetData(UserInvitationRequestEvent{
		RequestID: r.ID,
		User:      *u,
	})

	err = c.send(ctx, event)
	return
}

func (c *EventController) SendForgotPasswordRequest(ctx context.Context, r *ForgotPasswordRequest, u *User) (err error) {
	event := cloudevents.NewEvent()
	event.SetID(r.ID)
	event.SetType(ForgotPasswordEvent)
	event.SetSource(fmt.Sprintf("http://graphql-id/forgot-password/%s", r.UserID))
	event.SetData(UserInvitationRequestEvent{
		RequestID: r.ID,
		User:      *u,
	})

	err = c.send(ctx, event)
	return
}
