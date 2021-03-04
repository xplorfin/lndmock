package mock

import (
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	"github.com/xplorfin/docker-utils"
)

type LightningMocker struct {
	docker.Client
}

func NewLightningMocker() LightningMocker {
	return LightningMocker{
		docker.NewDockerClient(),
	}
}

// idempotently create volumes
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

func (c LightningMocker) Teardown() error {
	return c.TeardownSession()
}

func (c LightningMocker) CreateLndContainer(name string) (id string, err error) {
	created, err := c.CreateContainer(&container.Config{
		Image:      "ghcr.io/xplorfin/lnd:latest",
		Env:        EnvArgs,
		Tty:        false,
		Entrypoint: []string{"./start-lnd.sh"},
		Labels: map[string]string{
			"sessionId": c.SessionId,
		},
	}, &container.HostConfig{
		Links: []string{"btcd:blockchain"},
		Mounts: []mount.Mount{
			{
				Source: "shared",
				Target: "/rpc",
				Type:   mount.TypeVolume,
			},
			{
				Source: "lnd",
				Target: "/root/.lnd",
				Type:   mount.TypeVolume,
			},
		},
		// TODO name fix
	}, nil, nil, name)

	if err != nil {
		return id, err
	}

	err = c.ContainerStart(c.Ctx, created.ID, types.ContainerStartOptions{})
	if err != nil {
		return id, err
	}

	c.PrintContainerLogs(created.ID)

	return created.ID, nil
}
