package app

import (
	"context"
	"time"

	"github.com/sashabaranov/go-openai"
)

type moderator struct {
	c *openai.Client
}

const moderationTimeout = time.Second * 2

func newModerator(c *openai.Client) *moderator {
	return &moderator{c: c}
}

func (m *moderator) Moderate(input string) bool {
	ctx, cancel := context.WithTimeout(context.Background(), moderationTimeout)
	defer cancel()
	resp, err := m.c.Moderations(
		ctx,
		openai.ModerationRequest{Input: input},
	)
	if err != nil {
		// If we can't check the request, assume it is OK.
		return true
	}

	for _, result := range resp.Results {
		if result.Flagged {
			return false
		}
	}
	return true
}
