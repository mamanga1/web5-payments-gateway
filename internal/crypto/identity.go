// ============================================================================
// internal/crypto/identity.go - DID + ECDSA (estándar de Go)
// ============================================================================

package crypto

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"math/big"
	"time"

	"github.com/mr-tron/base58"
)

type DID struct {
	Method string
	Hash   []byte
}

func (d *DID) String() string {
	return d.Method + ":" + base58.Encode(d.Hash)
}

type Identity struct {
	DID            *DID
	PrivateKey     *ecdsa.PrivateKey
	PublicKey      *ecdsa.PublicKey
	Name           string
	CreatedAt      time.Time
	LastSeen       time.Time
	Reputation     uint64
	SignatureCurve string
}

func NewIdentity(name string) (*Identity, error) {
	privKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return nil, fmt.Errorf("failed to generate private key: %w", err)
	}
	pubKey := &privKey.PublicKey
	pubKeyBytes := elliptic.Marshal(elliptic.P256(), pubKey.X, pubKey.Y)
	hash := sha256.Sum256(pubKeyBytes)
	did := &DID{
		Method: "did:web5-mesh",
		Hash:   hash[:],
	}
	now := time.Now()
	return &Identity{
		DID:            did,
		PrivateKey:     privKey,
		PublicKey:      pubKey,
		Name:           name,
		CreatedAt:      now,
		LastSeen:       now,
		Reputation:     100,
		SignatureCurve: "P-256",
	}, nil
}

func (id *Identity) Sign(data []byte) ([]byte, error) {
	if id.PrivateKey == nil {
		return nil, fmt.Errorf("no private key available")
	}
	hash := sha256.Sum256(data)
	r, s, err := ecdsa.Sign(rand.Reader, id.PrivateKey, hash[:])
	if err != nil {
		return nil, fmt.Errorf("signing failed: %w", err)
	}
	signature := append(r.Bytes(), s.Bytes()...)
	return signature, nil
}

func (id *Identity) Verify(data []byte, signature []byte) bool {
	if id.PublicKey == nil {
		return false
	}
	if len(signature) < 64 {
		return false
	}
	hash := sha256.Sum256(data)
	r := new(big.Int).SetBytes(signature[:32])
	s := new(big.Int).SetBytes(signature[32:64])
	return ecdsa.Verify(id.PublicKey, hash[:], r, s)
}

func (id *Identity) GetDIDString() string {
	return id.DID.String()
}

func (id *Identity) GetPublicKeyHex() string {
	pubKeyBytes := elliptic.Marshal(elliptic.P256(), id.PublicKey.X, id.PublicKey.Y)
	return hex.EncodeToString(pubKeyBytes)
}
