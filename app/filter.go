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
			// Always reject on these ones.
			if result.Categories.HateThreatening ||
				result.Categories.SelfHarm ||
				result.Categories.Sexual ||
				result.Categories.SexualMinors {
				return false
			}
			// Sometimes you might have a creative conversation, like "an alien threatens to shoot everyone".
			// Ideally we don't flag it, because chat handles this pretty well anyway.
			cumulativeScore := result.CategoryScores.Hate + result.CategoryScores.HateThreatening +
				result.CategoryScores.Violence + result.CategoryScores.ViolenceGraphic
			if cumulativeScore < 1.0 {
				return true
			} else {
				return false
			}
		}
	}
	return true
}
