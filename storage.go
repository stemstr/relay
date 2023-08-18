package main

import (
	"context"
	"fmt"
	"log"

	"github.com/fiatjaf/relayer/v2/storage/postgresql"
	"github.com/nbd-wtf/go-nostr"
	"github.com/nbd-wtf/go-nostr/nip19"
	"github.com/stemstr/blastr"
)

func newStorage(cfg Config) *storage {
	store := &storage{
		PostgresBackend: &postgresql.PostgresBackend{
			DatabaseURL:       cfg.DatabaseURL,
			QueryLimit:        1000,
			QueryAuthorsLimit: 1000,
			QueryIDsLimit:     1000,
			QueryKindsLimit:   10,
			QueryTagsLimit:    20,
		},
		cfg: cfg,
	}

	if cfg.BlastrNsec != "" {
		store.blastr, _ = blastr.New(cfg.BlastrNsec)
	}

	return store
}

type storage struct {
	*postgresql.PostgresBackend
	cfg    Config
	blastr blastrIface
}

type blastrIface interface {
	Send(context.Context, nostr.Event) error
}

func (s *storage) BeforeSave(ctx context.Context, event *nostr.Event) {
}

func (s *storage) AfterSave(event *nostr.Event) {
	if event.Kind != 1808 {
		return
	}

	if s.blastr == nil {
		log.Println("blastr nil")
		return
	}

	shareEvent := generateShareEvent(event)
	if shareEvent != nil {
		s.blastr.Send(context.Background(), *shareEvent)
	}
}

const (
	stemstrNpub   = "npub1stemstrls4f5plqeqkeq43gtjhtycuqd9w25v5r5z5ygaq2n2sjsd6mul5"
	stemstrHexpub = "82f3b82c7f855340fc1905b20ac50b95d64c700d2b9546507415088e81535425"
)

func generateShareEvent(event *nostr.Event) *nostr.Event {
	npub, err := nip19.EncodePublicKey(event.PubKey)
	if err != nil {
		log.Printf("failed to encoded share event npub: %v", err)
		return nil
	}

	content := fmt.Sprintf(
		"🎉 Let's gooo! nostr:%s just shared a track on nostr:%s 🙌.\nCheck it out at https://stemstr.app/thread/%s/\n#music #tunestr",
		npub, stemstrNpub, event.ID,
	)

	shareEvent := nostr.Event{
		Kind: nostr.KindTextNote,
		Tags: nostr.Tags{
			{"p", event.PubKey},  // Tag author
			{"p", stemstrHexpub}, // Tag stemstr
			{"t", "music"},
			{"t", "tunestr"},
		},
		Content: content,
	}

	return &shareEvent
}
