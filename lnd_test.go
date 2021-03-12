package mock

import (
	"fmt"
	"testing"
)

// run through this https://git.io/JqcC4 workflow
func TestLightningMocker(t *testing.T) {
	mocker := NewLightningMocker()
	defer mocker.Teardown()
	err := mocker.CreateVolumes()
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

	bobPubKey, err := bobContainer.GetPubKey()
	if err != nil {
		t.Error(err)
	}

	_ = alicePubKey
	_ = bobPubKey

	fmt.Print("")
	_ = bobContainer
}
