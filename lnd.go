package mock

import (
	"fmt"
	"strings"

	"github.com/buger/jsonparser"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
)

// CreateLndContainer create's an lnd container with a given name
// and no channels
func (c LightningMocker) CreateLndContainer(name string) (ctn LndContainer, err error) {
	ctn.c = &c
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
		return ctn, err
	}
	ctn.id = created.ID

	err = c.ContainerStart(c.Ctx, ctn.id, types.ContainerStartOptions{})
	if err != nil {
		return ctn, err
	}

	c.PrintContainerLogs(created.ID)

	return ctn, nil
}

type LndContainer struct {
	// id of the current docker container
	id string
	// the lightning mocker object
	c *LightningMocker
}

// get the hostname of the container
func (l LndContainer) Hostname() (hostname string, err error) {
	// get alices hostname
	hostnameResult, err := l.c.Exec(l.id, HostnameCmd)
	if err != nil {
		return "", err
	}
	return hostnameResult.StdOut, err
}

// Address gets the address of the user
func (l LndContainer) Address() (address string, err error) {
	hostname, err := l.Hostname()
	if err != nil {
		return "", err
	}
	// because we don't know when the lnd server will start, we need to keep trying until we get an address
	hasAddress := false
	counter := 0
	for !hasAddress {
		counter += 1
		if counter > 100 {
			return address, err
		}
		rawAddress, err := l.c.Exec(l.id, []string{"lncli", fmt.Sprintf("--rpcserver=%s:10009", strings.ReplaceAll(hostname, "\n", "")), "--network=simnet", "newaddress", "np2wkh"})
		if err != nil {
			return "", err
		}
		hasAddress = rawAddress.ExitCode == 0
		if hasAddress {
			address, _ = jsonparser.GetString([]byte(rawAddress.StdOut), "address")
		}
	}
	return address, err
}