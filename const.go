package mock

// MiningAddressName is the name of the mining address argument passed to BtcdContainer
const MiningAddressName = "MINING_ADDRESS"

// NetworkCmd defines a constant for the network all command use
const NetworkCmd = "--network=simnet"

// EnvArgs defines a list of arguments that must be used with the rpc server
var EnvArgs = []string{
	"RPCUSER",
	"RPCPASS",
	"NETWORK=simnet",
	"DEBUG",
}

// Volumes defines mounted volumes
var Volumes = []string{"shared", "bitcoin", "lnd"}

// HostnameCmd defines the command to fetch the hostname of a given docker container
var HostnameCmd = []string{"hostname", "-i"}
