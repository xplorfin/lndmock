package mock

import (
	"fmt"
	"strconv"

	"github.com/docker/go-connections/nat"
)

// PortMap wraps nat.PortMap for easier port querying
type PortMap nat.PortMap

// NatMap returns a nat.Portap for use with docker
func (p PortMap) NatMap() nat.PortMap {
	return nat.PortMap(p)
}

// GetHostPort will return a host port for any given container port
func (p PortMap) GetHostPort(containerPort int) (hostPort int) {
	for containerPortData, hostPortData := range p {
		if containerPortData.Int() == containerPort {
			hostPort, err := strconv.Atoi(hostPortData[0].HostPort)
			if err != nil {
				panic(err)
			}
			return hostPort
		}
	}
	panic(fmt.Errorf("container port %d not found in PortMap", containerPort))
}
