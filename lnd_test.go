package mock

import (
	"testing"

	. "github.com/stretchr/testify/assert"
)

// run through this https://git.io/JqcC4 workflow
func TestLightningMocker(t *testing.T) {
	mocker := NewLightningMocker()
	defer func() {
		Nil(t, mocker.Teardown())
	}()
	err := mocker.Initialize()
	if err != nil {
		t.Error(err)
	}

	// start btcd as a prereq to lnd
	btcdContainer, err := mocker.CreateBtcdContainer()
	if err != nil {
		t.Error(err)
	}

	// start alice's lnd instance
	aliceContainer, err := mocker.CreateLndContainer("alice")
	if err != nil {
		t.Error(err)
	}

	// get alices hostname
	aliceAddress, err := aliceContainer.Address()
	if err != nil {
		t.Error(err)
	}

	alicePubKey, err := aliceContainer.GetPubKey()
	if err != nil {
		t.Error(err)
	}

	err = btcdContainer.MineToAddress(aliceAddress, 500)
	if err != nil {
		t.Error(err)
	}

	// start bob's lnd instance
	bobContainer, err := mocker.CreateLndContainer("bob")
	if err != nil {
		t.Error(err)
	}

	// give bob btc
	bobAddress, err := bobContainer.Address()
	if err != nil {
		t.Error(err)
	}

	err = btcdContainer.MineToAddress(bobAddress, 500)
	if err != nil {
		t.Error(err)
	}

	// remove until we can fix container link
	_ = alicePubKey
	// open alice->bob channel
	// error is currently cannot link toa  non-running container /btcd ad /bob/blockchain
	err = bobContainer.OpenChannel(alicePubKey, "alice", 100000)
	if err != nil {
		t.Error(err)
	}

	// get bob pub key
	bobPubKey, err := bobContainer.GetPubKey()
	if err != nil {
		t.Error(err)
	}

	// open bob->alice container
	err = aliceContainer.OpenChannel(bobPubKey, "bob", 100000)
	if err != nil {
		t.Error(err)
	}

	// broadcast channel opening transactions
	err = btcdContainer.MineToAddress(bobAddress, 3)
	if err != nil {
		t.Error(err)
	}
}
