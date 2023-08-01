package main

import (
	"testing"

	"github.com/nbd-wtf/go-nostr"
	"github.com/stretchr/testify/assert"
)

func TestFromStemstrClient(t *testing.T) {
	var tests = []struct {
		name     string
		event    *nostr.Event
		expected bool
	}{
		{
			name: "with client tag",
			event: &nostr.Event{
				Tags: nostr.Tags{
					nostr.Tag{"client", "stemstr.app"},
				},
			},
			expected: true,
		},
		{
			name: "without client tag",
			event: &nostr.Event{
				Tags: nostr.Tags{},
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp := fromStemstrClient(tt.event)
			assert.Equal(t, tt.expected, resp)
		})
	}
}
