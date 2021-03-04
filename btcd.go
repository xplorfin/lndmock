package mock

import (
	"fmt"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
)

func (c LightningMocker) CreateBtcdContainer() (id string, err error) {
	newEnvArgs := append(EnvArgs, MiningAddressName)
	created, err := c.CreateContainer(&container.Config{
		Image:      "ghcr.io/xplorfin/btcd:latest",
		Env:        newEnvArgs,
		Tty:        false,
		Entrypoint: []string{"./start-btcd.sh"},
	}, &container.HostConfig{
		Mounts: []mount.Mount{
			{
				Source: "shared",
				Target: "/rpc",
				Type:   mount.TypeVolume,
			},
			{
				Source: "bitcoin",
				Target: "/data",
				Type:   mount.TypeVolume,
			},
		},
	}, nil, nil, "btcd")

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

// recreate the btcd  container with the correct mining address
// TODO reduce redundant code here
func (c LightningMocker) RecreateBtcdContainerMining(containerId string, miningAddress string) (id string, err error) {
	// remove the old container
	err = c.StopContainer(containerId)
	if err != nil {
		return containerId, err
	}
	err = c.ContainerRemove(c.Ctx, containerId, types.ContainerRemoveOptions{})
	if err != nil {
		return containerId, err
	}

	newEnvArgs := append(EnvArgs, fmt.Sprintf("%s=%s", MiningAddressName, miningAddress))
	created, err := c.CreateContainer(&container.Config{
		Image:      "ghcr.io/xplorfin/btcd:latest",
		Env:        newEnvArgs,
		Tty:        false,
		Entrypoint: []string{"./start-btcd.sh"},
	}, &container.HostConfig{
		Mounts: []mount.Mount{
			{
				Source: "shared",
				Target: "/rpc",
				Type:   mount.TypeVolume,
			},
			{
				Source: "bitcoin",
				Target: "/data",
				Type:   mount.TypeVolume,
			},
		},
	}, nil, nil, "btcd")

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
