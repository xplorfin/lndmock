package mock

import "github.com/docker/docker/api/types/container"

// NetworkName is the name of the network we use in docker
const NetworkName container.NetworkMode = "lightning-network"

// MiningAddressName is the name of the mining address argument passed to BtcdContainer
const MiningAddressName = "MINING_ADDRESS"

// NetworkCmd defines a constant for the network all command use
const NetworkCmd = "--network=simnet"

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

// HostnameCmd defines the command to fetch the hostname of a given docker container
// this is a function to make it immutable
func HostnameCmd() []string {
	return []string{"hostname", "-i"}
}
