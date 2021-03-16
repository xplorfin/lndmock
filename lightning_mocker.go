package mock

import (
	"github.com/hashicorp/go-multierror"
	"github.com/xplorfin/docker-utils"
)

// LightningMocker defines the lnd mocker object
type LightningMocker struct {
	docker.Client
	networkID string
}

// NewLightningMocker creates a new lightning mock object with a given
// session id
func NewLightningMocker() LightningMocker {
	return LightningMocker{
		Client:    docker.NewDockerClient(),
		networkID: "",
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
