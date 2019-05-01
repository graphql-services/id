package id

import (
	"context"
	"fmt"
	"os"

	cloudevents "github.com/cloudevents/sdk-go"
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
			cloudevents.WithStructuredEncoding(),
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

func (c *EventController) SendUserInvitationRequest(ctx context.Context, r *UserInvitationRequest) (err error) {
	event := cloudevents.NewEvent()
	event.SetID(r.ID)
	event.SetType("com.graphql.id.user.invited")
	event.SetSource(fmt.Sprintf("http://graphql-id/invited-user/%s", r.UserID))
	event.SetData(r)

	err = c.send(ctx, event)
	return
}
