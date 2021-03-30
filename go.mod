module github.com/xplorfin/lndmock

go 1.14

require (
	github.com/brianvoe/gofakeit/v5 v5.11.2
	github.com/btcsuite/btcd v0.21.0-beta
	github.com/buger/jsonparser v1.1.1
	github.com/docker/docker v20.10.5+incompatible
	github.com/docker/go-connections v0.4.0
	github.com/hashicorp/go-multierror v1.1.1
	github.com/lightningnetwork/lnd v0.11.1-beta.rc5
	github.com/stretchr/testify v1.7.0
	github.com/xplorfin/docker-utils v0.13.0
	github.com/xplorfin/netutils v0.37.0
	google.golang.org/grpc v1.36.1
	gopkg.in/macaroon.v2 v2.1.0
)

replace (
	github.com/btcsuite/btcd v0.21.0-beta => github.com/xplorfin/btcd v0.21.0-hotfix
	google.golang.org/grpc => google.golang.org/grpc v1.29.1
)
