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
				Source: fmt.Sprintf("%s-lnd", name),
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
	// hostname of lnd container
	hostname string
	// address of lnd wallet
	address string
}

// get the hostname of the container
func (l *LndContainer) Hostname() (hostname string, err error) {
	if l.hostname == "" {
		// get alices hostname
		hostnameResult, err := l.c.Exec(l.id, HostnameCmd)
		if err != nil {
			return "", err
		}
		l.hostname = hostnameResult.StdOut
	}
	return l.hostname, err
}

// Address gets the address of the user
func (l *LndContainer) Address() (address string, err error) {
	if l.address == "" {
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
				return l.address, err
			}
			rawAddress, err := l.c.Exec(l.id, []string{"lncli", fmt.Sprintf("--rpcserver=%s:10009", strings.ReplaceAll(hostname, "\n", "")), NetworkCmd, "newaddress", "np2wkh"})
			if err != nil {
				return "", err
			}
			hasAddress = rawAddress.ExitCode == 0
			if hasAddress {
				l.address, _ = jsonparser.GetString([]byte(rawAddress.StdOut), "address")
			}
		}
	}
	return l.address, err
}

// GetPubKey of instance
func (l *LndContainer) GetPubKey() (pubKey string, err error) {
	// wait for start, TODO make this more efficient
	l.Address()

	hostname, err := l.Hostname()
	if err != nil {
		return "", err
	}

	info, err := l.c.Exec(l.id, []string{"lncli", fmt.Sprintf("--rpcserver=%s:10009", strings.ReplaceAll(hostname, "\n", "")), NetworkCmd, "getinfo"})
	if err != nil {
		return pubKey, err
	}
	pubKey, err = jsonparser.GetString([]byte(info.StdOut), "identity_pubkey")
	if err != nil {
		return pubKey, err
	}
	return pubKey, err
}

// broadcast a channel opening transaction
// note: blocks must be mined for channel to be established
func (l *LndContainer) OpenChannel(pubKey string, amount int) error {
	// wait for start, TODO make this more efficient
	l.Address()

	hostname, err := l.Hostname()
	if err != nil {
		return err
	}

	_, err = l.c.Exec(l.id, []string{"lncli", fmt.Sprintf("--rpcserver=%s:10009", strings.ReplaceAll(hostname, "\n", "")), NetworkCmd,
		"openchannel", fmt.Sprintf("--node_key=%s", pubKey), fmt.Sprintf("--local_amt=%d", amount)})
	return err
}