package mock

import (
	"fmt"
	"strconv"

	"github.com/docker/go-connections/nat"
	"github.com/hashicorp/go-multierror"
	"github.com/xplorfin/docker-utils"
	"github.com/xplorfin/netutils"
)

// LightningMocker defines the lnd mocker object
type LightningMocker struct {
	docker.Client
	networkID string
	portStack netutils.FreePortStack
}

// NewLightningMocker creates a new lightning mock object with a given
// session id
func NewLightningMocker() LightningMocker {
	return LightningMocker{
		Client:    docker.NewDockerClient(),
		networkID: "",
		portStack: netutils.NewFreeportStack(),
	}
}

// Initialize creates common resources by calling CreateNetworks and CreateVolumes
// and returning an error if necessary
func (c *LightningMocker) Initialize() (err error) {
	cve := c.CreateNetworks()
	if cve != nil {
		err = multierror.Append(err, cve)
	}
	cne := c.CreateVolumes()
	if cne != nil {
		err = multierror.Append(err, cne)
	}
	return err
}

// CreateNetworks sets up the networks needed for
func (c *LightningMocker) CreateNetworks() (err error) {
	c.networkID, err = c.CreateNetwork(string(NetworkName))
	return err
}

// CreateVolumes creates  new volumes if they don't exist
func (c LightningMocker) CreateVolumes() error {
	for _, volume := range Volumes() {
		if !c.VolumeExists(volume) {
			err := c.CreateVolume(volume)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// Teardown removes all containers created in the session
func (c LightningMocker) Teardown() error {
	return c.TeardownSession()
}

// portsToMap takes an arbitrary list of ports and converts them to a port map
// with the local machine (using free, random ports)
func (c *LightningMocker) portsToMap(ports []int) (pm nat.PortMap) {
	pm = make(nat.PortMap)
	for _, port := range ports {
		pm[nat.Port(fmt.Sprintf("%d/tcp", port))] = []nat.PortBinding{
			{HostIP: "0.0.0.0", HostPort: strconv.Itoa(c.portStack.GetPort())},
		}
	}
	return pm
}
