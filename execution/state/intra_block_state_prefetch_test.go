package state

import (
	"testing"

	"github.com/erigontech/erigon/common"
	"github.com/erigontech/erigon/execution/types/accounts"
	"github.com/holiman/uint256"
)

// mockReader returns minimal but non-nil account and code to validate prefetching
type mockReader struct{}

func (m *mockReader) ReadAccountData(address common.Address) (*accounts.Account, error) {
	a := &accounts.Account{}
	a.Balance = *uint256.NewInt(1)
	return a, nil
}
func (m *mockReader) ReadAccountDataForDebug(address common.Address) (*accounts.Account, error) {
	return m.ReadAccountData(address)
}
func (m *mockReader) ReadAccountStorage(address common.Address, key common.Hash) (uint256.Int, bool, error) {
	return uint256.Int{}, false, nil
}
func (m *mockReader) HasStorage(address common.Address) (bool, error)               { return false, nil }
func (m *mockReader) ReadAccountCode(address common.Address) ([]byte, error)        { return []byte{0x1}, nil }
func (m *mockReader) ReadAccountCodeSize(address common.Address) (int, error)       { return 1, nil }
func (m *mockReader) ReadAccountIncarnation(address common.Address) (uint64, error) { return 0, nil }
func (m *mockReader) SetTrace(trace bool, tracePrefix string)                       {}

func TestPreloadHotContracts(t *testing.T) {
	// Save original value to restore
	original := EnablePrefetch
	EnablePrefetch = true
	defer func() { EnablePrefetch = original }()

	mr := &mockReader{}
	ibs := New(mr)

	// Verify hot addresses are loaded
	for _, addr := range HotAddresses {
		if _, ok := ibs.stateObjects[addr]; !ok {
			t.Fatalf("expected preload to populate stateObjects with %s", addr.Hex())
		}
	}

	// Now disable and ensure not preloading
	EnablePrefetch = false
	ibs2 := New(mr)
	if len(ibs2.stateObjects) != 0 {
		t.Fatalf("expected no preloaded objects when EnablePrefetch=false, got %d", len(ibs2.stateObjects))
	}
}
