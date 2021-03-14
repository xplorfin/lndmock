package mock

import (
	"testing"

	"github.com/lightningnetwork/lnd/zpay32"
)

// Make sure we can create and decode an lnd invoice
func TestMockLndInvoice(t *testing.T) {
	mockEncodedInvoice, _ := MockLndInvoiceMainnet(t)
	_, err := zpay32.Decode(mockEncodedInvoice, &params)
	if err != nil {
		t.Error(err)
	}
}
