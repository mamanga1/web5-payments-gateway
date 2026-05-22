// ============================================================================
// internal/crypto/identity.go - DID + secp256k1 (adaptado de web5-mesh)
// ============================================================================

package crypto

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/decred/dcrd/dcrec/secp256k1/v4"
	"github.com/mr-tron/base58"
)

// DID estructura del identificador descentralizado
type DID struct {
	Method string
	Hash   []byte
}

func (d *DID) String() string {
	return d.Method + ":" + base58.Encode(d.Hash)
}

// Identity estructura con material criptográfico completo
type Identity struct {
	DID            *DID
	PrivateKey     *secp256k1.PrivateKey
	PublicKey      *secp256k1.PublicKey
	Name           string
	CreatedAt      time.Time
	LastSeen       time.Time
	Reputation     uint64
	SignatureCurve string
}

// NewIdentity genera una nueva identidad criptográfica con secp256k1
func NewIdentity(name string) (*Identity, error) {
	privKey, err := secp256k1.GeneratePrivateKey()
	if err != nil {
		return nil, fmt.Errorf("failed to generate private key: %w", err)
	}

	pubKey := privKey.PubKey()
	pubKeyCompressed := pubKey.SerializeCompressed()
	hash := sha256.Sum256(pubKeyCompressed)

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
		SignatureCurve: "secp256k1",
	}, nil
}

// Sign firma datos usando la clave privada secp256k1
func (id *Identity) Sign(data []byte) ([]byte, error) {
	if id.PrivateKey == nil {
		return nil, fmt.Errorf("no private key available")
	}
	hash := sha256.Sum256(data)
	signature, err := id.PrivateKey.Sign(hash[:])
	if err != nil {
		return nil, fmt.Errorf("signing failed: %w", err)
	}
	return signature.Serialize(), nil
}

// Verify verifica una firma usando la clave pública
func (id *Identity) Verify(data []byte, signature []byte) bool {
	if id.PublicKey == nil {
		return false
	}
	if len(signature) < 64 {
		return false
	}
	hash := sha256.Sum256(data)
	var sig secp256k1.Signature
	if err := sig.ParseDERSignature(signature); err != nil {
		// Si falla, intentar formato [R||S]
		var r, s [32]byte
		copy(r[:], signature[:32])
		copy(s[:], signature[32:64])
		sig.SetRS(r, s)
	}
	return sig.Verify(hash[:], id.PublicKey)
}

// GetDIDString retorna el DID como string legible
func (id *Identity) GetDIDString() string {
	return id.DID.String()
}

// GetPublicKeyHex retorna la clave pública en formato hexadecimal
func (id *Identity) GetPublicKeyHex() string {
	return hex.EncodeToString(id.PublicKey.SerializeCompressed())
}
