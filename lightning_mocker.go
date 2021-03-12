package mock

import "github.com/xplorfin/docker-utils"

// lnd mocker object
type LightningMocker struct {
	docker.Client
}

// NewLightningMocker creates a new lightning mock object with a given
// session id
func NewLightningMocker() LightningMocker {
	return LightningMocker{
		docker.NewDockerClient(),
	}
}

// CreateVolumes creates  new volumes if they don't exist
func (c LightningMocker) CreateVolumes() error {
	for _, volume := range Volumes {
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
