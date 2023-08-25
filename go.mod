module github.com/stemstr/relay

go 1.20

require (
	github.com/fiatjaf/relayer/v2 v2.1.0
	github.com/jmoiron/sqlx v1.3.1
	github.com/nbd-wtf/go-nostr v0.20.0
	github.com/stemstr/blastr v0.1.0
	github.com/stretchr/testify v1.8.4
	golang.org/x/time v0.0.0-20210220033141-f8bda1e9f3ba
	gopkg.in/yaml.v3 v3.0.1
)

// We run off our fork
replace github.com/fiatjaf/relayer/v2 => github.com/stemstr/relayer/v2 v2.1.4

require (
	github.com/andybalholm/brotli v1.0.5 // indirect
	github.com/bits-and-blooms/bitset v1.8.0 // indirect
	github.com/bits-and-blooms/bloom/v3 v3.5.0 // indirect
	github.com/btcsuite/btcd/btcec/v2 v2.3.2 // indirect
	github.com/btcsuite/btcd/btcutil v1.1.3 // indirect
	github.com/btcsuite/btcd/chaincfg/chainhash v1.0.2 // indirect
	github.com/cespare/xxhash/v2 v2.2.0 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/decred/dcrd/crypto/blake256 v1.0.1 // indirect
	github.com/decred/dcrd/dcrec/secp256k1/v4 v4.2.0 // indirect
	github.com/dgraph-io/ristretto v0.1.1 // indirect
	github.com/dustin/go-humanize v1.0.1 // indirect
	github.com/fasthttp/websocket v1.5.3 // indirect
	github.com/gobwas/httphead v0.1.0 // indirect
	github.com/gobwas/pool v0.2.1 // indirect
	github.com/gobwas/ws v1.2.1 // indirect
	github.com/golang/glog v1.1.1 // indirect
	github.com/josharian/intern v1.0.0 // indirect
	github.com/klauspost/compress v1.16.5 // indirect
	github.com/kr/pretty v0.3.1 // indirect
	github.com/lib/pq v1.10.3 // indirect
	github.com/mailru/easyjson v0.7.7 // indirect
	github.com/mattn/go-sqlite3 v1.14.16 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/puzpuzpuz/xsync v1.5.2 // indirect
	github.com/rogpeppe/go-internal v1.11.0 // indirect
	github.com/rs/cors v1.7.0 // indirect
	github.com/savsgio/gotils v0.0.0-20230208104028-c358bd845dee // indirect
	github.com/spaolacci/murmur3 v1.1.0 // indirect
	github.com/tidwall/gjson v1.15.0 // indirect
	github.com/tidwall/match v1.1.1 // indirect
	github.com/tidwall/pretty v1.2.1 // indirect
	github.com/valyala/bytebufferpool v1.0.0 // indirect
	github.com/valyala/fasthttp v1.47.0 // indirect
	github.com/willf/bitset v1.1.11 // indirect
	golang.org/x/exp v0.0.0-20230801115018-d63ba01acd4b // indirect
	golang.org/x/sys v0.11.0 // indirect
	gopkg.in/check.v1 v1.0.0-20201130134442-10cb98267c6c // indirect
)
