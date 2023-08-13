package main

import (
	"context"

	"github.com/fiatjaf/relayer/v2/storage/postgresql"
	"github.com/nbd-wtf/go-nostr"
)

func newStorage(cfg Config) *storage {
	return &storage{
		&postgresql.PostgresBackend{
			DatabaseURL:       cfg.DatabaseURL,
			QueryLimit:        1000,
			QueryAuthorsLimit: 1000,
			QueryIDsLimit:     1000,
			QueryKindsLimit:   10,
			QueryTagsLimit:    20,
		},
	}
}

type storage struct {
	*postgresql.PostgresBackend
}

func (s *storage) BeforeSave(ctx context.Context, evt *nostr.Event) {
}

func (s *storage) AfterSave(evt *nostr.Event) {
}
