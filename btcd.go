package mock

import (
	"fmt"
	"strconv"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
)

// Create a new btcd container
func (c LightningMocker) CreateBtcdContainer() (ctn BtcdContainer, err error) {
	ctn.c = &c
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
		return ctn, err
	}

	ctn.id = created.ID
	err = c.ContainerStart(c.Ctx, created.ID, types.ContainerStartOptions{})
	if err != nil {
		return ctn, err
	}

	c.PrintContainerLogs(created.ID)

	return ctn, nil
}

// btcd container object
type BtcdContainer struct {
	// id of the current docker container
	id string
	// the lightning mocker object
	c *LightningMocker
}

// mine a given number of block rewards to an address
func (b *BtcdContainer) MineToAddress(address string, blocks int) (err error) {
	b.id, err = b.recreateWithMiningAddress(b.id, address)
	if err != nil {
		return err
	}
	// generate n-blocks
	_, err = b.c.Exec(b.id, []string{"/start-btcctl.sh", "generate", strconv.Itoa(blocks)})

	return err
}

// recreate the btcd  container with a mining address (any subsequent blocks rewards will
// go to this address)
func (b *BtcdContainer) recreateWithMiningAddress(containerId string, miningAddress string) (id string, err error) {
	// remove the old container
	err = b.c.StopContainer(containerId)
	if err != nil {
		return containerId, err
	}
	err = b.c.ContainerRemove(b.c.Ctx, containerId, types.ContainerRemoveOptions{})
	if err != nil {
		return containerId, err
	}

	newEnvArgs := append(EnvArgs, fmt.Sprintf("%s=%s", MiningAddressName, miningAddress))
	created, err := b.c.CreateContainer(&container.Config{
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

	err = b.c.ContainerStart(b.c.Ctx, created.ID, types.ContainerStartOptions{})
	if err != nil {
		return id, err
	}

	b.c.PrintContainerLogs(created.ID)

	return created.ID, nil
}
