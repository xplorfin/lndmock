package mock

const MiningAddressName = "MINING_ADDRESS"
const NetworkCmd = "--network=simnet"

var EnvArgs = []string{
	"RPCUSER",
	"RPCPASS",
	"NETWORK=simnet",
	"DEBUG",
}

var Volumes = []string{"shared", "bitcoin", "lnd"}

var HostnameCmd = []string{"hostname", "-i"}
