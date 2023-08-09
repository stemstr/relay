package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/fiatjaf/relayer/v2"
	"github.com/nbd-wtf/go-nostr"
	"github.com/nbd-wtf/go-nostr/nip19"
)

func adminHandler(cfg Config, db relayer.Storage) func(http.ResponseWriter, *http.Request) {
	const template = "admin.html"

	return func(w http.ResponseWriter, r *http.Request) {
		if !auth(r, cfg.Admins) {
			w.Header().Add("WWW-Authenticate", `Basic realm="username and password required"`)
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("Unauthorized"))
			return
		}

		var (
			id       = r.URL.Query().Get("id")
			kindStr  = r.URL.Query().Get("kind")
			pk       = r.URL.Query().Get("pubkey")
			limitStr = r.URL.Query().Get("limit")
		)

		var kind *int
		if kindStr != "" {
			if i, err := strconv.Atoi(kindStr); err == nil {
				kind = &i
			}
		}

		limit := 100
		if limitStr != "" {
			if i, err := strconv.Atoi(limitStr); err == nil {
				limit = i
			}
		}

		events, err := getEvents(db, id, pk, kind, limit)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}

		t, ok := templates[template]
		if !ok {
			log.Printf("template %s not found", template)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("missing template"))
			return
		}

		tkn := make([]byte, 16)
		rand.Read(tkn)
		nonce := fmt.Sprintf("%x", tkn)
		csp := fmt.Sprintf("script-src: 'self' 'unsafe-inline' 'nonce-%s'", nonce)

		data := map[string]any{
			"events": events,
			"nonce":  nonce,
		}

		w.Header().Set("Content-Type", "text/html")
		w.Header().Set("Content-Security-Policy", csp)
		if err := t.Execute(w, data); err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
		}
	}
}

func adminDeleteHandler(cfg Config, db relayer.Storage) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if !strings.EqualFold(r.Method, http.MethodDelete) {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		if !auth(r, cfg.Admins) {
			w.Header().Add("WWW-Authenticate", `Basic realm="username and password required"`)
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("Unauthorized"))
			return
		}

		var (
			id = r.URL.Query().Get("id")
		)

		filter := nostr.Filter{}
		if id != "" {
			filter.IDs = []string{id}
		}

		if len(filter.IDs) == 0 {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("must provide id param"))
			return
		}

		ctx := context.Background()
		events, err := db.QueryEvents(ctx, &filter)
		if err != nil {
			fmt.Printf("[error] query events: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}

		for event := range events {
			if err := db.DeleteEvent(ctx, event.ID, event.PubKey); err != nil {
				log.Printf("[error] DeleteEvent %s: %s\n", event.ID, err)
			} else {
				log.Printf("deleted event: %s\n", event.String())
			}
		}

		w.WriteHeader(http.StatusNoContent)
		w.Write([]byte(""))
		return
	}
}

func auth(r *http.Request, admins map[string]string) bool {
	user, pass, ok := r.BasicAuth()
	if !ok {
		fmt.Println("Error parsing basic auth")
		return false
	}

	adminPass, ok := admins[user]
	if !ok {
		fmt.Printf("invalid user: %s\n", user)
		return false
	}

	if !strings.EqualFold(pass, adminPass) {
		fmt.Printf("invalid password for user: %s\n", user)
		return false
	}

	return true
}

type uiEvent struct {
	nostr.Event
	PrettyTime string
	Npub       string
}

func getEvents(db relayer.Storage, id, pk string, kind *int, limit int) ([]uiEvent, error) {
	filter := nostr.Filter{
		Limit: limit,
	}
	if id != "" {
		filter.IDs = []string{id}
	}
	if pk != "" {
		filter.Authors = []string{pk}
	}
	if kind != nil {
		filter.Kinds = []int{*kind}
	}

	ctx := context.Background()
	events, err := db.QueryEvents(ctx, &filter)
	if err != nil {
		return nil, fmt.Errorf("query events: %w", err)
	}

	var out []uiEvent
	for event := range events {
		npub, _ := nip19.EncodePublicKey(event.PubKey)

		out = append(out, uiEvent{
			Event:      *event,
			PrettyTime: event.CreatedAt.Time().Format(time.RFC822),
			Npub:       npub,
		})
	}

	return out, nil
}
