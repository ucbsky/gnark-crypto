// Copyright 2020 Consensys Software Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Code generated by consensys/gnark-crypto DO NOT EDIT

package fft

import (
	"math/bits"

	"github.com/consensys/gnark-crypto/ecc/bls12-381/fr"
)

// BitReverse applies the bit-reversal permutation to a.
// len(a) must be a power of 2 (as in every single function in this file)
func BitReverse(a []fr.Element) {
	n := uint64(len(a))
	nn := uint64(64 - bits.TrailingZeros64(n))

	for i := uint64(0); i < n; i++ {
		irev := bits.Reverse64(i) >> nn
		if irev > i {
			a[i], a[irev] = a[irev], a[i]
		}
	}
}

func deriveLogTileSize(logN uint64) uint64 {
	q := uint64(9)

	for int(logN)-int(2*q) <= 0 {
		q--
	}

	return q
}

func bitReverseCobraInPlace(buf []fr.Element) {
	logN := uint64(bits.Len64(uint64(len(buf))) - 1)
	logTileSize := deriveLogTileSize(logN)
	logBLen := logN - 2*logTileSize
	bLen := uint64(1) << logBLen
	bShift := logBLen + logTileSize
	tileSize := uint64(1) << logTileSize

	t := make([]fr.Element, tileSize*tileSize)

	for b := uint64(0); b < bLen; b++ {

		for a := uint64(0); a < tileSize; a++ {
			aRev := bits.Reverse64(a) >> (64 - logTileSize)
			for c := uint64(0); c < tileSize; c++ {
				tIdx := (aRev << logTileSize) | c
				idx := (a << (bShift)) | (b << logTileSize) | c
				t[tIdx] = buf[idx]
			}
		}

		bRev := bits.Reverse64(b) >> (64 - logBLen)

		for c := uint64(0); c < tileSize; c++ {
			cRev := bits.Reverse64(c) >> (64 - logTileSize)
			for aRev := uint64(0); aRev < tileSize; aRev++ {
				a := bits.Reverse64(aRev) >> (64 - logTileSize)
				idx := (a << (bShift)) | (b << logTileSize) | c
				idxRev := (cRev << (bShift)) | (bRev << logTileSize) | aRev
				if idx < idxRev {
					tIdx := (aRev << logTileSize) | c
					buf[idxRev], t[tIdx] = t[tIdx], buf[idxRev]
				}
			}
		}

		for a := uint64(0); a < tileSize; a++ {
			aRev := bits.Reverse64(a) >> (64 - logTileSize)
			for c := uint64(0); c < tileSize; c++ {
				cRev := bits.Reverse64(c) >> (64 - logTileSize)
				idx := (a << (bShift)) | (b << logTileSize) | c
				idxRev := (cRev << (bShift)) | (bRev << logTileSize) | aRev
				if idx < idxRev {
					tIdx := (aRev << logTileSize) | c
					buf[idx], t[tIdx] = t[tIdx], buf[idx]
				}
			}
		}
	}
}

func bitReverseCobraInPlace_9_21(buf []fr.Element) {
	const (
		logTileSize = uint64(9)
		tileSize    = uint64(1) << logTileSize
		logN        = 21
		logBLen     = logN - 2*logTileSize
		bShift      = logBLen + logTileSize
		bLen        = uint64(1) << logBLen
	)

	var t [tileSize * tileSize]fr.Element

	for b := uint64(0); b < bLen; b++ {
		bRev := (bits.Reverse64(b) >> (64 - logBLen)) << logTileSize

		for a := uint64(0); a < tileSize; a++ {
			aRev := (bits.Reverse64(a) >> 55) << logTileSize
			for c := uint64(0); c < tileSize; c++ {
				idx := (a << (bShift)) | (b << logTileSize) | c
				t[aRev|c] = buf[idx]
			}
		}

		for c := uint64(0); c < tileSize; c++ {
			cRev := (bits.Reverse64(c) >> 55 << bShift) | bRev
			for aRev := uint64(0); aRev < tileSize; aRev++ {
				a := bits.Reverse64(aRev) >> 55
				idx := (a << bShift) | (b << logTileSize) | c
				idxRev := cRev | aRev
				if idx < idxRev {
					tIdx := (aRev << logTileSize) | c
					buf[idxRev], t[tIdx] = t[tIdx], buf[idxRev]
				}
			}
		}

		for a := uint64(0); a < tileSize; a++ {
			aRev := bits.Reverse64(a) >> 55
			for c := uint64(0); c < tileSize; c++ {
				cRev := (bits.Reverse64(c) >> 55) << bShift
				idx := (a << bShift) | (b << logTileSize) | c
				idxRev := cRev | bRev | aRev
				if idx < idxRev {
					tIdx := (aRev << logTileSize) | c
					buf[idx], t[tIdx] = t[tIdx], buf[idx]
				}
			}
		}
	}

}

func bitReverseCobraInPlace_9_22(buf []fr.Element) {
	const (
		logTileSize = uint64(9)
		tileSize    = uint64(1) << logTileSize
		logN        = 22
		logBLen     = logN - 2*logTileSize
		bShift      = logBLen + logTileSize
		bLen        = uint64(1) << logBLen
	)

	var t [tileSize * tileSize]fr.Element

	for b := uint64(0); b < bLen; b++ {
		bRev := (bits.Reverse64(b) >> (64 - logBLen)) << logTileSize

		for a := uint64(0); a < tileSize; a++ {
			aRev := (bits.Reverse64(a) >> 55) << logTileSize
			for c := uint64(0); c < tileSize; c++ {
				idx := (a << (bShift)) | (b << logTileSize) | c
				t[aRev|c] = buf[idx]
			}
		}

		for c := uint64(0); c < tileSize; c++ {
			cRev := (bits.Reverse64(c) >> 55 << bShift) | bRev
			for aRev := uint64(0); aRev < tileSize; aRev++ {
				a := bits.Reverse64(aRev) >> 55
				idx := (a << bShift) | (b << logTileSize) | c
				idxRev := cRev | aRev
				if idx < idxRev {
					tIdx := (aRev << logTileSize) | c
					buf[idxRev], t[tIdx] = t[tIdx], buf[idxRev]
				}
			}
		}

		for a := uint64(0); a < tileSize; a++ {
			aRev := bits.Reverse64(a) >> 55
			for c := uint64(0); c < tileSize; c++ {
				cRev := (bits.Reverse64(c) >> 55) << bShift
				idx := (a << bShift) | (b << logTileSize) | c
				idxRev := cRev | bRev | aRev
				if idx < idxRev {
					tIdx := (aRev << logTileSize) | c
					buf[idx], t[tIdx] = t[tIdx], buf[idx]
				}
			}
		}
	}

}

func bitReverseCobraInPlace_9_23(buf []fr.Element) {
	const (
		logTileSize = uint64(9)
		tileSize    = uint64(1) << logTileSize
		logN        = 23
		logBLen     = logN - 2*logTileSize
		bShift      = logBLen + logTileSize
		bLen        = uint64(1) << logBLen
	)

	var t [tileSize * tileSize]fr.Element

	for b := uint64(0); b < bLen; b++ {
		bRev := (bits.Reverse64(b) >> (64 - logBLen)) << logTileSize

		for a := uint64(0); a < tileSize; a++ {
			aRev := (bits.Reverse64(a) >> 55) << logTileSize
			for c := uint64(0); c < tileSize; c++ {
				idx := (a << (bShift)) | (b << logTileSize) | c
				t[aRev|c] = buf[idx]
			}
		}

		for c := uint64(0); c < tileSize; c++ {
			cRev := (bits.Reverse64(c) >> 55 << bShift) | bRev
			for aRev := uint64(0); aRev < tileSize; aRev++ {
				a := bits.Reverse64(aRev) >> 55
				idx := (a << bShift) | (b << logTileSize) | c
				idxRev := cRev | aRev
				if idx < idxRev {
					tIdx := (aRev << logTileSize) | c
					buf[idxRev], t[tIdx] = t[tIdx], buf[idxRev]
				}
			}
		}

		for a := uint64(0); a < tileSize; a++ {
			aRev := bits.Reverse64(a) >> 55
			for c := uint64(0); c < tileSize; c++ {
				cRev := (bits.Reverse64(c) >> 55) << bShift
				idx := (a << bShift) | (b << logTileSize) | c
				idxRev := cRev | bRev | aRev
				if idx < idxRev {
					tIdx := (aRev << logTileSize) | c
					buf[idx], t[tIdx] = t[tIdx], buf[idx]
				}
			}
		}
	}

}

func bitReverseCobraInPlace_9_24(buf []fr.Element) {
	const (
		logTileSize = uint64(9)
		tileSize    = uint64(1) << logTileSize
		logN        = 24
		logBLen     = logN - 2*logTileSize
		bShift      = logBLen + logTileSize
		bLen        = uint64(1) << logBLen
	)

	var t [tileSize * tileSize]fr.Element

	for b := uint64(0); b < bLen; b++ {
		bRev := (bits.Reverse64(b) >> (64 - logBLen)) << logTileSize

		for a := uint64(0); a < tileSize; a++ {
			aRev := (bits.Reverse64(a) >> 55) << logTileSize
			for c := uint64(0); c < tileSize; c++ {
				idx := (a << (bShift)) | (b << logTileSize) | c
				t[aRev|c] = buf[idx]
			}
		}

		for c := uint64(0); c < tileSize; c++ {
			cRev := (bits.Reverse64(c) >> 55 << bShift) | bRev
			for aRev := uint64(0); aRev < tileSize; aRev++ {
				a := bits.Reverse64(aRev) >> 55
				idx := (a << bShift) | (b << logTileSize) | c
				idxRev := cRev | aRev
				if idx < idxRev {
					tIdx := (aRev << logTileSize) | c
					buf[idxRev], t[tIdx] = t[tIdx], buf[idxRev]
				}
			}
		}

		for a := uint64(0); a < tileSize; a++ {
			aRev := bits.Reverse64(a) >> 55
			for c := uint64(0); c < tileSize; c++ {
				cRev := (bits.Reverse64(c) >> 55) << bShift
				idx := (a << bShift) | (b << logTileSize) | c
				idxRev := cRev | bRev | aRev
				if idx < idxRev {
					tIdx := (aRev << logTileSize) | c
					buf[idx], t[tIdx] = t[tIdx], buf[idx]
				}
			}
		}
	}

}

func bitReverseCobraInPlace_9_25(buf []fr.Element) {
	const (
		logTileSize = uint64(9)
		tileSize    = uint64(1) << logTileSize
		logN        = 25
		logBLen     = logN - 2*logTileSize
		bShift      = logBLen + logTileSize
		bLen        = uint64(1) << logBLen
	)

	var t [tileSize * tileSize]fr.Element

	for b := uint64(0); b < bLen; b++ {
		bRev := (bits.Reverse64(b) >> (64 - logBLen)) << logTileSize

		for a := uint64(0); a < tileSize; a++ {
			aRev := (bits.Reverse64(a) >> 55) << logTileSize
			for c := uint64(0); c < tileSize; c++ {
				idx := (a << (bShift)) | (b << logTileSize) | c
				t[aRev|c] = buf[idx]
			}
		}

		for c := uint64(0); c < tileSize; c++ {
			cRev := (bits.Reverse64(c) >> 55 << bShift) | bRev
			for aRev := uint64(0); aRev < tileSize; aRev++ {
				a := bits.Reverse64(aRev) >> 55
				idx := (a << bShift) | (b << logTileSize) | c
				idxRev := cRev | aRev
				if idx < idxRev {
					tIdx := (aRev << logTileSize) | c
					buf[idxRev], t[tIdx] = t[tIdx], buf[idxRev]
				}
			}
		}

		for a := uint64(0); a < tileSize; a++ {
			aRev := bits.Reverse64(a) >> 55
			for c := uint64(0); c < tileSize; c++ {
				cRev := (bits.Reverse64(c) >> 55) << bShift
				idx := (a << bShift) | (b << logTileSize) | c
				idxRev := cRev | bRev | aRev
				if idx < idxRev {
					tIdx := (aRev << logTileSize) | c
					buf[idx], t[tIdx] = t[tIdx], buf[idx]
				}
			}
		}
	}

}

func bitReverseCobraInPlace_9_26(buf []fr.Element) {
	const (
		logTileSize = uint64(9)
		tileSize    = uint64(1) << logTileSize
		logN        = 26
		logBLen     = logN - 2*logTileSize
		bShift      = logBLen + logTileSize
		bLen        = uint64(1) << logBLen
	)

	var t [tileSize * tileSize]fr.Element

	for b := uint64(0); b < bLen; b++ {
		bRev := (bits.Reverse64(b) >> (64 - logBLen)) << logTileSize

		for a := uint64(0); a < tileSize; a++ {
			aRev := (bits.Reverse64(a) >> 55) << logTileSize
			for c := uint64(0); c < tileSize; c++ {
				idx := (a << (bShift)) | (b << logTileSize) | c
				t[aRev|c] = buf[idx]
			}
		}

		for c := uint64(0); c < tileSize; c++ {
			cRev := (bits.Reverse64(c) >> 55 << bShift) | bRev
			for aRev := uint64(0); aRev < tileSize; aRev++ {
				a := bits.Reverse64(aRev) >> 55
				idx := (a << bShift) | (b << logTileSize) | c
				idxRev := cRev | aRev
				if idx < idxRev {
					tIdx := (aRev << logTileSize) | c
					buf[idxRev], t[tIdx] = t[tIdx], buf[idxRev]
				}
			}
		}

		for a := uint64(0); a < tileSize; a++ {
			aRev := bits.Reverse64(a) >> 55
			for c := uint64(0); c < tileSize; c++ {
				cRev := (bits.Reverse64(c) >> 55) << bShift
				idx := (a << bShift) | (b << logTileSize) | c
				idxRev := cRev | bRev | aRev
				if idx < idxRev {
					tIdx := (aRev << logTileSize) | c
					buf[idx], t[tIdx] = t[tIdx], buf[idx]
				}
			}
		}
	}

}

func bitReverseCobraInPlace_9_27(buf []fr.Element) {
	const (
		logTileSize = uint64(9)
		tileSize    = uint64(1) << logTileSize
		logN        = 27
		logBLen     = logN - 2*logTileSize
		bShift      = logBLen + logTileSize
		bLen        = uint64(1) << logBLen
	)

	var t [tileSize * tileSize]fr.Element

	for b := uint64(0); b < bLen; b++ {
		bRev := (bits.Reverse64(b) >> (64 - logBLen)) << logTileSize

		for a := uint64(0); a < tileSize; a++ {
			aRev := (bits.Reverse64(a) >> 55) << logTileSize
			for c := uint64(0); c < tileSize; c++ {
				idx := (a << (bShift)) | (b << logTileSize) | c
				t[aRev|c] = buf[idx]
			}
		}

		for c := uint64(0); c < tileSize; c++ {
			cRev := (bits.Reverse64(c) >> 55 << bShift) | bRev
			for aRev := uint64(0); aRev < tileSize; aRev++ {
				a := bits.Reverse64(aRev) >> 55
				idx := (a << bShift) | (b << logTileSize) | c
				idxRev := cRev | aRev
				if idx < idxRev {
					tIdx := (aRev << logTileSize) | c
					buf[idxRev], t[tIdx] = t[tIdx], buf[idxRev]
				}
			}
		}

		for a := uint64(0); a < tileSize; a++ {
			aRev := bits.Reverse64(a) >> 55
			for c := uint64(0); c < tileSize; c++ {
				cRev := (bits.Reverse64(c) >> 55) << bShift
				idx := (a << bShift) | (b << logTileSize) | c
				idxRev := cRev | bRev | aRev
				if idx < idxRev {
					tIdx := (aRev << logTileSize) | c
					buf[idx], t[tIdx] = t[tIdx], buf[idx]
				}
			}
		}
	}

}

func BitReverseNew(buf []fr.Element) {
	switch len(buf) {
	case 1 << 21:
		bitReverseCobraInPlace_9_21(buf)
	case 1 << 22:
		bitReverseCobraInPlace_9_22(buf)
	case 1 << 23:
		bitReverseCobraInPlace_9_23(buf)
	case 1 << 24:
		bitReverseCobraInPlace_9_24(buf)
	case 1 << 25:
		bitReverseCobraInPlace_9_25(buf)
	case 1 << 26:
		bitReverseCobraInPlace_9_26(buf)
	case 1 << 27:
		bitReverseCobraInPlace_9_27(buf)
	default:
		if len(buf) > 1<<27 {
			bitReverseCobraInPlace(buf)
		} else {
			BitReverse(buf)
		}
	}
}
