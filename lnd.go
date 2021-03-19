package mock

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"strings"

	"github.com/lightningnetwork/lnd/macaroons"

	"github.com/buger/jsonparser"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	"github.com/lightningnetwork/lnd/lnrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"gopkg.in/macaroon.v2"
)

// CreateLndContainer create's an lnd container with a given name
// and no channels
func (c LightningMocker) CreateLndContainer(name string) (ctn LndContainer, err error) {
	ctn.c = &c
	ctn.PortMap = c.portsToMap([]int{10009, 8080, 9735})
	created, err := c.CreateContainer(&container.Config{
		Image:      "ghcr.io/xplorfin/lnd:latest",
		Env:        EnvArgs(),
		Tty:        false,
		Entrypoint: []string{"./start-lnd.sh"},
		Labels:     c.GetSessionLabels(),
	}, &container.HostConfig{
		PortBindings: ctn.PortMap.NatMap(),
		NetworkMode:  NetworkName,
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

// LndContainer object contains methods that allow us to interact with a created
// lnd container
type LndContainer struct {
	// id of the current docker container
	id string
	// the lightning mocker object
	c *LightningMocker
	// address of lnd wallet
	address string
	// PortMap is the mapping of ports to the host binding
	PortMap PortMap
}

// Address gets the address of the user
func (l *LndContainer) Address() (address string, err error) {
	if l.address == "" {
		// because we don't know when the lnd server will start, we need to keep trying until we get an address
		hasAddress := false
		counter := 0
		for !hasAddress {
			counter++
			if counter > 100 {
				return l.address, err
			}
			// TODO we might be able to replace the hostname here with the container command
			rawAddress, err := l.c.Exec(l.id, append(LnCLIPrefix(), "newaddress", "np2wkh"))
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
	info, err := l.c.Exec(l.id, append(LnCLIPrefix(), "getinfo"))
	if err != nil {
		return pubKey, err
	}
	pubKey, err = jsonparser.GetString([]byte(info.StdOut), "identity_pubkey")
	if err != nil {
		return pubKey, err
	}
	return pubKey, err
}

// Connect connects to the lightning peer
func (l *LndContainer) Connect(pubKey, host string) (err error) {
	// we use address to make sure the wallet is unlocked
	// TODO: clean this up
	_, err = l.Address()
	if err != nil {
		return err
	}

	_, err = l.c.Exec(l.id, append(LnCLIPrefix(), "connect", fmt.Sprintf("%s@%s", pubKey, host)))
	return err
}

// OpenChannel connects to the peer and broadcasts a channel
// opening transaction to the mempool.
// Note: blocks must be mined for channel to be established
func (l *LndContainer) OpenChannel(pubKey, host string, amount int) (err error) {
	err = l.Connect(pubKey, host)
	if err != nil {
		return err
	}
	// open the channel
	_, err = l.c.Exec(l.id, append(LnCLIPrefix(),
		"openchannel", fmt.Sprintf("--node_key=%s", pubKey), fmt.Sprintf("--local_amt=%d", amount)))
	return err
}

// GetTLSCert gets the tls cert for LndContainer
func (l *LndContainer) GetTLSCert() (cert *tls.Config, rawCert string, err error) {
	rawResult, err := l.c.Exec(l.id, []string{"cat", "/root/.lnd/tls.cert"})
	if err != nil {
		return cert, rawCert, err
	}

	rawCert = rawResult.StdOut

	cp := x509.NewCertPool()
	if !cp.AppendCertsFromPEM([]byte(rawCert)) {
		return cert, rawCert, err
	}

	cert = &tls.Config{
		InsecureSkipVerify: true,
		RootCAs:            cp,
	}

	return cert, rawCert, err
}

// GetAdminMacaroon fetches the admin macaroon from LndContainer
func (l *LndContainer) GetAdminMacaroon() (mac *macaroon.Macaroon, err error) {
	rawMacaroonRes, err := l.c.Exec(l.id, []string{"base64", "/root/.lnd/data/chain/bitcoin/simnet/admin.macaroon"})
	if err != nil {
		return mac, err
	}
	rawMac := strings.ReplaceAll(rawMacaroonRes.StdOut, "\n", "")
	decoded, err := macaroon.Base64Decode([]byte(rawMac))
	if err != nil {
		return nil, err
	}

	mac = &macaroon.Macaroon{}
	err = mac.UnmarshalBinary(decoded)
	if err != nil {
		return nil, err
	}

	return mac, err
}

// GrpcConnection generates a grpc connection to a
func (l *LndContainer) GrpcConnection() (conn *grpc.ClientConn, err error) {
	cert, _, err := l.GetTLSCert()
	if err != nil {
		return nil, err
	}
	mac, err := l.GetAdminMacaroon()
	if err != nil {
		return nil, err
	}
	return grpc.DialContext(
		l.c.Ctx,
		fmt.Sprintf("localhost:%d",
			l.PortMap.GetHostPort(10009)),
		grpc.WithTransportCredentials(credentials.NewTLS(cert)),
		grpc.WithPerRPCCredentials(macaroons.NewMacaroonCredential(mac)),
	)
}

// RPCClient gets an authenticated
func (l *LndContainer) RPCClient() (rpcClient lnrpc.LightningClient, err error) {
	grpcConn, err := l.GrpcConnection()
	if err != nil {
		return nil, err
	}
	return lnrpc.NewLightningClient(grpcConn), nil
}
