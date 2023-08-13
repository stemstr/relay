package main

import (
	"context"
	"testing"

	"github.com/nbd-wtf/go-nostr"
	"github.com/stretchr/testify/assert"
)

func TestAfterSave(t *testing.T) {
	const sourceEventPubkey = "e6230af84359c4bedd38207a349dfbf8d25ad3b767518e7231fcae03907b9282"

	var tests = []struct {
		name            string
		sourceEvent     *nostr.Event
		expectedContent string
		expectedTags    nostr.Tags
	}{
		{
			name: "handle 1808",
			sourceEvent: &nostr.Event{
				ID:      "source.event.id",
				PubKey:  sourceEventPubkey,
				Kind:    1808,
				Content: "a sample 1808",
			},
			expectedContent: "🎉 Let's gooo! nostr:npub1uc3s47zrt8ztahfcyparf80mlrf945ahvagcuu33ljhq8yrmj2pqefmzr7 just shared a track on nostr:npub1stemstrls4f5plqeqkeq43gtjhtycuqd9w25v5r5z5ygaq2n2sjsd6mul5 🙌.\nCheck it out at https://stemstr.app/thread/source.event.id",
			expectedTags: nostr.Tags{
				{"p", sourceEventPubkey}, // Tag author
				{"p", stemstrHexpub},     // Tag stemstr
			},
		},
		{
			name: "ignore 1",
			sourceEvent: &nostr.Event{
				ID:      "source.event.id",
				PubKey:  sourceEventPubkey,
				Kind:    1,
				Content: "a sample kind 1",
			},
			expectedContent: "",
			expectedTags:    nostr.Tags{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			store := &storage{
				blastr: &mockBlastr{
					sendAsserter: func(event nostr.Event) {
						// Make sure the message is properly rendered
						assert.Equal(t, tt.expectedContent, event.Content)
						// And the mentions are tagged
						assert.Equal(t, tt.expectedTags, event.Tags)
					},
				},
			}

			store.AfterSave(tt.sourceEvent)
		})
	}
}

type mockBlastr struct {
	sendAsserter func(nostr.Event)
}

func (m *mockBlastr) Send(ctx context.Context, event nostr.Event) error {
	m.sendAsserter(event)
	return nil
}
