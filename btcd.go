package mock

import (
	"fmt"
	"strconv"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
)

// CreateBtcdContainer creates a btccontainer with a mining address address
func (c LightningMocker) createBtcdContainerWithAddress(address string) (ctn BtcdContainer, err error) {
	ctn.c = &c
	newEnvArgs := append(EnvArgs(), fmt.Sprintf("%s=%s", MiningAddressName, address))
	created, err := c.CreateContainer(&container.Config{
		Image:      "ghcr.io/xplorfin/btcd:latest",
		Env:        newEnvArgs,
		Tty:        false,
		Entrypoint: []string{"./start-btcd.sh"},
		Labels:     c.GetSessionLabels(),
	}, &container.HostConfig{
		NetworkMode: NetworkName,
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
	}, nil, nil, "blockchain")

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

// CreateBtcdContainer creates a BtcdContainer and starts it so it can
// respond to rpc requests. Mining must be done manually and a mining address
// should be set using BtcdContainer.MineToAddress.
func (c LightningMocker) CreateBtcdContainer() (ctn BtcdContainer, err error) {
	return c.createBtcdContainerWithAddress("")
}

// BtcdContainer object contains methods that allow us to interact with a created
// btcd container
type BtcdContainer struct {
	// id of the current docker container
	id string
	// c reference to the lightning mocker object
	c *LightningMocker
}

// MineToAddress mines a given number of block rewards to an address
func (b *BtcdContainer) MineToAddress(address string, blocks int) (err error) {
	b.id, err = b.recreateWithMiningAddress(b.id, address)
	if err != nil {
		return err
	}
	// generate n-blocks
	_, err = b.c.Exec(b.id, []string{"/start-btcctl.sh", "generate", strconv.Itoa(blocks)})

	return err
}

// recreateWithMiningAddress recreates the btcd  container with a mining address
// (any subsequent blocks rewards will go to this address)
func (b *BtcdContainer) recreateWithMiningAddress(containerID string, miningAddress string) (id string, err error) {
	// remove the old container
	err = b.c.StopContainer(containerID)
	if err != nil {
		return containerID, err
	}
	err = b.c.ContainerRemove(b.c.Ctx, containerID, types.ContainerRemoveOptions{})
	if err != nil {
		return containerID, err
	}

	ctn, err := b.c.createBtcdContainerWithAddress(miningAddress)
	if err != nil {
		return containerID, err
	}

	return ctn.id, nil
}
