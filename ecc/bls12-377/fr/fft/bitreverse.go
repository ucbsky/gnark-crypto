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

	"github.com/consensys/gnark-crypto/ecc/bls12-377/fr"
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

func BitReverseCobraInPlace(buf []fr.Element) {
	// TODO @gbotrel do a switch on the len(buf)
	bitReverseCobraInPlace(buf)
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

func bitReverseCobraInPlace_4(buf []fr.Element) {
	const logTileSize = uint64(4)
	const tileSize = uint64(1) << logTileSize

	logN := uint64(bits.Len64(uint64(len(buf))) - 1)
	logBLen := logN - 2*logTileSize
	bShift := logBLen + logTileSize
	bLen := uint64(1) << logBLen

	var t [tileSize * tileSize]fr.Element

	for b := uint64(0); b < bLen; b++ {
		bRev := bits.Reverse64(b) >> (64 - logBLen)

		for a := uint64(0); a < tileSize; a++ {
			aRev := bits.Reverse64(a) >> 60
			for c := uint64(0); c < tileSize; c++ {
				tIdx := (aRev << logTileSize) | c
				idx := (a << (bShift)) | (b << logTileSize) | c
				t[tIdx] = buf[idx]
			}
		}

		for c := uint64(0); c < tileSize; c++ {
			cRev := bits.Reverse64(c) >> 60
			for aRev := uint64(0); aRev < tileSize; aRev++ {
				a := bits.Reverse64(aRev) >> 60
				idx := (a << (bShift)) | (b << logTileSize) | c
				idxRev := (cRev << (bShift)) | (bRev << logTileSize) | aRev
				if idx < idxRev {
					tIdx := (aRev << logTileSize) | c
					buf[idxRev], t[tIdx] = t[tIdx], buf[idxRev]
				}
			}
		}

		for a := uint64(0); a < tileSize; a++ {
			aRev := bits.Reverse64(a) >> 60
			for c := uint64(0); c < tileSize; c++ {
				cRev := bits.Reverse64(c) >> 60
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

func bitReverseCobraInPlace_4unrolled(buf []fr.Element) {
	const logTileSize = uint64(4)
	const tileSize = uint64(1) << logTileSize

	logN := uint64(bits.Len64(uint64(len(buf))) - 1)
	logBLen := logN - 2*logTileSize
	bShift := logBLen + logTileSize
	bLen := uint64(1) << logBLen

	var t [tileSize * tileSize]fr.Element

	for b := uint64(0); b < bLen; b++ {

		// for a := uint64(0); a < tileSize; a++ {
		// 	aRev := bits.Reverse64(a) >> 60
		// 	for c := uint64(0); c < tileSize; c++ {
		// 		tIdx := (aRev << logTileSize) | c
		// 		idx := (a << (bShift)) | (b << logTileSize) | c
		// 		t[tIdx] = buf[idx]
		// 	}
		// }
		t[0] = buf[(b << logTileSize)]
		t[1] = buf[(b<<logTileSize)|1]
		t[2] = buf[(b<<logTileSize)|2]
		t[3] = buf[(b<<logTileSize)|3]
		t[4] = buf[(b<<logTileSize)|4]
		t[5] = buf[(b<<logTileSize)|5]
		t[6] = buf[(b<<logTileSize)|6]
		t[7] = buf[(b<<logTileSize)|7]
		t[8] = buf[(b<<logTileSize)|8]
		t[9] = buf[(b<<logTileSize)|9]
		t[10] = buf[(b<<logTileSize)|10]
		t[11] = buf[(b<<logTileSize)|11]
		t[12] = buf[(b<<logTileSize)|12]
		t[13] = buf[(b<<logTileSize)|13]
		t[14] = buf[(b<<logTileSize)|14]
		t[15] = buf[(b<<logTileSize)|15]
		t[128] = buf[(1<<bShift)|(b<<logTileSize)]
		t[129] = buf[(1<<bShift)|(b<<logTileSize)|1]
		t[130] = buf[(1<<bShift)|(b<<logTileSize)|2]
		t[131] = buf[(1<<bShift)|(b<<logTileSize)|3]
		t[132] = buf[(1<<bShift)|(b<<logTileSize)|4]
		t[133] = buf[(1<<bShift)|(b<<logTileSize)|5]
		t[134] = buf[(1<<bShift)|(b<<logTileSize)|6]
		t[135] = buf[(1<<bShift)|(b<<logTileSize)|7]
		t[136] = buf[(1<<bShift)|(b<<logTileSize)|8]
		t[137] = buf[(1<<bShift)|(b<<logTileSize)|9]
		t[138] = buf[(1<<bShift)|(b<<logTileSize)|10]
		t[139] = buf[(1<<bShift)|(b<<logTileSize)|11]
		t[140] = buf[(1<<bShift)|(b<<logTileSize)|12]
		t[141] = buf[(1<<bShift)|(b<<logTileSize)|13]
		t[142] = buf[(1<<bShift)|(b<<logTileSize)|14]
		t[143] = buf[(1<<bShift)|(b<<logTileSize)|15]
		t[64] = buf[(2<<bShift)|(b<<logTileSize)]
		t[65] = buf[(2<<bShift)|(b<<logTileSize)|1]
		t[66] = buf[(2<<bShift)|(b<<logTileSize)|2]
		t[67] = buf[(2<<bShift)|(b<<logTileSize)|3]
		t[68] = buf[(2<<bShift)|(b<<logTileSize)|4]
		t[69] = buf[(2<<bShift)|(b<<logTileSize)|5]
		t[70] = buf[(2<<bShift)|(b<<logTileSize)|6]
		t[71] = buf[(2<<bShift)|(b<<logTileSize)|7]
		t[72] = buf[(2<<bShift)|(b<<logTileSize)|8]
		t[73] = buf[(2<<bShift)|(b<<logTileSize)|9]
		t[74] = buf[(2<<bShift)|(b<<logTileSize)|10]
		t[75] = buf[(2<<bShift)|(b<<logTileSize)|11]
		t[76] = buf[(2<<bShift)|(b<<logTileSize)|12]
		t[77] = buf[(2<<bShift)|(b<<logTileSize)|13]
		t[78] = buf[(2<<bShift)|(b<<logTileSize)|14]
		t[79] = buf[(2<<bShift)|(b<<logTileSize)|15]
		t[192] = buf[(3<<bShift)|(b<<logTileSize)]
		t[193] = buf[(3<<bShift)|(b<<logTileSize)|1]
		t[194] = buf[(3<<bShift)|(b<<logTileSize)|2]
		t[195] = buf[(3<<bShift)|(b<<logTileSize)|3]
		t[196] = buf[(3<<bShift)|(b<<logTileSize)|4]
		t[197] = buf[(3<<bShift)|(b<<logTileSize)|5]
		t[198] = buf[(3<<bShift)|(b<<logTileSize)|6]
		t[199] = buf[(3<<bShift)|(b<<logTileSize)|7]
		t[200] = buf[(3<<bShift)|(b<<logTileSize)|8]
		t[201] = buf[(3<<bShift)|(b<<logTileSize)|9]
		t[202] = buf[(3<<bShift)|(b<<logTileSize)|10]
		t[203] = buf[(3<<bShift)|(b<<logTileSize)|11]
		t[204] = buf[(3<<bShift)|(b<<logTileSize)|12]
		t[205] = buf[(3<<bShift)|(b<<logTileSize)|13]
		t[206] = buf[(3<<bShift)|(b<<logTileSize)|14]
		t[207] = buf[(3<<bShift)|(b<<logTileSize)|15]
		t[32] = buf[(4<<bShift)|(b<<logTileSize)]
		t[33] = buf[(4<<bShift)|(b<<logTileSize)|1]
		t[34] = buf[(4<<bShift)|(b<<logTileSize)|2]
		t[35] = buf[(4<<bShift)|(b<<logTileSize)|3]
		t[36] = buf[(4<<bShift)|(b<<logTileSize)|4]
		t[37] = buf[(4<<bShift)|(b<<logTileSize)|5]
		t[38] = buf[(4<<bShift)|(b<<logTileSize)|6]
		t[39] = buf[(4<<bShift)|(b<<logTileSize)|7]
		t[40] = buf[(4<<bShift)|(b<<logTileSize)|8]
		t[41] = buf[(4<<bShift)|(b<<logTileSize)|9]
		t[42] = buf[(4<<bShift)|(b<<logTileSize)|10]
		t[43] = buf[(4<<bShift)|(b<<logTileSize)|11]
		t[44] = buf[(4<<bShift)|(b<<logTileSize)|12]
		t[45] = buf[(4<<bShift)|(b<<logTileSize)|13]
		t[46] = buf[(4<<bShift)|(b<<logTileSize)|14]
		t[47] = buf[(4<<bShift)|(b<<logTileSize)|15]
		t[160] = buf[(5<<bShift)|(b<<logTileSize)]
		t[161] = buf[(5<<bShift)|(b<<logTileSize)|1]
		t[162] = buf[(5<<bShift)|(b<<logTileSize)|2]
		t[163] = buf[(5<<bShift)|(b<<logTileSize)|3]
		t[164] = buf[(5<<bShift)|(b<<logTileSize)|4]
		t[165] = buf[(5<<bShift)|(b<<logTileSize)|5]
		t[166] = buf[(5<<bShift)|(b<<logTileSize)|6]
		t[167] = buf[(5<<bShift)|(b<<logTileSize)|7]
		t[168] = buf[(5<<bShift)|(b<<logTileSize)|8]
		t[169] = buf[(5<<bShift)|(b<<logTileSize)|9]
		t[170] = buf[(5<<bShift)|(b<<logTileSize)|10]
		t[171] = buf[(5<<bShift)|(b<<logTileSize)|11]
		t[172] = buf[(5<<bShift)|(b<<logTileSize)|12]
		t[173] = buf[(5<<bShift)|(b<<logTileSize)|13]
		t[174] = buf[(5<<bShift)|(b<<logTileSize)|14]
		t[175] = buf[(5<<bShift)|(b<<logTileSize)|15]
		t[96] = buf[(6<<bShift)|(b<<logTileSize)]
		t[97] = buf[(6<<bShift)|(b<<logTileSize)|1]
		t[98] = buf[(6<<bShift)|(b<<logTileSize)|2]
		t[99] = buf[(6<<bShift)|(b<<logTileSize)|3]
		t[100] = buf[(6<<bShift)|(b<<logTileSize)|4]
		t[101] = buf[(6<<bShift)|(b<<logTileSize)|5]
		t[102] = buf[(6<<bShift)|(b<<logTileSize)|6]
		t[103] = buf[(6<<bShift)|(b<<logTileSize)|7]
		t[104] = buf[(6<<bShift)|(b<<logTileSize)|8]
		t[105] = buf[(6<<bShift)|(b<<logTileSize)|9]
		t[106] = buf[(6<<bShift)|(b<<logTileSize)|10]
		t[107] = buf[(6<<bShift)|(b<<logTileSize)|11]
		t[108] = buf[(6<<bShift)|(b<<logTileSize)|12]
		t[109] = buf[(6<<bShift)|(b<<logTileSize)|13]
		t[110] = buf[(6<<bShift)|(b<<logTileSize)|14]
		t[111] = buf[(6<<bShift)|(b<<logTileSize)|15]
		t[224] = buf[(7<<bShift)|(b<<logTileSize)]
		t[225] = buf[(7<<bShift)|(b<<logTileSize)|1]
		t[226] = buf[(7<<bShift)|(b<<logTileSize)|2]
		t[227] = buf[(7<<bShift)|(b<<logTileSize)|3]
		t[228] = buf[(7<<bShift)|(b<<logTileSize)|4]
		t[229] = buf[(7<<bShift)|(b<<logTileSize)|5]
		t[230] = buf[(7<<bShift)|(b<<logTileSize)|6]
		t[231] = buf[(7<<bShift)|(b<<logTileSize)|7]
		t[232] = buf[(7<<bShift)|(b<<logTileSize)|8]
		t[233] = buf[(7<<bShift)|(b<<logTileSize)|9]
		t[234] = buf[(7<<bShift)|(b<<logTileSize)|10]
		t[235] = buf[(7<<bShift)|(b<<logTileSize)|11]
		t[236] = buf[(7<<bShift)|(b<<logTileSize)|12]
		t[237] = buf[(7<<bShift)|(b<<logTileSize)|13]
		t[238] = buf[(7<<bShift)|(b<<logTileSize)|14]
		t[239] = buf[(7<<bShift)|(b<<logTileSize)|15]
		t[16] = buf[(8<<bShift)|(b<<logTileSize)]
		t[17] = buf[(8<<bShift)|(b<<logTileSize)|1]
		t[18] = buf[(8<<bShift)|(b<<logTileSize)|2]
		t[19] = buf[(8<<bShift)|(b<<logTileSize)|3]
		t[20] = buf[(8<<bShift)|(b<<logTileSize)|4]
		t[21] = buf[(8<<bShift)|(b<<logTileSize)|5]
		t[22] = buf[(8<<bShift)|(b<<logTileSize)|6]
		t[23] = buf[(8<<bShift)|(b<<logTileSize)|7]
		t[24] = buf[(8<<bShift)|(b<<logTileSize)|8]
		t[25] = buf[(8<<bShift)|(b<<logTileSize)|9]
		t[26] = buf[(8<<bShift)|(b<<logTileSize)|10]
		t[27] = buf[(8<<bShift)|(b<<logTileSize)|11]
		t[28] = buf[(8<<bShift)|(b<<logTileSize)|12]
		t[29] = buf[(8<<bShift)|(b<<logTileSize)|13]
		t[30] = buf[(8<<bShift)|(b<<logTileSize)|14]
		t[31] = buf[(8<<bShift)|(b<<logTileSize)|15]
		t[144] = buf[(9<<bShift)|(b<<logTileSize)]
		t[145] = buf[(9<<bShift)|(b<<logTileSize)|1]
		t[146] = buf[(9<<bShift)|(b<<logTileSize)|2]
		t[147] = buf[(9<<bShift)|(b<<logTileSize)|3]
		t[148] = buf[(9<<bShift)|(b<<logTileSize)|4]
		t[149] = buf[(9<<bShift)|(b<<logTileSize)|5]
		t[150] = buf[(9<<bShift)|(b<<logTileSize)|6]
		t[151] = buf[(9<<bShift)|(b<<logTileSize)|7]
		t[152] = buf[(9<<bShift)|(b<<logTileSize)|8]
		t[153] = buf[(9<<bShift)|(b<<logTileSize)|9]
		t[154] = buf[(9<<bShift)|(b<<logTileSize)|10]
		t[155] = buf[(9<<bShift)|(b<<logTileSize)|11]
		t[156] = buf[(9<<bShift)|(b<<logTileSize)|12]
		t[157] = buf[(9<<bShift)|(b<<logTileSize)|13]
		t[158] = buf[(9<<bShift)|(b<<logTileSize)|14]
		t[159] = buf[(9<<bShift)|(b<<logTileSize)|15]
		t[80] = buf[(10<<bShift)|(b<<logTileSize)]
		t[81] = buf[(10<<bShift)|(b<<logTileSize)|1]
		t[82] = buf[(10<<bShift)|(b<<logTileSize)|2]
		t[83] = buf[(10<<bShift)|(b<<logTileSize)|3]
		t[84] = buf[(10<<bShift)|(b<<logTileSize)|4]
		t[85] = buf[(10<<bShift)|(b<<logTileSize)|5]
		t[86] = buf[(10<<bShift)|(b<<logTileSize)|6]
		t[87] = buf[(10<<bShift)|(b<<logTileSize)|7]
		t[88] = buf[(10<<bShift)|(b<<logTileSize)|8]
		t[89] = buf[(10<<bShift)|(b<<logTileSize)|9]
		t[90] = buf[(10<<bShift)|(b<<logTileSize)|10]
		t[91] = buf[(10<<bShift)|(b<<logTileSize)|11]
		t[92] = buf[(10<<bShift)|(b<<logTileSize)|12]
		t[93] = buf[(10<<bShift)|(b<<logTileSize)|13]
		t[94] = buf[(10<<bShift)|(b<<logTileSize)|14]
		t[95] = buf[(10<<bShift)|(b<<logTileSize)|15]
		t[208] = buf[(11<<bShift)|(b<<logTileSize)]
		t[209] = buf[(11<<bShift)|(b<<logTileSize)|1]
		t[210] = buf[(11<<bShift)|(b<<logTileSize)|2]
		t[211] = buf[(11<<bShift)|(b<<logTileSize)|3]
		t[212] = buf[(11<<bShift)|(b<<logTileSize)|4]
		t[213] = buf[(11<<bShift)|(b<<logTileSize)|5]
		t[214] = buf[(11<<bShift)|(b<<logTileSize)|6]
		t[215] = buf[(11<<bShift)|(b<<logTileSize)|7]
		t[216] = buf[(11<<bShift)|(b<<logTileSize)|8]
		t[217] = buf[(11<<bShift)|(b<<logTileSize)|9]
		t[218] = buf[(11<<bShift)|(b<<logTileSize)|10]
		t[219] = buf[(11<<bShift)|(b<<logTileSize)|11]
		t[220] = buf[(11<<bShift)|(b<<logTileSize)|12]
		t[221] = buf[(11<<bShift)|(b<<logTileSize)|13]
		t[222] = buf[(11<<bShift)|(b<<logTileSize)|14]
		t[223] = buf[(11<<bShift)|(b<<logTileSize)|15]
		t[48] = buf[(12<<bShift)|(b<<logTileSize)]
		t[49] = buf[(12<<bShift)|(b<<logTileSize)|1]
		t[50] = buf[(12<<bShift)|(b<<logTileSize)|2]
		t[51] = buf[(12<<bShift)|(b<<logTileSize)|3]
		t[52] = buf[(12<<bShift)|(b<<logTileSize)|4]
		t[53] = buf[(12<<bShift)|(b<<logTileSize)|5]
		t[54] = buf[(12<<bShift)|(b<<logTileSize)|6]
		t[55] = buf[(12<<bShift)|(b<<logTileSize)|7]
		t[56] = buf[(12<<bShift)|(b<<logTileSize)|8]
		t[57] = buf[(12<<bShift)|(b<<logTileSize)|9]
		t[58] = buf[(12<<bShift)|(b<<logTileSize)|10]
		t[59] = buf[(12<<bShift)|(b<<logTileSize)|11]
		t[60] = buf[(12<<bShift)|(b<<logTileSize)|12]
		t[61] = buf[(12<<bShift)|(b<<logTileSize)|13]
		t[62] = buf[(12<<bShift)|(b<<logTileSize)|14]
		t[63] = buf[(12<<bShift)|(b<<logTileSize)|15]
		t[176] = buf[(13<<bShift)|(b<<logTileSize)]
		t[177] = buf[(13<<bShift)|(b<<logTileSize)|1]
		t[178] = buf[(13<<bShift)|(b<<logTileSize)|2]
		t[179] = buf[(13<<bShift)|(b<<logTileSize)|3]
		t[180] = buf[(13<<bShift)|(b<<logTileSize)|4]
		t[181] = buf[(13<<bShift)|(b<<logTileSize)|5]
		t[182] = buf[(13<<bShift)|(b<<logTileSize)|6]
		t[183] = buf[(13<<bShift)|(b<<logTileSize)|7]
		t[184] = buf[(13<<bShift)|(b<<logTileSize)|8]
		t[185] = buf[(13<<bShift)|(b<<logTileSize)|9]
		t[186] = buf[(13<<bShift)|(b<<logTileSize)|10]
		t[187] = buf[(13<<bShift)|(b<<logTileSize)|11]
		t[188] = buf[(13<<bShift)|(b<<logTileSize)|12]
		t[189] = buf[(13<<bShift)|(b<<logTileSize)|13]
		t[190] = buf[(13<<bShift)|(b<<logTileSize)|14]
		t[191] = buf[(13<<bShift)|(b<<logTileSize)|15]
		t[112] = buf[(14<<bShift)|(b<<logTileSize)]
		t[113] = buf[(14<<bShift)|(b<<logTileSize)|1]
		t[114] = buf[(14<<bShift)|(b<<logTileSize)|2]
		t[115] = buf[(14<<bShift)|(b<<logTileSize)|3]
		t[116] = buf[(14<<bShift)|(b<<logTileSize)|4]
		t[117] = buf[(14<<bShift)|(b<<logTileSize)|5]
		t[118] = buf[(14<<bShift)|(b<<logTileSize)|6]
		t[119] = buf[(14<<bShift)|(b<<logTileSize)|7]
		t[120] = buf[(14<<bShift)|(b<<logTileSize)|8]
		t[121] = buf[(14<<bShift)|(b<<logTileSize)|9]
		t[122] = buf[(14<<bShift)|(b<<logTileSize)|10]
		t[123] = buf[(14<<bShift)|(b<<logTileSize)|11]
		t[124] = buf[(14<<bShift)|(b<<logTileSize)|12]
		t[125] = buf[(14<<bShift)|(b<<logTileSize)|13]
		t[126] = buf[(14<<bShift)|(b<<logTileSize)|14]
		t[127] = buf[(14<<bShift)|(b<<logTileSize)|15]
		t[240] = buf[(15<<bShift)|(b<<logTileSize)]
		t[241] = buf[(15<<bShift)|(b<<logTileSize)|1]
		t[242] = buf[(15<<bShift)|(b<<logTileSize)|2]
		t[243] = buf[(15<<bShift)|(b<<logTileSize)|3]
		t[244] = buf[(15<<bShift)|(b<<logTileSize)|4]
		t[245] = buf[(15<<bShift)|(b<<logTileSize)|5]
		t[246] = buf[(15<<bShift)|(b<<logTileSize)|6]
		t[247] = buf[(15<<bShift)|(b<<logTileSize)|7]
		t[248] = buf[(15<<bShift)|(b<<logTileSize)|8]
		t[249] = buf[(15<<bShift)|(b<<logTileSize)|9]
		t[250] = buf[(15<<bShift)|(b<<logTileSize)|10]
		t[251] = buf[(15<<bShift)|(b<<logTileSize)|11]
		t[252] = buf[(15<<bShift)|(b<<logTileSize)|12]
		t[253] = buf[(15<<bShift)|(b<<logTileSize)|13]
		t[254] = buf[(15<<bShift)|(b<<logTileSize)|14]
		t[255] = buf[(15<<bShift)|(b<<logTileSize)|15]

		bRev := bits.Reverse64(b) >> (64 - logBLen)

		for c := uint64(0); c < tileSize; c++ {
			cRev := bits.Reverse64(c) >> 60
			for aRev := uint64(0); aRev < tileSize; aRev++ {
				a := bits.Reverse64(aRev) >> 60
				idx := (a << (bShift)) | (b << logTileSize) | c
				idxRev := (cRev << (bShift)) | (bRev << logTileSize) | aRev
				if idx < idxRev {
					tIdx := (aRev << logTileSize) | c
					buf[idxRev], t[tIdx] = t[tIdx], buf[idxRev]
				}
			}
		}

		for a := uint64(0); a < tileSize; a++ {
			aRev := bits.Reverse64(a) >> 60
			for c := uint64(0); c < tileSize; c++ {
				cRev := bits.Reverse64(c) >> 60
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

func bitReverseCobraInPlace_5(buf []fr.Element) {
	const logTileSize = uint64(5)
	const tileSize = uint64(1) << logTileSize

	logN := uint64(bits.Len64(uint64(len(buf))) - 1)
	logBLen := logN - 2*logTileSize
	bShift := logBLen + logTileSize
	bLen := uint64(1) << logBLen

	var t [tileSize * tileSize]fr.Element

	for b := uint64(0); b < bLen; b++ {
		bRev := bits.Reverse64(b) >> (64 - logBLen)

		for a := uint64(0); a < tileSize; a++ {
			aRev := bits.Reverse64(a) >> 59
			for c := uint64(0); c < tileSize; c++ {
				tIdx := (aRev << logTileSize) | c
				idx := (a << (bShift)) | (b << logTileSize) | c
				t[tIdx] = buf[idx]
			}
		}

		for c := uint64(0); c < tileSize; c++ {
			cRev := bits.Reverse64(c) >> 59
			for aRev := uint64(0); aRev < tileSize; aRev++ {
				a := bits.Reverse64(aRev) >> 59
				idx := (a << (bShift)) | (b << logTileSize) | c
				idxRev := (cRev << (bShift)) | (bRev << logTileSize) | aRev
				if idx < idxRev {
					tIdx := (aRev << logTileSize) | c
					buf[idxRev], t[tIdx] = t[tIdx], buf[idxRev]
				}
			}
		}

		for a := uint64(0); a < tileSize; a++ {
			aRev := bits.Reverse64(a) >> 59
			for c := uint64(0); c < tileSize; c++ {
				cRev := bits.Reverse64(c) >> 59
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

func bitReverseCobraInPlace_6(buf []fr.Element) {
	const logTileSize = uint64(6)
	const tileSize = uint64(1) << logTileSize

	logN := uint64(bits.Len64(uint64(len(buf))) - 1)
	logBLen := logN - 2*logTileSize
	bShift := logBLen + logTileSize
	bLen := uint64(1) << logBLen

	var t [tileSize * tileSize]fr.Element

	for b := uint64(0); b < bLen; b++ {
		bRev := bits.Reverse64(b) >> (64 - logBLen)

		for a := uint64(0); a < tileSize; a++ {
			aRev := bits.Reverse64(a) >> 58
			for c := uint64(0); c < tileSize; c++ {
				tIdx := (aRev << logTileSize) | c
				idx := (a << (bShift)) | (b << logTileSize) | c
				t[tIdx] = buf[idx]
			}
		}

		for c := uint64(0); c < tileSize; c++ {
			cRev := bits.Reverse64(c) >> 58
			for aRev := uint64(0); aRev < tileSize; aRev++ {
				a := bits.Reverse64(aRev) >> 58
				idx := (a << (bShift)) | (b << logTileSize) | c
				idxRev := (cRev << (bShift)) | (bRev << logTileSize) | aRev
				if idx < idxRev {
					tIdx := (aRev << logTileSize) | c
					buf[idxRev], t[tIdx] = t[tIdx], buf[idxRev]
				}
			}
		}

		for a := uint64(0); a < tileSize; a++ {
			aRev := bits.Reverse64(a) >> 58
			for c := uint64(0); c < tileSize; c++ {
				cRev := bits.Reverse64(c) >> 58
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

func bitReverseCobraInPlace_7(buf []fr.Element) {
	const logTileSize = uint64(7)
	const tileSize = uint64(1) << logTileSize

	logN := uint64(bits.Len64(uint64(len(buf))) - 1)
	logBLen := logN - 2*logTileSize
	bShift := logBLen + logTileSize
	bLen := uint64(1) << logBLen

	var t [tileSize * tileSize]fr.Element

	for b := uint64(0); b < bLen; b++ {
		bRev := bits.Reverse64(b) >> (64 - logBLen)

		for a := uint64(0); a < tileSize; a++ {
			aRev := bits.Reverse64(a) >> 57
			for c := uint64(0); c < tileSize; c++ {
				tIdx := (aRev << logTileSize) | c
				idx := (a << (bShift)) | (b << logTileSize) | c
				t[tIdx] = buf[idx]
			}
		}

		for c := uint64(0); c < tileSize; c++ {
			cRev := bits.Reverse64(c) >> 57
			for aRev := uint64(0); aRev < tileSize; aRev++ {
				a := bits.Reverse64(aRev) >> 57
				idx := (a << (bShift)) | (b << logTileSize) | c
				idxRev := (cRev << (bShift)) | (bRev << logTileSize) | aRev
				if idx < idxRev {
					tIdx := (aRev << logTileSize) | c
					buf[idxRev], t[tIdx] = t[tIdx], buf[idxRev]
				}
			}
		}

		for a := uint64(0); a < tileSize; a++ {
			aRev := bits.Reverse64(a) >> 57
			for c := uint64(0); c < tileSize; c++ {
				cRev := bits.Reverse64(c) >> 57
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

func bitReverseCobraInPlace_8(buf []fr.Element) {
	const logTileSize = uint64(8)
	const tileSize = uint64(1) << logTileSize

	logN := uint64(bits.Len64(uint64(len(buf))) - 1)
	logBLen := logN - 2*logTileSize
	bShift := logBLen + logTileSize
	bLen := uint64(1) << logBLen

	var t [tileSize * tileSize]fr.Element

	for b := uint64(0); b < bLen; b++ {
		bRev := bits.Reverse64(b) >> (64 - logBLen)

		for a := uint64(0); a < tileSize; a++ {
			aRev := bits.Reverse64(a) >> 56
			for c := uint64(0); c < tileSize; c++ {
				tIdx := (aRev << logTileSize) | c
				idx := (a << (bShift)) | (b << logTileSize) | c
				t[tIdx] = buf[idx]
			}
		}

		for c := uint64(0); c < tileSize; c++ {
			cRev := bits.Reverse64(c) >> 56
			for aRev := uint64(0); aRev < tileSize; aRev++ {
				a := bits.Reverse64(aRev) >> 56
				idx := (a << (bShift)) | (b << logTileSize) | c
				idxRev := (cRev << (bShift)) | (bRev << logTileSize) | aRev
				if idx < idxRev {
					tIdx := (aRev << logTileSize) | c
					buf[idxRev], t[tIdx] = t[tIdx], buf[idxRev]
				}
			}
		}

		for a := uint64(0); a < tileSize; a++ {
			aRev := bits.Reverse64(a) >> 56
			for c := uint64(0); c < tileSize; c++ {
				cRev := bits.Reverse64(c) >> 56
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

func bitReverseCobraInPlace_9(buf []fr.Element) {
	const logTileSize = uint64(9)
	const tileSize = uint64(1) << logTileSize

	logN := uint64(bits.Len64(uint64(len(buf))) - 1)
	logBLen := logN - 2*logTileSize
	bShift := logBLen + logTileSize
	bLen := uint64(1) << logBLen

	var t [tileSize * tileSize]fr.Element

	for b := uint64(0); b < bLen; b++ {
		bRev := bits.Reverse64(b) >> (64 - logBLen)

		for a := uint64(0); a < tileSize; a++ {
			aRev := bits.Reverse64(a) >> 55
			for c := uint64(0); c < tileSize; c++ {
				tIdx := (aRev << logTileSize) | c
				idx := (a << (bShift)) | (b << logTileSize) | c
				t[tIdx] = buf[idx]
			}
		}

		for c := uint64(0); c < tileSize; c++ {
			cRev := bits.Reverse64(c) >> 55
			for aRev := uint64(0); aRev < tileSize; aRev++ {
				a := bits.Reverse64(aRev) >> 55
				idx := (a << (bShift)) | (b << logTileSize) | c
				idxRev := (cRev << (bShift)) | (bRev << logTileSize) | aRev
				if idx < idxRev {
					tIdx := (aRev << logTileSize) | c
					buf[idxRev], t[tIdx] = t[tIdx], buf[idxRev]
				}
			}
		}

		for a := uint64(0); a < tileSize; a++ {
			aRev := bits.Reverse64(a) >> 55
			for c := uint64(0); c < tileSize; c++ {
				cRev := bits.Reverse64(c) >> 55
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

func bitReverseCobraInPlace_10(buf []fr.Element) {
	const logTileSize = uint64(10)
	const tileSize = uint64(1) << logTileSize

	logN := uint64(bits.Len64(uint64(len(buf))) - 1)
	logBLen := logN - 2*logTileSize
	bShift := logBLen + logTileSize
	bLen := uint64(1) << logBLen

	var t [tileSize * tileSize]fr.Element

	for b := uint64(0); b < bLen; b++ {
		bRev := bits.Reverse64(b) >> (64 - logBLen)

		for a := uint64(0); a < tileSize; a++ {
			aRev := bits.Reverse64(a) >> 54
			for c := uint64(0); c < tileSize; c++ {
				tIdx := (aRev << logTileSize) | c
				idx := (a << (bShift)) | (b << logTileSize) | c
				t[tIdx] = buf[idx]
			}
		}

		for c := uint64(0); c < tileSize; c++ {
			cRev := bits.Reverse64(c) >> 54
			for aRev := uint64(0); aRev < tileSize; aRev++ {
				a := bits.Reverse64(aRev) >> 54
				idx := (a << (bShift)) | (b << logTileSize) | c
				idxRev := (cRev << (bShift)) | (bRev << logTileSize) | aRev
				if idx < idxRev {
					tIdx := (aRev << logTileSize) | c
					buf[idxRev], t[tIdx] = t[tIdx], buf[idxRev]
				}
			}
		}

		for a := uint64(0); a < tileSize; a++ {
			aRev := bits.Reverse64(a) >> 54
			for c := uint64(0); c < tileSize; c++ {
				cRev := bits.Reverse64(c) >> 54
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
