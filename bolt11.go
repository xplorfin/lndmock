package mock

import (
	"encoding/hex"
	"fmt"
	"testing"
	"time"

	"github.com/brianvoe/gofakeit/v5"
	"github.com/btcsuite/btcd/btcec"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/lightningnetwork/lnd/lnwire"
	"github.com/lightningnetwork/lnd/zpay32"
)

var params = chaincfg.MainNetParams

// create a mock lnd invoice on mainnet
func MockLndInvoiceMain(t *testing.T) (encoded string, decoded *zpay32.Invoice) {
	// amounts from https://git.io/JttyN
	var (
		testPaymentHash     [32]byte
		testDescriptionHash [32]byte
	)

	testPrivKeyBytes, _ := hex.DecodeString("e126f68f7eafcc8b74f54d269fe206be715000f94dac067d1c04a8ca3b2db734")
	testPrivKey, testPubKey := btcec.PrivKeyFromBytes(btcec.S256(), testPrivKeyBytes)

	testMessageSigner := zpay32.MessageSigner{
		SignCompact: func(hash []byte) ([]byte, error) {
			sig, err := btcec.SignCompact(btcec.S256(),
				testPrivKey, hash, true)
			if err != nil {
				return nil, fmt.Errorf("can't sign the "+
					"message: %v", err)
			}
			return sig, nil
		},
	}
	testDescription := gofakeit.Sentence(10)

	testAmount := lnwire.MilliSatoshi(gofakeit.RandomUint([]uint{2400000000000, 250000000, 2000000000}))

	testPaymentHashSlice, _ := hex.DecodeString("0001020304050607080900010203040506070809000102030405060708090102")
	testDescriptionHashSlice := chainhash.HashB([]byte(gofakeit.Sentence(20)))

	copy(testPaymentHash[:], testPaymentHashSlice[:])
	copy(testDescriptionHash[:], testDescriptionHashSlice[:])

	decoded = &zpay32.Invoice{
		Net:       &params,
		MilliSat:  &testAmount,
		Timestamp: gofakeit.DateRange(time.Now().AddDate(-1, 0, 0), time.Now()),
		//DescriptionHash: &testDescriptionHash,
		PaymentHash: &testPaymentHash,
		Description: &testDescription,
		Destination: testPubKey,
		// If no features were set, we'll populate an empty feature vector.
		Features: lnwire.NewFeatureVector(
			nil, lnwire.Features),
	}

	encoded, err := decoded.Encode(testMessageSigner)
	if err != nil {
		t.Error(err)
	}
	return encoded, decoded
}
