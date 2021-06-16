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

package kzg

import (
	"errors"
	"io"
	"math/big"

	"github.com/consensys/gnark-crypto/ecc/bls12-381"
	"github.com/consensys/gnark-crypto/ecc/bls12-381/fr"
	"github.com/consensys/gnark-crypto/ecc/bls12-381/fr/fft"
	"github.com/consensys/gnark-crypto/ecc/bls12-381/fr/polynomial"
	fiatshamir "github.com/consensys/gnark-crypto/fiat-shamir"
	"github.com/consensys/gnark-crypto/internal/parallel"
)

var (
	errNbDigestsNeqNbPolynomials = errors.New("number of digests is not the same as the number of polynomials")
	errUnsupportedSize           = errors.New("the size of the polynomials exceeds the capacity of the SRS")
)

var (
	ErrVerifyOpeningProof            = errors.New("error verifying opening proof")
	ErrVerifyBatchOpeningSinglePoint = errors.New("error verifying batch opening proof at single point")
)

// Digest commitment of a polynomial.
type Digest struct {
	data bls12381.G1Affine
}

// Scheme stores KZG data
type Scheme struct {

	// Domain to perform polynomial division. The size of the domain is the lowest power of 2 greater than Size.
	Domain fft.Domain

	// SRS stores the result of the MPC
	SRS struct {
		G1 []bls12381.G1Affine  // [gen [alpha]gen , [alpha**2]gen, ... ]
		G2 [2]bls12381.G2Affine // [gen, [alpha]gen ]
	}
}

// Proof KZG proof for opening at a single point.
type Proof struct {

	// Point at which the polynomial is evaluated
	Point fr.Element

	// ClaimedValue purported value
	ClaimedValue fr.Element

	// H quotient polynomial (f - f(z))/(x-z)
	H bls12381.G1Affine
}

// NewScheme returns a new KZG scheme.
// This should be used for testing purpose only.
func NewScheme(size int, alpha fr.Element) *Scheme {

	s := &Scheme{}

	d := fft.NewDomain(uint64(size), 0, false)
	s.Domain = *d
	s.SRS.G1 = make([]bls12381.G1Affine, size)

	var bAlpha big.Int
	alpha.ToBigIntRegular(&bAlpha)

	_, _, gen1Aff, gen2Aff := bls12381.Generators()
	s.SRS.G1[0] = gen1Aff
	s.SRS.G2[0] = gen2Aff
	s.SRS.G2[1].ScalarMultiplication(&gen2Aff, &bAlpha)

	alphas := make([]fr.Element, size-1)
	alphas[0] = alpha
	for i := 1; i < len(alphas); i++ {
		alphas[i].Mul(&alphas[i-1], &alpha)
	}
	for i := 0; i < len(alphas); i++ {
		alphas[i].FromMont()
	}
	g1s := bls12381.BatchScalarMultiplicationG1(&gen1Aff, alphas)
	copy(s.SRS.G1[1:], g1s)

	return s
}

// Clone returns a copy of d
func (d *Digest) Clone() *Digest {
	var res Digest
	res.data.Set(&d.data)
	return &res
}

// Marshal serializes the point as in bls12381.G1Affine.
func (d *Digest) Marshal() []byte {
	return d.data.Marshal()
}

// Marshal serializes the point as in bls12381.G1Affine.
func (d *Digest) Unmarshal(buf []byte) error {
	err := d.data.Unmarshal(buf)
	if err != nil {
		return err
	}
	return nil
}

// Add adds two digest. The API and behaviour mimics bls12381.G1Affine's,
// i.e. the caller is modified.
func (d *Digest) Add(d1, d2 *Digest) *Digest {
	var p1, p2 bls12381.G1Jac
	p1.FromAffine(&d1.data)
	p2.FromAffine(&d2.data)
	p1.AddAssign(&p2)
	d.data.FromJacobian(&p1)
	return d
}

// Sub adds two digest. The API and behaviour mimics bls12381.G1Affine's,
// i.e. the caller is modified.
func (d *Digest) Sub(d1, d2 *Digest) *Digest {
	var p1, p2 bls12381.G1Jac
	p1.FromAffine(&d1.data)
	p2.FromAffine(&d2.data)
	p1.SubAssign(&p2)
	d.data.FromJacobian(&p1)
	return d
}

// Add adds two digest. The API and behaviour mimics bls12381.G1Affine's,
// i.e. the caller is modified.
func (d *Digest) ScalarMul(d1 *Digest, s big.Int) *Digest {
	var p1 bls12381.G1Affine
	p1.Set(&d1.data)
	p1.ScalarMultiplication(&p1, &s)
	d.data.Set(&p1)
	return d
}

// Marshal serializes a proof as H||point||claimed_value.
// The point H is not compressed.
func (p *Proof) Marshal() []byte {

	var res [4 * fr.Bytes]byte

	bH := p.H.RawBytes()
	copy(res[:], bH[:])
	be := p.Point.Bytes()
	copy(res[2*fr.Bytes:], be[:])
	be = p.ClaimedValue.Bytes()
	copy(res[3*fr.Bytes:], be[:])

	return res[:]
}

// GetClaimedValue returns the serialized claimed value.
func (p *Proof) GetClaimedValue() []byte {
	return p.ClaimedValue.Marshal()
}

type BatchProofsSinglePoint struct {
	// Point at which the polynomials are evaluated
	Point fr.Element

	// ClaimedValues purported values
	ClaimedValues []fr.Element

	// H quotient polynomial Sum_i gamma**i*(f - f(z))/(x-z)
	H bls12381.G1Affine
}

// Marshal serializes a proof as H||point||claimed_values.
// The point H is not compressed.
func (p *BatchProofsSinglePoint) Marshal() []byte {
	nbClaimedValues := len(p.ClaimedValues)

	// 2 for H, 1 for point, nbClaimedValues for the claimed values
	res := make([]byte, (3+nbClaimedValues)*fr.Bytes)

	bH := p.H.RawBytes()
	copy(res, bH[:])
	be := p.Point.Bytes()
	copy(res[2*fr.Bytes:], be[:])
	offset := 3 * fr.Bytes
	for i := 0; i < nbClaimedValues; i++ {
		be = p.ClaimedValues[i].Bytes()
		copy(res[offset:], be[:])
		offset += fr.Bytes
	}

	return res
}

// GetClaimedValues returns a slice of the claimed values,
// serialized.
func (p *BatchProofsSinglePoint) GetClaimedValues() [][]byte {
	res := make([][]byte, len(p.ClaimedValues))
	for i := 0; i < len(p.ClaimedValues); i++ {
		res[i] = p.ClaimedValues[i].Marshal()
	}
	return res
}

// WriteTo writes binary encoding of the scheme data.
// It writes only the SRS, the fft fomain is reconstructed
// from it.
func (s *Scheme) WriteTo(w io.Writer) (int64, error) {

	// encode the fft
	n, err := s.Domain.WriteTo(w)
	if err != nil {
		return n, err
	}

	// encode the SRS
	enc := bls12381.NewEncoder(w)

	toEncode := []interface{}{
		&s.SRS.G2[0],
		&s.SRS.G2[1],
		s.SRS.G1,
	}

	for _, v := range toEncode {
		if err := enc.Encode(v); err != nil {
			return n + enc.BytesWritten(), err
		}
	}

	return n + enc.BytesWritten(), nil
}

// ReadFrom decodes KZG data from reader.
// The kzg data should have been encoded using WriteTo.
// Only the points from the SRS are actually encoded in the
// reader, the fft domain is reconstructed from it.
func (s *Scheme) ReadFrom(r io.Reader) (int64, error) {

	// decode the fft
	n, err := s.Domain.ReadFrom(r)
	if err != nil {
		return n, err
	}

	// decode the SRS
	dec := bls12381.NewDecoder(r)

	toDecode := []interface{}{
		&s.SRS.G2[0],
		&s.SRS.G2[1],
		&s.SRS.G1,
	}

	for _, v := range toDecode {
		if err := dec.Decode(v); err != nil {
			return n + dec.BytesRead(), err
		}
	}

	return n + dec.BytesRead(), nil

}

// Commit commits to a polynomial using a multi exponentiation with the SRS.
// It is assumed that the polynomial is in canonical form, in Montgomery form.
func (s *Scheme) Commit(p polynomial.Polynomial) (Digest, error) {

	if p.Degree() >= s.Domain.Cardinality {
		return Digest{}, errUnsupportedSize
	}

	var _res bls12381.G1Affine

	// ensure we don't modify p
	pCopy := make(polynomial.Polynomial, s.Domain.Cardinality)
	copy(pCopy, p)

	parallel.Execute(len(p), func(start, end int) {
		for i := start; i < end; i++ {
			pCopy[i].FromMont()
		}
	})
	_res.MultiExp(s.SRS.G1, pCopy)

	var res Digest
	res.data.Set(&_res)

	return res, nil
}

// Open computes an opening proof of p at _val.
// Returns a MockProof, which is an empty interface.
func (s *Scheme) Open(point *fr.Element, p polynomial.Polynomial) (Proof, error) {

	if p.Degree() >= s.Domain.Cardinality {
		panic("[Open] Size of polynomial exceeds the size supported by the scheme")
	}

	// build the proof
	res := Proof{
		Point:        *point,
		ClaimedValue: p.Eval(point),
	}

	// compute H

	h := dividePolyByXminusA(s.Domain, p, res.ClaimedValue, res.Point)

	// commit to H
	c, err := s.Commit(h)
	if err != nil {
		return Proof{}, err
	}
	res.H.Set(&c.data)

	return res, nil
}

// Verify verifies a KZG opening proof at a single point
func (s *Scheme) Verify(d *Digest, proof *Proof) error {

	var _commitment bls12381.G1Affine
	_commitment.Set(&d.data)

	// comm(f(a))
	var claimedValueG1Aff bls12381.G1Affine
	var claimedValueBigInt big.Int
	proof.ClaimedValue.ToBigIntRegular(&claimedValueBigInt)
	claimedValueG1Aff.ScalarMultiplication(&s.SRS.G1[0], &claimedValueBigInt)

	// [f(alpha) - f(a)]G1Jac
	var fminusfaG1Jac, tmpG1Jac bls12381.G1Jac
	fminusfaG1Jac.FromAffine(&_commitment)
	tmpG1Jac.FromAffine(&claimedValueG1Aff)
	fminusfaG1Jac.SubAssign(&tmpG1Jac)

	// [-H(alpha)]G1Aff
	var negH bls12381.G1Affine
	negH.Neg(&proof.H)

	// [alpha-a]G2Jac
	var alphaMinusaG2Jac, genG2Jac, alphaG2Jac bls12381.G2Jac
	var pointBigInt big.Int
	proof.Point.ToBigIntRegular(&pointBigInt)
	genG2Jac.FromAffine(&s.SRS.G2[0])
	alphaG2Jac.FromAffine(&s.SRS.G2[1])
	alphaMinusaG2Jac.ScalarMultiplication(&genG2Jac, &pointBigInt).
		Neg(&alphaMinusaG2Jac).
		AddAssign(&alphaG2Jac)

	// [alpha-a]G2Aff
	var xminusaG2Aff bls12381.G2Affine
	xminusaG2Aff.FromJacobian(&alphaMinusaG2Jac)

	// [f(alpha) - f(a)]G1Aff
	var fminusfaG1Aff bls12381.G1Affine
	fminusfaG1Aff.FromJacobian(&fminusfaG1Jac)

	// e([-H(alpha)]G1Aff, G2gen).e([-H(alpha)]G1Aff, [alpha-a]G2Aff) ==? 1
	check, err := bls12381.PairingCheck(
		[]bls12381.G1Affine{fminusfaG1Aff, negH},
		[]bls12381.G2Affine{s.SRS.G2[0], xminusaG2Aff},
	)
	if err != nil {
		return err
	}
	if !check {
		return ErrVerifyOpeningProof
	}
	return nil
}

// BatchOpenSinglePoint creates a batch opening proof at _val of a list of polynomials.
// It's an interactive protocol, made non interactive using Fiat Shamir.
// point is the point at which the polynomials are opened.
// digests is the list of committed polynomials to open, need to derive the challenge using Fiat Shamir.
// polynomials is the list of polynomials to open.
func (s *Scheme) BatchOpenSinglePoint(point *fr.Element, digests []Digest, polynomials []polynomial.Polynomial) (BatchProofsSinglePoint, error) {

	nbDigests := len(digests)
	if nbDigests != len(polynomials) {
		return BatchProofsSinglePoint{}, errNbDigestsNeqNbPolynomials
	}

	var res BatchProofsSinglePoint

	// compute the purported values
	res.ClaimedValues = make([]fr.Element, len(polynomials))
	for i := 0; i < len(polynomials); i++ {
		res.ClaimedValues[i] = polynomials[i].Eval(point)
	}

	// set the point at which the evaluation is done
	res.Point.Set(point)

	// derive the challenge gamma, binded to the point and the commitments
	fs := fiatshamir.NewTranscript(fiatshamir.SHA256, "gamma")
	if err := fs.Bind("gamma", res.Point.Marshal()); err != nil {
		return BatchProofsSinglePoint{}, err
	}
	for i := 0; i < len(digests); i++ {
		if err := fs.Bind("gamma", digests[i].Marshal()); err != nil {
			return BatchProofsSinglePoint{}, err
		}
	}
	gammaByte, err := fs.ComputeChallenge("gamma")
	if err != nil {
		return BatchProofsSinglePoint{}, err
	}
	var gamma fr.Element
	gamma.SetBytes(gammaByte)

	// compute sum_i gamma**i*f and sum_i gamma**i*f(a)
	var sumGammaiTimesEval fr.Element
	sumGammaiTimesEval.Set(&res.ClaimedValues[nbDigests-1])
	sumGammaiTimesPol := polynomials[nbDigests-1].Clone()
	for i := nbDigests - 2; i >= 0; i-- {
		sumGammaiTimesEval.Mul(&sumGammaiTimesEval, &gamma).
			Add(&sumGammaiTimesEval, &res.ClaimedValues[i])
		sumGammaiTimesPol.ScaleInPlace(&gamma)
		sumGammaiTimesPol.Add(polynomials[i], sumGammaiTimesPol)
	}

	// compute H
	h := dividePolyByXminusA(s.Domain, sumGammaiTimesPol, sumGammaiTimesEval, res.Point)
	c, err := s.Commit(h)
	if err != nil {
		return BatchProofsSinglePoint{}, err
	}

	res.H.Set(&c.data)

	return res, nil
}

// BatchVerifySinglePoint verifies a batched opening proof at a single point of a list of polynomials.
// point: point at which the polynomials are evaluated
// claimedValues: claimed values of the polynomials at _val
// commitments: list of commitments to the polynomials which are opened
// batchOpeningProof: the batched opening proof at a single point of the polynomials.
func (s *Scheme) BatchVerifySinglePoint(digests []Digest, batchOpeningProof *BatchProofsSinglePoint) error {

	nbDigests := len(digests)

	// check consistancy between numbers of claims vs number of digests
	if len(digests) != len(batchOpeningProof.ClaimedValues) {
		return errNbDigestsNeqNbPolynomials
	}

	// derive the challenge gamma, binded to the point and the commitments
	fs := fiatshamir.NewTranscript(fiatshamir.SHA256, "gamma")
	err := fs.Bind("gamma", batchOpeningProof.Point.Marshal())
	if err != nil {
		return err
	}
	for i := 0; i < len(digests); i++ {
		err := fs.Bind("gamma", digests[i].Marshal())
		if err != nil {
			return err
		}
	}
	gammaByte, err := fs.ComputeChallenge("gamma")
	if err != nil {
		return err
	}
	var gamma fr.Element
	gamma.SetBytes(gammaByte)

	var sumGammaiTimesEval fr.Element
	sumGammaiTimesEval.Set(&batchOpeningProof.ClaimedValues[nbDigests-1])
	for i := nbDigests - 2; i >= 0; i-- {
		sumGammaiTimesEval.Mul(&sumGammaiTimesEval, &gamma).
			Add(&sumGammaiTimesEval, &batchOpeningProof.ClaimedValues[i])
	}

	var sumGammaiTimesEvalBigInt big.Int
	sumGammaiTimesEval.ToBigIntRegular(&sumGammaiTimesEvalBigInt)
	var sumGammaiTimesEvalG1Aff bls12381.G1Affine
	sumGammaiTimesEvalG1Aff.ScalarMultiplication(&s.SRS.G1[0], &sumGammaiTimesEvalBigInt)

	var acc fr.Element
	acc.SetOne()
	gammai := make([]fr.Element, len(digests))
	gammai[0].SetOne().FromMont()
	for i := 1; i < len(digests); i++ {
		acc.Mul(&acc, &gamma)
		gammai[i].Set(&acc).FromMont()
	}
	var sumGammaiTimesDigestsG1Aff bls12381.G1Affine
	_digests := make([]bls12381.G1Affine, len(digests))
	for i := 0; i < len(digests); i++ {
		_digests[i].Set(&digests[i].data)
	}

	sumGammaiTimesDigestsG1Aff.MultiExp(_digests, gammai)

	// sum_i [gamma**i * (f-f(a))]G1
	var sumGammiDiffG1Aff bls12381.G1Affine
	var t1, t2 bls12381.G1Jac
	t1.FromAffine(&sumGammaiTimesDigestsG1Aff)
	t2.FromAffine(&sumGammaiTimesEvalG1Aff)
	t1.SubAssign(&t2)
	sumGammiDiffG1Aff.FromJacobian(&t1)

	// [alpha-a]G2Jac
	var alphaMinusaG2Jac, genG2Jac, alphaG2Jac bls12381.G2Jac
	var pointBigInt big.Int
	batchOpeningProof.Point.ToBigIntRegular(&pointBigInt)
	genG2Jac.FromAffine(&s.SRS.G2[0])
	alphaG2Jac.FromAffine(&s.SRS.G2[1])
	alphaMinusaG2Jac.ScalarMultiplication(&genG2Jac, &pointBigInt).
		Neg(&alphaMinusaG2Jac).
		AddAssign(&alphaG2Jac)

	// [alpha-a]G2Aff
	var xminusaG2Aff bls12381.G2Affine
	xminusaG2Aff.FromJacobian(&alphaMinusaG2Jac)

	// [-H(alpha)]G1Aff
	var negH bls12381.G1Affine
	negH.Neg(&batchOpeningProof.H)

	// check the pairing equation
	check, err := bls12381.PairingCheck(
		[]bls12381.G1Affine{sumGammiDiffG1Aff, negH},
		[]bls12381.G2Affine{s.SRS.G2[0], xminusaG2Aff},
	)
	if err != nil {
		return err
	}
	if !check {
		return ErrVerifyBatchOpeningSinglePoint
	}
	return nil
}
