package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/bits-and-blooms/bloom/v3"
	"github.com/fiatjaf/relayer/v2"
	"github.com/jmoiron/sqlx"
	"github.com/nbd-wtf/go-nostr"
	"github.com/nbd-wtf/go-nostr/nip11"
	"golang.org/x/time/rate"
)

func newRelay(cfg Config, subscriptionsDB *sqlx.DB) (*Relay, error) {
	if len(cfg.AllowedKinds) == 0 {
		cfg.AllowedKinds = defaultAllowedKinds
	}

	r := Relay{
		cfg:     cfg,
		storage: newStorage(cfg),
		updates: make(chan nostr.Event),

		subscriptionsDB: subscriptionsDB,
	}

	opts := []relayer.Option{
		relayer.WithPerConnectionLimiter(rate.Every(time.Millisecond*100), 10),
	}
	server, err := relayer.NewServer(r, opts...)
	if err != nil {
		return nil, fmt.Errorf("relayer new server: %w", err)
	}
	r.server = server

	return &r, nil
}

type Relay struct {
	cfg     Config
	server  *relayer.Server
	storage *storage
	updates chan nostr.Event

	subscriptionsDB *sqlx.DB
}

func (r *Relay) GetNIP11InformationDocument() nip11.RelayInformationDocument {
	return nip11.RelayInformationDocument{
		Name:          r.Name(),
		Description:   r.cfg.Nip11Description,
		PubKey:        r.cfg.Nip11Pubkey,
		Contact:       r.cfg.Nip11Contact,
		SupportedNIPs: []int{9, 11, 12, 15, 16, 20, 45, 78, 94},
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
	// Reject any kinds not explicitly allowed
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

	// Reject events that are too large
	jsonb, _ := json.Marshal(evt)
	if len(jsonb) > 10000 {
		log.Printf("rejected event, too large")
		return false
	}

	// Require subscription for some events
	if kindRequiresSubscription(evt.Kind) {
		if ok, err := isSubscribed(r.subscriptionsDB, evt.PubKey); !ok {
			if err != nil {
				log.Printf("isSubscribed: %v\n", err)
			}

			log.Printf("rejected event, no sub. pubkey: %v", evt.PubKey)
			return false
		}
	}

	// Kind 1's must be from Stemstr client else reference a known event.
	if evt.Kind == 1 {
		if !fromStemstrClient(evt) && !referencesExistingEvent(r.storage.seenEvents, evt) {
			log.Printf("rejected event, not from stemstr.app or referencing known event: %s", string(jsonb))
			return false
		}
	}

	// Reactions and reposts must reference a known event.
	if evt.Kind == nostr.KindReaction || evt.Kind == 6 || evt.Kind == 16 {
		if !referencesExistingEvent(r.storage.seenEvents, evt) {
			fmt.Printf("rejected event, does not reference known event: %v\n", string(jsonb))
			return false
		}
	}

	// 1808s are only allowed from Stemstr client.
	if evt.Kind == 1808 && !fromStemstrClient(evt) {
		log.Printf("rejected event, not from stemstr.app: %s", string(jsonb))
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

var defaultAllowedKinds = []int{
	nostr.KindSetMetadata,
	nostr.KindTextNote,
	nostr.KindContactList,
	nostr.KindReaction,
	1808,  // Stemstr Music Track
	6,     // NIP-18: Repost
	16,    // NIP-18: Generic Repost
	9735,  // NIP-57: Zap Receipt
	30078, // NIP-78: Application-specific Data
	1063,  // NIP-94: File Metadata
}

func fromStemstrClient(event *nostr.Event) bool {
	clientTag := event.Tags.GetFirst([]string{"client"})
	if clientTag == nil {
		return false
	}

	return strings.EqualFold(clientTag.Value(), "stemstr.app")
}

// referencesExistingEvent returns true if the given event has an e tag
// to an event in the provided bloom filter.
func referencesExistingEvent(f *bloom.BloomFilter, event *nostr.Event) bool {
	// Has no e tags, cannot reference existing event
	eTags := event.Tags.GetAll([]string{"e"})
	if eTags == nil || len(eTags) == 0 {
		return false
	}

	for _, tag := range eTags {
		referencedEventID := tag.Value()
		if referencedEventID == "" {
			// e tag has no value, probably a malformed event.
			// skip it and check the others
			continue
		}

		if f.Test([]byte(referencedEventID)) {
			// This new event DOES reference an existing event
			return true
		}
	}

	return false
}
