package mock

import (
	"testing"

	"github.com/lightningnetwork/lnd/lnrpc"

	. "github.com/stretchr/testify/assert"
)

// run through this https://git.io/JqcC4 workflow
func TestLightningMocker(t *testing.T) {
	mocker := NewLightningMocker()
	defer func() {
		Nil(t, mocker.Teardown())
	}()
	err := mocker.Initialize()
	Nil(t, err)

	// start btcd as a prereq to lnd
	btcdContainer, err := mocker.CreateBtcdContainer()
	Nil(t, err)

	// start alice's lnd instance
	aliceContainer, err := mocker.CreateLndContainer("alice")
	Nil(t, err)

	// get alices hostname
	aliceAddress, err := aliceContainer.Address()
	Nil(t, err)

	alicePubKey, err := aliceContainer.GetPubKey()
	Nil(t, err)

	err = btcdContainer.MineToAddress(aliceAddress, 500)
	Nil(t, err)

	// start bob's lnd instance
	bobContainer, err := mocker.CreateLndContainer("bob")
	Nil(t, err)

	// give bob btc
	bobAddress, err := bobContainer.Address()
	Nil(t, err)

	err = btcdContainer.MineToAddress(bobAddress, 500)
	Nil(t, err)

	// open alice->bob channel
	err = bobContainer.OpenChannel(alicePubKey, "alice", 100000)
	Nil(t, err)

	// get bob pub key
	bobPubKey, err := bobContainer.GetPubKey()
	Nil(t, err)

	// open bob->alice container
	err = aliceContainer.OpenChannel(bobPubKey, "bob", 100000)
	Nil(t, err)

	// broadcast channel opening transactions
	err = btcdContainer.MineToAddress(bobAddress, 3)
	Nil(t, err)

	testRPCClient(t, bobContainer)
	testRPCClient(t, aliceContainer)
}

func testRPCClient(t *testing.T, c LndContainer) {
	client, err := c.RPCClient()
	Nil(t, err)

	req := lnrpc.GetInfoRequest{}
	res, err := client.GetInfo(c.c.Ctx, &req)
	_ = res
	Nil(t, err)
}
