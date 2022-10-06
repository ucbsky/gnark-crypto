// Copyright 2020 ConsenSys Software Inc.
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

package sumcheck

import (
	"github.com/consensys/gnark-crypto/ecc/bn254/fr"
	"github.com/consensys/gnark-crypto/ecc/bn254/fr/polynomial"
)

// This does not make use of parallelism and represents polynomials as lists of coefficients
// It is currently geared towards arithmetic hashes. Once we have a more unified hash function interface, this can be generified.

// Claims to a multi-sumcheck statement. i.e. one of the form ∑_{0≤i<2ⁿ} fⱼ(i) = cⱼ for 1 ≤ j ≤ m.
// Later evolving into a claim of the form gⱼ = ∑_{0≤i<2ⁿ⁻ʲ} g(r₁, r₂, ..., rⱼ₋₁, Xⱼ, i...)
type Claims interface {
	Combine(a fr.Element) polynomial.Polynomial // Combine into the 0ᵗʰ sumcheck subclaim. Create g := ∑_{1≤j≤m} aʲ⁻¹fⱼ for which now we seek to prove ∑_{0≤i<2ⁿ} g(i) = c := ∑_{1≤j≤m} aʲ⁻¹cⱼ. Return g₁.
	Next(fr.Element) polynomial.Polynomial      // Return the evaluations gⱼ(k) for 1 ≤ k < degⱼ(g). Update the claim to gⱼ₊₁ for the input value as rⱼ
	VarsNum() int                               //number of variables
	ClaimsNum() int                             //number of claims
	ProveFinalEval(r []fr.Element) interface{}  //in case it is difficult for the verifier to compute g(r₁, ..., rₙ) on its own, the prover can provide the value and a proof
}

// LazyClaims is the Claims data structure on the verifier side. It is "lazy" in that it has to compute fewer things.
type LazyClaims interface {
	ClaimsNum() int                      // ClaimsNum = m
	VarsNum() int                        // VarsNum = n
	CombinedSum(a fr.Element) fr.Element // CombinedSum returns c = ∑_{1≤j≤m} aʲ⁻¹cⱼ
	Degree(i int) int                    //Degree of the total claim in the i'th variable
	VerifyFinalEval(r []fr.Element, combinationCoeff fr.Element, purportedValue fr.Element, proof interface{}) bool
}

// Proof of a multi-sumcheck statement.
type Proof struct {
	PartialSumPolys []polynomial.Polynomial `json:"partialSumPolys"`
	FinalEvalProof  interface{}             `json:"finalEvalProof"` //in case it is difficult for the verifier to compute g(r₁, ..., rₙ) on its own, the prover can provide the value and a proof
}

// TODO: User unfriendly. Fix
func elementSliceToInterfaceSlice(elementSlice []fr.Element) (interfaceSlice []interface{}) {

	interfaceSlice = make([]interface{}, len(elementSlice))
	for i := range elementSlice {
		interfaceSlice[i] = &elementSlice[i]
	}
	return
}

// Prove create a non-interactive sumcheck proof
// transcript must have a hash function specified and seeded with a
func Prove(claims Claims, transcript ArithmeticTranscript) (proof Proof) {
	// TODO: Are claims supposed to already be incorporated in the challengeSeed? Given the business with the commitments

	var combinationCoeff fr.Element
	if claims.ClaimsNum() >= 2 {
		combinationCoeff = transcript.Next()
	}

	varsNum := claims.VarsNum()
	proof.PartialSumPolys = make([]polynomial.Polynomial, varsNum)
	proof.PartialSumPolys[0] = claims.Combine(combinationCoeff)
	challenges := make([]fr.Element, varsNum)

	for j := 0; j+1 < varsNum; j++ {
		challenges[j] = transcript.Next(elementSliceToInterfaceSlice(proof.PartialSumPolys[j])...)
		proof.PartialSumPolys[j+1] = claims.Next(challenges[j])
	}

	challenges[varsNum-1] = transcript.Next(elementSliceToInterfaceSlice(proof.PartialSumPolys[varsNum-1])...)

	proof.FinalEvalProof = claims.ProveFinalEval(challenges)

	return
}

func Verify(claims LazyClaims, proof Proof, transcript ArithmeticTranscript) bool {
	var combinationCoeff fr.Element

	if claims.ClaimsNum() >= 2 {
		combinationCoeff = transcript.Next()
	}

	r := make([]fr.Element, claims.VarsNum())

	// Just so that there is enough room for gJ to be reused
	maxDegree := claims.Degree(0)
	for j := 1; j < claims.VarsNum(); j++ {
		if d := claims.Degree(j); d > maxDegree {
			maxDegree = d
		}
	}
	gJ := make(polynomial.Polynomial, maxDegree+1) //At the end of iteration j, gJ = ∑_{i < 2ⁿ⁻ʲ⁻¹} g(X₁, ..., Xⱼ₊₁, i...)		NOTE: n is shorthand for claims.VarsNum()
	gJR := claims.CombinedSum(combinationCoeff)    // At the beginning of iteration j, gJR = ∑_{i < 2ⁿ⁻ʲ} g(r₁, ..., rⱼ, i...)

	for j := 0; j < claims.VarsNum(); j++ {
		if len(proof.PartialSumPolys[j]) != claims.Degree(j) {
			return false //Malformed proof
		}
		copy(gJ[1:], proof.PartialSumPolys[j])
		gJ[0].Sub(&gJR, &proof.PartialSumPolys[j][0]) // Requirement that gⱼ(0) + gⱼ(1) = gⱼ₋₁(r)
		// gJ is ready

		//Prepare for the next iteration
		r[j] = transcript.Next(elementSliceToInterfaceSlice(proof.PartialSumPolys[j])...)
		// This is an extremely inefficient way of interpolating. TODO: Interpolate without symbolically computing a polynomial
		gJCoeffs := polynomial.InterpolateOnRange(gJ[:(claims.Degree(j) + 1)])
		gJR = gJCoeffs.Eval(&r[j])
	}

	return claims.VerifyFinalEval(r, combinationCoeff, gJR, proof.FinalEvalProof)
}

// -------- fiatshamir  --------- TODO: Replace with existing fiat-shamir impl

// This is an implementation of Fiat-Shamir optimized for in-circuit verifiers.
// i.e. the hashes used operate on and return field elements.

type ArithmeticTranscript interface {
	Update(...interface{})
	Next(...interface{}) fr.Element
	NextN(int, ...interface{}) []fr.Element
}

// This is a very bad fiat-shamir challenge generator
type MessageCounter struct {
	state   uint64
	step    uint64
	updated bool
}

func (m *MessageCounter) Update(i ...interface{}) {
	m.state += m.step
	m.updated = true
}

func (m *MessageCounter) Next(i ...interface{}) (challenge fr.Element) {
	if !m.updated || len(i) != 0 {
		m.Update(i)
	}
	challenge.SetUint64(m.state)
	m.updated = false
	return
}

func (m *MessageCounter) NextN(N int, i ...interface{}) (challenges []fr.Element) {
	challenges = make([]fr.Element, N)
	for n := 0; n < N; n++ {
		challenges[n] = m.Next(i)
		i = []interface{}{}
	}
	return
}

func NewMessageCounter(startState, step int) ArithmeticTranscript {
	transcript := &MessageCounter{state: uint64(startState), step: uint64(step)}
	transcript.Update([]byte{})
	return transcript
}

func NewMessageCounterGenerator(startState, step int) func() ArithmeticTranscript {
	return func() ArithmeticTranscript {
		return NewMessageCounter(startState, step)
	}
}
