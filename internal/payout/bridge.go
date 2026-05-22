// ============================================================================
// internal/payout/bridge.go - Conexión con la Red MAIA
// ============================================================================
// Este puente conecta el gateway con tu red MAIA (la Xeon del búnker).
// Liquida el peaje acumulado en tokens MAIA de forma asíncrona.
// ============================================================================

package payout

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"time"
)

// MAIABridge conecta el gateway con la red MAIA (Xeon del búnker)
type MAIABridge struct {
	endpoint string        // Dirección UDP del nodo MAIA (ej: "xeon.maia.net:4242")
	timeout  time.Duration
}

// NewMAIABridge crea un nuevo puente hacia la red MAIA
func NewMAIABridge(endpoint string) *MAIABridge {
	return &MAIABridge{
		endpoint: endpoint,
		timeout:  5 * time.Second,
	}
}

// TollSettlement representa una liquidación de peaje a enviar a la red MAIA
type TollSettlement struct {
	Type      string    `json:"type"`       // "TOLL_SETTLEMENT"
	FromNode  string    `json:"from_node"`  // DID del nodo gateway
	Amount    float64   `json:"amount"`     // Cantidad de MAIA a liquidar
	AssetID   string    `json:"asset_id"`   // "maia-liquidity-pool"
	Timestamp time.Time `json:"timestamp"`
}

// SettleToll envía una liquidación de peaje a la red MAIA por UDP
func (b *MAIABridge) SettleToll(fromNodeDID string, amount float64) error {
	if amount <= 0 {
		return fmt.Errorf("monto inválido para liquidar: %.6f", amount)
	}

	settlement := TollSettlement{
		Type:      "TOLL_SETTLEMENT",
		FromNode:  fromNodeDID,
		Amount:    amount,
		AssetID:   "maia-liquidity-pool",
		Timestamp: time.Now(),
	}

	data, err := json.Marshal(settlement)
	if err != nil {
		return fmt.Errorf("error serializando liquidación: %w", err)
	}

	log.Printf("[BRIDGE] Enviando liquidación de %.6f MAIA a %s", amount, b.endpoint)

	// Resolver dirección UDP
	udpAddr, err := net.ResolveUDPAddr("udp", b.endpoint)
	if err != nil {
		return fmt.Errorf("error resolviendo endpoint MAIA: %w", err)
	}

	// Conectar y enviar
	conn, err := net.DialUDP("udp", nil, udpAddr)
	if err != nil {
		return fmt.Errorf("error conectando con red MAIA: %w", err)
	}
	defer conn.Close()

	// Setear timeout
	if err := conn.SetWriteDeadline(time.Now().Add(b.timeout)); err != nil {
		return fmt.Errorf("error seteando timeout: %w", err)
	}

	_, err = conn.Write(data)
	if err != nil {
		return fmt.Errorf("error enviando liquidación a MAIA: %w", err)
	}

	log.Printf("[BRIDGE] Liquidación enviada con éxito: %.6f MAIA", amount)
	return nil
}

// SettleMultipleTolls envía múltiples liquidaciones en un solo lote (optimización)
func (b *MAIABridge) SettleMultipleTolls(settlements []TollSettlement) error {
	if len(settlements) == 0 {
		return nil
	}

	data, err := json.Marshal(settlements)
	if err != nil {
		return fmt.Errorf("error serializando lote: %w", err)
	}

	udpAddr, err := net.ResolveUDPAddr("udp", b.endpoint)
	if err != nil {
		return fmt.Errorf("error resolviendo endpoint MAIA: %w", err)
	}

	conn, err := net.DialUDP("udp", nil, udpAddr)
	if err != nil {
		return fmt.Errorf("error conectando con red MAIA: %w", err)
	}
	defer conn.Close()

	if err := conn.SetWriteDeadline(time.Now().Add(b.timeout)); err != nil {
		return fmt.Errorf("error seteando timeout: %w", err)
	}

	_, err = conn.Write(data)
	if err != nil {
		return fmt.Errorf("error enviando lote a MAIA: %w", err)
	}

	log.Printf("[BRIDGE] Lote de %d liquidaciones enviado con éxito", len(settlements))
	return nil
}
