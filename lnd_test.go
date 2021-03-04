package mock

import (
	"fmt"
	"strings"
	"testing"

	"github.com/buger/jsonparser"
)

func TestLightningMocker(t *testing.T) {
	mocker := NewLightningMocker()
	defer mocker.Teardown()
	err := mocker.CreateVolumes()
	if err != nil {
		t.Error(err)
	}

	// start btcd as a prereq to lnd
	btcdId, err := mocker.CreateBtcdContainer()
	if err != nil {
		t.Error(err)
	}

	// start alice's lnd instance
	aliceId, err := mocker.CreateLndContainer("alice")
	if err != nil {
		t.Error(err)
	}

	// get alices hostname
	aliceHostname, err := mocker.Exec(aliceId, HostnameCmd)
	if err != nil {
		t.Error(err)
	}

	// because we don't know when the lnd server will start, we need to keep trying until we get an address
	hasAddress := false
	aliceAddress := ""
	counter := 0
	for !hasAddress {
		counter += 1
		if counter > 100 {
			t.Error("cannot start container")
			break
		}
		aliceAddressRaw, err := mocker.Exec(aliceId, []string{"lncli", fmt.Sprintf("--rpcserver=%s:10009", strings.ReplaceAll(aliceHostname.StdOut, "\n", "")), "--network=simnet", "newaddress", "np2wkh"})
		if err != nil {
			break
		}
		hasAddress = aliceAddressRaw.ExitCode == 0
		if hasAddress {
			aliceAddress, _ = jsonparser.GetString([]byte(aliceAddressRaw.StdOut), "address")
		}
	}
	// recreate the btcd container with alice's mining address
	btcdId, err = mocker.RecreateBtcdContainerMining(btcdId, aliceAddress)
	if err != nil {
		t.Error(err)
	}

	// give alice btc
	_, err = mocker.Exec(btcdId, []string{"/start-btcctl.sh", "generate", "400"})
	if err != nil {
		t.Error(err)
	}

	// start bob's lnd instance
	bobId, err := mocker.CreateLndContainer("bob")
	if err != nil {
		t.Error(err)
	}

	_ = bobId
}
