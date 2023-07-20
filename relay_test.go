package main

import (
	"context"
	"testing"

	"github.com/nbd-wtf/go-nostr"
	"github.com/stretchr/testify/assert"
)

func TestPubkeyIsAllowed(t *testing.T) {
	var tests = []struct {
		name     string
		pubkeys  []string
		pubkey   string
		expected bool
	}{
		{"default open", []string{}, "123", true},
		{"whitelisted", []string{"123"}, "123", true},
		{"whitelisted", []string{"123", "456"}, "123", true},
		{"not whitelisted", []string{"123", "456"}, "789", false},
	}

	for _, tt := range tests {
		result := pubkeyIsAllowed(tt.pubkeys, tt.pubkey)
		assert.Equal(t, tt.expected, result)
	}
}

func TestAcceptEvent(t *testing.T) {
	var tests = []struct {
		name     string
		relay    Relay
		event    nostr.Event
		accepted bool
	}{
		{
			"kind 1 from allowed pubkey",
			Relay{
				cfg: Config{
					AllowedKinds:   defaultAllowedKinds,
					AllowedPubkeys: []string{"123"},
				},
			},
			nostr.Event{
				Kind:   1,
				PubKey: "123",
			},
			true,
		},
		{
			"kind 1808 from allowed pubkey",
			Relay{
				cfg: Config{
					AllowedKinds:   defaultAllowedKinds,
					AllowedPubkeys: []string{"123"},
				},
			},
			nostr.Event{
				Kind:   1808,
				PubKey: "123",
			},
			true,
		},
		{
			"kind 1 from not allowed pubkey",
			Relay{
				cfg: Config{
					AllowedKinds:   defaultAllowedKinds,
					AllowedPubkeys: []string{"123"},
				},
			},
			nostr.Event{
				Kind:   1,
				PubKey: "456",
			},
			false,
		},
		{
			"kind 1808 from not allowed pubkey",
			Relay{
				cfg: Config{
					AllowedKinds:   defaultAllowedKinds,
					AllowedPubkeys: []string{"123"},
				},
			},
			nostr.Event{
				Kind:   1808,
				PubKey: "456",
			},
			false,
		},
	}

	ctx := context.Background()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.relay.AcceptEvent(ctx, &tt.event)
			assert.Equal(t, tt.accepted, result)
		})
	}
}
