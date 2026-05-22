// ============================================================================
// internal/payout/router.go - Peaje 1% + Lógica de Pagos
// ============================================================================
// Este es el corazón del gateway. Procesa eventos de pago, aplica el peaje
// obligatorio del 1% en tokens MAIA, y enruta los fondos.
// ============================================================================

package payout

import (
	"encoding/json"
	"fmt"
	"log"
	"sync"
	"time"
)

// PaymentEvent representa un pago confirmado desde un Bridge Node externo
type PaymentEvent struct {
	Type      string  `json:"type"`       // "PAYMENT_CONFIRMED"
	Sender    string  `json:"sender"`     // DID del Bridge emisor
	Receiver  string  `json:"receiver"`   // DID del Nodo Productor (el que cobra)
	Amount    float64 `json:"amount"`     // Monto en moneda original (USD, EUR, ARS)
	Currency  string  `json:"currency"`   // "USD", "EUR", "ARS", etc.
	AssetID   string  `json:"asset_id"`   // Identificador del lote/animal/producto
	Timestamp int64   `json:"timestamp"`
	Signature []byte  `json:"signature"`  // Firma criptográfica del Bridge
}

// TollConfig configuración del peaje
type TollConfig struct {
	Rate      float64 // Porcentaje del peaje (ej: 0.01 = 1%)
	TokenName string  // "MAIA"
	AssetID   string  // ID del pool de liquidez MAIA
}

// Router es el motor principal de pagos y peajes
type Router struct {
	config      TollConfig
	localLedger map[string]float64 // Saldo local indexado por DID (off-chain)
	mu          sync.RWMutex
	nodeDID     string
}

// NewRouter crea un nuevo router con la configuración de peaje
func NewRouter(nodeDID string, tollRate float64) *Router {
	return &Router{
		config: TollConfig{
			Rate:      tollRate,
			TokenName: "MAIA",
			AssetID:   "maia-liquidity-pool",
		},
		localLedger: make(map[string]float64),
		nodeDID:     nodeDID,
	}
}

// ProcessPayment procesa un evento de pago, aplica el peaje y actualiza el ledger local
func (r *Router) ProcessPayment(event *PaymentEvent) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	log.Printf("[ROUTER] Procesando pago de %.2f %s para productor %s", event.Amount, event.Currency, event.Receiver)

	// 1. Calcular el peaje (1% por defecto)
	tollAmount := event.Amount * r.config.Rate
	if tollAmount <= 0 {
		return fmt.Errorf("peaje inválido: %.6f (monto original: %.2f)", tollAmount, event.Amount)
	}

	log.Printf("[ROUTER] Peaje del %.0f%%: %.6f %s en tokens %s", r.config.Rate*100, tollAmount, event.Currency, r.config.TokenName)

	// 2. Acreditar el peaje en el ledger local del nodo operador
	r.localLedger[r.nodeDID] += tollAmount
	log.Printf("[ROUTER] Peaje acreditado en ledger local. Saldo operador: %.6f MAIA", r.localLedger[r.nodeDID])

	// 3. Registrar el evento para el productor (pendiente de liquidación bancaria)
	log.Printf("[ROUTER] Pago de %.2f %s registrado para productor %s. Asset ID: %s", 
		event.Amount-tollAmount, event.Currency, event.Receiver, event.AssetID)

	return nil
}

// GetBalance devuelve el saldo actual del ledger local para un DID
func (r *Router) GetBalance(did string) float64 {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.localLedger[did]
}

// SettleWithBridge envía el peaje acumulado a la red MAIA (se llama periódicamente)
func (r *Router) SettleWithBridge(bridgeEndpoint string) error {
	r.mu.Lock()
	balance := r.localLedger[r.nodeDID]
	if balance <= 0 {
		r.mu.Unlock()
		log.Printf("[ROUTER] Sin saldo para liquidar en MAIA")
		return nil
	}
	// Limpiar ledger local antes de enviar (para evitar doble liquidación)
	r.localLedger[r.nodeDID] = 0
	r.mu.Unlock()

	log.Printf("[ROUTER] Liquidando %.6f MAIA en bridge %s", balance, bridgeEndpoint)
	// Aquí va la llamada real al bridge (lo vamos a implementar en bridge.go)

	return nil
}
