// Package hyperloglog implements the HyperLogLog algorithm for
// cardinality estimation. In English: it counts things. It counts
// things using very small amounts of memory compared to the number of
// objects it is counting.
//
// For a full description of the algorithm, see the paper HyperLogLog:
// the analysis of a near-optimal cardinality estimation algorithm by
// Flajolet, et. al.
package hyperloglog

import (
	"fmt"
	"math"
)

var (
	exp32 = math.Pow(2, 32)
)

// A HyperLogLog is a deterministic cardinality estimator.  This version
// exports its fields so that it is suitable for saving eg. to a database.
type HyperLogLog struct {
	M         uint    // Number of registers
	B         uint32  // Number of bits used to determine register index
	Alpha     float64 // Bias correction constant
	Registers []uint8
}

// Compute bias correction alpha_m.
func getAlpha(m uint) (result float64) {
	switch m {
	case 16:
		result = 0.673
	case 32:
		result = 0.697
	case 64:
		result = 0.709
	default:
		result = 0.7213 / (1.0 + 1.079/float64(m))
	}
	return result
}

// New creates a HyperLogLog with the given number of registers. More
// registers leads to lower error in your estimated count, at the
// expense of memory.
//
// Choose a power of two number of registers, depending on the amount
// of memory you're willing to use and the error you're willing to
// tolerate. Each register uses one byte of memory.
//
// Standard error will be: σ ≈ 1.04 / sqrt(registers)
// The estimates provided by hyperloglog are expected to be within σ, 2σ, 3σ
// of the exact count in respectively 65%, 95%, 99% of all the cases.
func New(registers uint) (*HyperLogLog, error) {
	if (registers & (registers - 1)) != 0 {
		return nil, fmt.Errorf("number of registers %d not a power of two", registers)
	}
	h := &HyperLogLog{}
	h.M = registers
	h.B = uint32(math.Ceil(math.Log2(float64(registers))))
	h.Alpha = getAlpha(registers)
	h.Registers = make([]uint8, h.M)
	return h, nil
}

// Reset all internal variables and set the count to zero.
func (h *HyperLogLog) Reset() {
	for i := range h.Registers {
		h.Registers[i] = 0
	}
}

// Calculate the position of the leftmost 1-bit.
func rho(val uint32, max uint32) uint8 {
	r := uint32(1)
	for val&0x80000000 == 0 && r <= max {
		r++
		val <<= 1
	}
	return uint8(r)
}

// Add to the count. val should be a 32 bit unsigned integer from a
// good hash function.
func (h *HyperLogLog) Add(val uint32) {
	k := 32 - h.B
	r := rho(val<<h.B, k)
	j := val >> uint(k)
	if r > h.Registers[j] {
		h.Registers[j] = r
	}
}

// Count returns the estimated cardinality.
func (h *HyperLogLog) Count() uint64 {
	return h.count(true)
}

// CountWithoutLargeRangeCorrection returns the estimated cardinality, without applying
// the large range correction proposed by Flajolet et al. as it can lead to significant
// overcounting.
//
// See https://github.com/DataDog/hyperloglog/pull/15
func (h *HyperLogLog) CountWithoutLargeRangeCorrection() uint64 {
	return h.count(false)
}

func (h *HyperLogLog) count(withLargeRangeCorrection bool) uint64 {
	sum := 0.0
	m := float64(h.M)
	for _, val := range h.Registers {
		sum += 1.0 / math.Pow(2.0, float64(val))
	}
	estimate := h.Alpha * m * m / sum
	if estimate <= 5.0/2.0*m {
		// Small range correction
		v := 0
		for _, r := range h.Registers {
			if r == 0 {
				v++
			}
		}
		if v > 0 {
			estimate = m * math.Log(m/float64(v))
		}
	} else if estimate > 1.0/30.0*exp32 && withLargeRangeCorrection {
		// Large range correction
		estimate = -exp32 * math.Log(1-estimate/exp32)
	}
	return uint64(estimate)
}

// Merge another HyperLogLog into this one. The number of registers in
// each must be the same.
func (h *HyperLogLog) Merge(other *HyperLogLog) error {
	if h.M != other.M {
		return fmt.Errorf("number of registers doesn't match: %d != %d",
			h.M, other.M)
	}
	for j, r := range other.Registers {
		if r > h.Registers[j] {
			h.Registers[j] = r
		}
	}
	return nil
}
