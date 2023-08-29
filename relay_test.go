package main

import (
	"testing"

	"github.com/bits-and-blooms/bloom/v3"
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

func TestReferencesExistingEvent(t *testing.T) {
	var tests = []struct {
		name             string
		event            *nostr.Event
		existingEventIDs []string
		expected         bool
	}{
		{
			name: "reaction on existing event",
			event: &nostr.Event{
				Kind: 7,
				Tags: nostr.Tags{
					nostr.Tag{"e", "12345"},
				},
			},
			existingEventIDs: []string{"12345", "xxxx"},
			expected:         true,
		},
		{
			name: "repost on existing event",
			event: &nostr.Event{
				Kind: 6,
				Tags: nostr.Tags{
					nostr.Tag{"e", "12345"},
				},
			},
			existingEventIDs: []string{"12345", "xxxx"},
			expected:         true,
		},
		{
			name: "generic repost on existing event",
			event: &nostr.Event{
				Kind: 16,
				Tags: nostr.Tags{
					nostr.Tag{"e", "12345"},
				},
			},
			existingEventIDs: []string{"12345", "xxxx"},
			expected:         true,
		},
		{
			name: "event does not exist",
			event: &nostr.Event{
				Kind: 16,
				Tags: nostr.Tags{
					nostr.Tag{"e", "yyy"},
				},
			},
			existingEventIDs: []string{"12345", "xxxx"},
			expected:         false,
		},
		{
			name: "no e tags",
			event: &nostr.Event{
				Kind: 16,
			},
			existingEventIDs: []string{"12345", "xxxx"},
			expected:         false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := bloom.NewWithEstimates(1000, 0.01)
			for _, id := range tt.existingEventIDs {
				f.Add([]byte(id))
			}

			resp := referencesExistingEvent(f, tt.event)
			assert.Equal(t, tt.expected, resp)
		})
	}
}
