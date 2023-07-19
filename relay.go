package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/fiatjaf/relayer/v2"
	"github.com/fiatjaf/relayer/v2/storage/postgresql"
	"github.com/nbd-wtf/go-nostr"
	"github.com/nbd-wtf/go-nostr/nip11"
)

var allowedKinds = []int{
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

func newRelay(port int, databaseURL, infoPubkey, infoContact, infoDesc, infoVersion string) (*Relay, error) {
	r := Relay{
		port: port,
		storage: &postgresql.PostgresBackend{
			DatabaseURL:       databaseURL,
			QueryLimit:        1000,
			QueryAuthorsLimit: 1000,
			QueryIDsLimit:     1000,
			QueryKindsLimit:   10,
			QueryTagsLimit:    20,
		},
		updates:     make(chan nostr.Event),
		infoPubkey:  infoPubkey,
		infoContact: infoContact,
		infoDesc:    infoDesc,
		infoVersion: infoVersion,
	}

	if err := r.storage.Init(); err != nil {
		return nil, fmt.Errorf("relay init: %w", err)
	}

	return &r, nil
}

type Relay struct {
	port        int
	storage     *postgresql.PostgresBackend
	updates     chan nostr.Event
	infoPubkey  string
	infoContact string
	infoDesc    string
	infoVersion string
}

func (r *Relay) GetNIP11InformationDocument() nip11.RelayInformationDocument {
	supportedNIPs := []int{9, 11, 12, 15, 16, 20, 78, 94}
	return nip11.RelayInformationDocument{
		Name:          r.Name(),
		Description:   r.infoDesc,
		PubKey:        r.infoPubkey,
		Contact:       r.infoContact,
		SupportedNIPs: supportedNIPs,
		Software:      "https://github.com/Stemstr",
		Version:       r.infoVersion,
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
	// block events that are too large
	jsonb, _ := json.Marshal(evt)
	if len(jsonb) > 10000 {
		return false
	}

	allowed := false
	for _, kind := range allowedKinds {
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
	server, err := relayer.NewServer(r)
	if err != nil {
		return fmt.Errorf("relayer new server: %w", err)
	}

	log.Printf("listening on 0.0.0.0:%v\n", r.port)
	return server.Start("0.0.0.0", r.port)
}
