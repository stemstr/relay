package main

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/jmoiron/sqlx"
)

var subRequiredKinds = []int{1, 6, 16, 1808}

func kindRequiresSubscription(kind int) bool {
	for _, k := range subRequiredKinds {
		if kind == k {
			return true
		}
	}

	return false
}

func isSubscribed(db *sqlx.DB, pubkey string) (bool, error) {
	const query = `SELECT id 
FROM subscription
WHERE pubkey=$1
	AND expires_at > NOW()
ORDER BY created_at DESC
LIMIT 1;
`

	var id string
	if err := db.Get(&id, query, pubkey); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, nil
		}

		return false, fmt.Errorf("db.Get sub: %w", err)
	}

	return true, nil
}
