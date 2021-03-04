package mock

import (
	"testing"

	"github.com/lightningnetwork/lnd/zpay32"
)

func TestMockLndInvoice(t *testing.T) {
	mockEncodedInvoice, _ := MockLndInvoiceMain(t)
	_, err := zpay32.Decode(mockEncodedInvoice, &params)
	if err != nil {
		t.Error(err)
	}
}
