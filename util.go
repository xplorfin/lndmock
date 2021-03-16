package mock

import (
	"github.com/docker/docker/api/types/container"
)

// NetworkName is the name of the network we use in docker
const NetworkName container.NetworkMode = "lightning-network"

// MiningAddressName is the name of the mining address argument passed to BtcdContainer
const MiningAddressName = "MINING_ADDRESS"

// LnCLIPrefix defines a constant for the network all command use
// this is a function to make it immutable
func LnCLIPrefix() []string {
	return []string{"lncli", "--rpcserver=localhost:10009", "--network=simnet"}
}

// EnvArgs defines a list of arguments that must be used with the rpc server
// this is a function to make it immutable
func EnvArgs() []string {
	return []string{
		"RPCUSER",
		"RPCPASS",
		"NETWORK=simnet",
		"DEBUG",
	}
}

// Volumes defines mounted volumes
// this is a function to make it immutable
func Volumes() []string {
	return []string{"shared", "bitcoin", "lnd"}
}
