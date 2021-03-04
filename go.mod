module github.com/xplorfin/lndmock

go 1.14

require (
	github.com/brianvoe/gofakeit/v5 v5.11.2
	github.com/btcsuite/btcd v0.21.0-beta.0.20201208033208-6bd4c64a54fa
	github.com/buger/jsonparser v1.1.1
	github.com/docker/docker v20.10.5+incompatible
	github.com/lightningnetwork/lnd v0.12.1-beta.rc6
	github.com/xplorfin/docker-utils v0.4.0
)

replace github.com/btcsuite/btcd v0.21.0-beta => github.com/xplorfin/btcd v0.21.0-hotfix
