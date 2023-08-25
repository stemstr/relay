package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	"github.com/bits-and-blooms/bloom/v3"
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
	cfg        Config
	blastr     blastrIface
	seenEvents *bloom.BloomFilter
}

type blastrIface interface {
	Send(context.Context, nostr.Event) error
}

func (s *storage) Init() error {
	// First call the shadowed relayer Init
	if err := s.PostgresBackend.Init(); err != nil {
		return err
	}

	// Now do our own init
	if s.cfg.BloomFilterSize > 0 && s.cfg.BloomFilterFP > 0 {
		log.Printf("bloom filter size: %v fp: %v\n", s.cfg.BloomFilterSize, s.cfg.BloomFilterFP)
		s.seenEvents = bloom.NewWithEstimates(s.cfg.BloomFilterSize, s.cfg.BloomFilterFP)
	} else {
		log.Printf("defaulting bloom filter size: 1,000,000 fp: 0.01\n")
		s.seenEvents = bloom.NewWithEstimates(1_000_000, 0.01)
	}

	if err := s.initSeenEvents(); err != nil {
		return fmt.Errorf("initSeenEvents: %w", err)
	}

	return nil
}

func (s *storage) BeforeSave(ctx context.Context, event *nostr.Event) {
}

func (s *storage) AfterSave(event *nostr.Event) {
	// Update the Bloom Filter
	s.seenEvents.Add([]byte(event.ID))

	switch event.Kind {
	case 1808:
		shareEvent := generateShareEvent(event)
		if shareEvent != nil && s.blastr != nil {
			s.blastr.Send(context.Background(), *shareEvent)
		}
	}
}

func (s *storage) initSeenEvents() error {
	ids, err := s.getAllEventIDs()
	if err != nil {
		return err
	}

	for _, id := range ids {
		s.seenEvents.Add([]byte(id))
	}

	log.Printf("seenFilter count: %d\n", s.seenEvents.ApproximatedSize())
	return nil
}

func (s *storage) getAllEventIDs() ([]string, error) {
	var ids []string
	err := s.DB.Select(&ids, "SELECT id FROM event")
	if err != nil && err != sql.ErrNoRows {
		return nil, fmt.Errorf("failed to select events: %w", err)
	}

	return ids, nil
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
		"🎉 Let's gooo! nostr:%s just shared a track on nostr:%s 🙌.\n\nCheck it out at https://stemstr.app/thread/%s/\n\n#stemstr #music #tunestr",
		npub, stemstrNpub, event.ID,
	)

	shareEvent := nostr.Event{
		Kind: nostr.KindTextNote,
		Tags: nostr.Tags{
			{"p", event.PubKey},  // Tag author
			{"p", stemstrHexpub}, // Tag stemstr
			{"t", "stemstr"},
			{"t", "music"},
			{"t", "tunestr"},
		},
		Content: content,
	}

	return &shareEvent
}
