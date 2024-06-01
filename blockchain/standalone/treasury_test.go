// Copyright (c) 2020 The Decred developers
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package standalone

import (
	"errors"
	"testing"
)

const (
	// mainNetTVI is the treasury vote interval for mainnet.
	mainNetTVI = 288

	// mainNetTVIMul is the treasury vote interval multiplier for mainnet.
	mainNetTVIMul = 12
)

// TestIsTreasuryVoteInterval ensures that the function to determine if a given
// block height is on a treasury vote interval returns the expected results
// including for edge conditions.
func TestIsTreasuryVoteInterval(t *testing.T) {
	tests := []struct {
		name   string // test description
		height uint64 // height to check
		tvi    uint64 // treasury vote interval
		want   bool   // expected result
	}{{
		name:   "0 is never considered a TVI",
		height: 0,
		tvi:    mainNetTVI,
		want:   false,
	}, {
		name:   "TVI - 1",
		height: mainNetTVI - 1,
		tvi:    mainNetTVI,
		want:   false,
	}, {
		name:   "exactly TVI",
		height: mainNetTVI,
		tvi:    mainNetTVI,
		want:   true,
	}, {
		name:   "TVI + 1",
		height: mainNetTVI + 1,
		tvi:    mainNetTVI,
		want:   false,
	}, {
		name:   "Multiple of TVI",
		height: 2 * mainNetTVI,
		tvi:    mainNetTVI,
		want:   true,
	}}

	for _, test := range tests {
		result := IsTreasuryVoteInterval(test.height, test.tvi)
		if result != test.want {
			t.Errorf("%q: Unexpected result -- got %v, want %v", test.name,
				result, test.want)
		}
	}
}

// TestCalcTSpendWindow ensures that the function that calculates the start and
// end of a treasury spend voting window returns the expected results including
// error conditions.
func TestCalcTSpendWindow(t *testing.T) {
	tests := []struct {
		name      string // test description
		expiry    uint32 // expiry to calculate the window for
		tvi       uint64 // treasury vote interval
		tvimul    uint64 // treasury vote interval multiplier
		err       error  // expected error
		wantStart uint32 // expected start result
		wantEnd   uint32 // expected end result
	}{{
		name:      "zero is not a valid expiry",
		expiry:    0,
		tvi:       mainNetTVI,
		tvimul:    mainNetTVIMul,
		err:       ErrInvalidTSpendExpiry,
		wantStart: 0,
		wantEnd:   0,
	}, {
		name:      "min required expiry - 1",
		expiry:    mainNetTVI*mainNetTVIMul + 1,
		tvi:       mainNetTVI,
		tvimul:    mainNetTVIMul,
		err:       ErrInvalidTSpendExpiry,
		wantStart: 0,
		wantEnd:   0,
	}, {
		name:      "not a TVI + 2",
		expiry:    mainNetTVI*mainNetTVIMul + 3,
		tvi:       mainNetTVI,
		tvimul:    mainNetTVIMul,
		err:       ErrInvalidTSpendExpiry,
		wantStart: 0,
		wantEnd:   0,
	}, {
		name:      "5 is not a valid start or end for a tvi 11, mul 3",
		expiry:    5,
		tvi:       11,
		tvimul:    3,
		err:       ErrInvalidTSpendExpiry,
		wantStart: 0,
		wantEnd:   0,
	}, {
		name:      "first possible valid mainnet params",
		expiry:    mainNetTVI*mainNetTVIMul + 2,
		tvi:       mainNetTVI,
		tvimul:    mainNetTVIMul,
		err:       nil,
		wantStart: 0,
		wantEnd:   mainNetTVI * mainNetTVIMul,
	}, {
		name:      "second possible valid mainnet params",
		expiry:    mainNetTVI*mainNetTVIMul*2 + 2,
		tvi:       mainNetTVI,
		tvimul:    mainNetTVIMul,
		err:       nil,
		wantStart: mainNetTVI * mainNetTVIMul,
		wantEnd:   mainNetTVI * mainNetTVIMul * 2,
	}, {
		name:      "5186 for tvi 288, mul 7 is window [3168, 5184)",
		expiry:    5186,
		tvi:       288,
		tvimul:    7,
		err:       nil,
		wantStart: 5186 - 288*7 - 2,
		wantEnd:   5186 - 2,
	}}

	for _, test := range tests {
		// Calculate result and ensure the expected error is produced.
		start, end, err := CalcTSpendWindow(test.expiry, test.tvi, test.tvimul)
		if !errors.Is(err, test.err) {
			t.Errorf("%q: unexpected error -- got %v, want %v", test.name, err,
				test.err)
			continue
		}

		// Ensure the expected start value is calculated.
		if start != test.wantStart {
			t.Errorf("%q: unexpected start val -- got %v, want %v", test.name,
				start, test.wantStart)
			continue
		}

		// Ensure the expected end value is calculated.
		if end != test.wantEnd {
			t.Errorf("%q: unexpected end val -- got %v, want %v", test.name,
				end, test.wantEnd)
			continue
		}
	}
}

// TestCalcTSpendExpiry ensures that the function to calculate the expiry for
// the treasury spend voting window for a given height produces the expected
// results.
func TestCalcTSpendExpiry(t *testing.T) {
	tests := []struct {
		name       string // test description
		height     int64  // height to calculate the expiry for
		tvi        uint64 // treasury vote interval
		tvimul     uint64 // treasury vote interval multiplier
		wantExpiry uint32 // expected expiry
	}{{
		name:       "mul 1, tvi 288, first block in first tvi",
		height:     0,
		tvi:        288,
		tvimul:     1,
		wantExpiry: 578,
	}, {
		name:       "mul 1, tvi 288, last block in first tvi",
		height:     287,
		tvi:        288,
		tvimul:     1,
		wantExpiry: 578,
	}, {
		name:       "mul 1, tvi 288, first block in second tvi",
		height:     288,
		tvi:        288,
		tvimul:     1,
		wantExpiry: 866,
	}, {
		name:       "mul 2, tvi 288, first block in first tvi",
		height:     0,
		tvi:        288,
		tvimul:     2,
		wantExpiry: 866,
	}, {
		name:       "mul 2, tvi 288, last block in first tvi",
		height:     287,
		tvi:        288,
		tvimul:     2,
		wantExpiry: 866,
	}, {
		name:       "mul 2, tvi 288, first block in second tvi",
		height:     288,
		tvi:        288,
		tvimul:     2,
		wantExpiry: 1154,
	}, {
		name:       "mul 60, tvi 4, block in middle of 13th tvi",
		height:     810,
		tvi:        60,
		tvimul:     4,
		wantExpiry: 1082,
	}, {
		name:       "mul 7, tvi 288, first block in 10th tvi",
		height:     2880,
		tvi:        288,
		tvimul:     7,
		wantExpiry: 5186,
	}}

	for _, test := range tests {
		// Calculate expiry and ensure the expected value is produced.
		// Ensure the expected start value is calculated.
		expiry := CalcTSpendExpiry(test.height, test.tvi, test.tvimul)
		if expiry != test.wantExpiry {
			t.Errorf("%q: unexpected expiry -- got %v, want %v", test.name,
				expiry, test.wantExpiry)
			continue
		}
	}
}

// TestInsideTSpendWindow ensures that the function to determine if a given
// height is within a treasury spend voting window returns the expected results.
func TestInsideTSpendWindow(t *testing.T) {
	tests := []struct {
		name   string // test description
		height int64  // height to check
		expiry uint32 // expiry for the voting window
		tvi    uint64 // treasury vote interval
		tvimul uint64 // treasury vote interval multiplier
		want   bool   // expected result
	}{{
		name:   "invalid expiry but otherwise correct",
		height: 3167,
		expiry: 5185,
		tvi:    288,
		tvimul: 7,
		want:   false,
	}, {
		name:   "one block before window start",
		height: 3167,
		expiry: 5186,
		tvi:    288,
		tvimul: 7,
		want:   false,
	}, {
		name:   "exactly window start",
		height: 3168,
		expiry: 5186,
		tvi:    288,
		tvimul: 7,
		want:   true,
	}, {
		name:   "last block of window",
		height: 5184,
		expiry: 5186,
		tvi:    288,
		tvimul: 7,
		want:   true,
	}, {
		name:   "one block after window",
		height: 5185,
		expiry: 5186,
		tvi:    288,
		tvimul: 7,
		want:   false,
	}}

	for _, test := range tests {
		result := InsideTSpendWindow(test.height, test.expiry, test.tvi,
			test.tvimul)
		if result != test.want {
			t.Errorf("%q: unexpected result -- got %v, want %v", test.name,
				result, test.want)
			continue
		}
	}
}
