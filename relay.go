package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/fiatjaf/relayer/v2"
	"github.com/fiatjaf/relayer/v2/storage/postgresql"
	"github.com/nbd-wtf/go-nostr"
	"github.com/nbd-wtf/go-nostr/nip11"
)

func newRelay(cfg Config) (*Relay, error) {
	if len(cfg.AllowedKinds) == 0 {
		cfg.AllowedKinds = defaultAllowedKinds
	}

	r := Relay{
		cfg: cfg,
		storage: &postgresql.PostgresBackend{
			DatabaseURL:       cfg.DatabaseURL,
			QueryLimit:        1000,
			QueryAuthorsLimit: 1000,
			QueryIDsLimit:     1000,
			QueryKindsLimit:   10,
			QueryTagsLimit:    20,
		},
		updates: make(chan nostr.Event),
	}

	if err := r.storage.Init(); err != nil {
		return nil, fmt.Errorf("relay init: %w", err)
	}

	server, err := relayer.NewServer(r)
	if err != nil {
		return nil, fmt.Errorf("relayer new server: %w", err)
	}
	r.server = server

	return &r, nil
}

type Relay struct {
	cfg     Config
	server  *relayer.Server
	storage *postgresql.PostgresBackend
	updates chan nostr.Event
}

func (r *Relay) GetNIP11InformationDocument() nip11.RelayInformationDocument {
	return nip11.RelayInformationDocument{
		Name:          r.Name(),
		Description:   r.cfg.Nip11Description,
		PubKey:        r.cfg.Nip11Pubkey,
		Contact:       r.cfg.Nip11Contact,
		SupportedNIPs: []int{9, 11, 12, 15, 16, 20, 78, 94},
		Software:      "https://github.com/Stemstr",
		Version:       r.cfg.Nip11Version,
	}
}

func (r Relay) Name() string {
	return "Stemstr relay"
}

func (r Relay) Storage(ctx context.Context) relayer.Storage {
	return r.storage
}

func (r Relay) OnInitialized(*relayer.Server) {}

func (r Relay) Init() error {
	return nil
}

func (r Relay) AcceptEvent(ctx context.Context, evt *nostr.Event) bool {
	// NOTE: beta release gates kind: 1808,1 to only whitelisted users.
	if evt.Kind == 1808 || evt.Kind == 1 {
		if !pubkeyIsAllowed(r.cfg.AllowedPubkeys, evt.PubKey) {
			return false
		}
	}

	// block events that are too large
	jsonb, _ := json.Marshal(evt)
	if len(jsonb) > 10000 {
		return false
	}

	allowed := false
	for _, kind := range r.cfg.AllowedKinds {
		if evt.Kind == kind {
			allowed = true
			break
		}
	}
	if !allowed {
		return false
	}

	fmt.Printf("relay: received event: %v\n", string(jsonb))

	return true
}

func (relay Relay) InjectEvents() chan nostr.Event {
	return relay.updates
}

func (r Relay) Start() error {
	log.Printf("listening on 0.0.0.0:%v\n", r.cfg.Port)
	return r.server.Start("0.0.0.0", r.cfg.Port)
}

func pubkeyIsAllowed(pubkeys []string, pubkey string) bool {
	// If no whitelist of pubkeys are provided, it's allowed
	if len(pubkeys) == 0 {
		return true
	}

	allowed := false
	for _, allowedPubkey := range pubkeys {
		if strings.EqualFold(allowedPubkey, pubkey) {
			allowed = true
			break
		}
	}

	return allowed
}

var defaultAllowedKinds = []int{
	nostr.KindSetMetadata,
	nostr.KindTextNote,
	nostr.KindContactList,
	nostr.KindBoost,
	nostr.KindReaction,
	1808,  // Stemstr Music Track
	6,     // NIP-18: Repost
	18,    // NIP-18: Generic Repost
	9735,  // NIP-57: Zaps
	30078, // NIP-78: Application-specific Data
	1063,  // NIP-94: File Metadata
}
