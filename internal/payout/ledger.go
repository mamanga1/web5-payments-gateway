// ============================================================================
// internal/payout/ledger.go - Ledger Local Off-Chain
// ============================================================================
// Mantiene un registro contable local de saldos y operaciones.
// No depende de Internet. Solo para reconciliación y auditoría interna.
// ============================================================================

package payout

import (
	"encoding/json"
	"log"
	"os"
	"sync"
	"time"
)

// LedgerEntry representa una operación contable
type LedgerEntry struct {
	ID        string    `json:"id"`         // UUID o hash de la operación
	Type      string    `json:"type"`       // "TOLL_CREDIT", "TOLL_SETTLED", "PAYMENT_ROUTED"
	FromDID   string    `json:"from_did"`   // DID del nodo que origina (opcional)
	ToDID     string    `json:"to_did"`     // DID del nodo destino (opcional)
	Amount    float64   `json:"amount"`     // Monto en MAIA
	AssetID   string    `json:"asset_id"`   // "maia-liquidity-pool"
	Timestamp time.Time `json:"timestamp"`
}

// Ledger es el libro contable local (off-chain)
type Ledger struct {
	entries   []LedgerEntry
	balances  map[string]float64 // DID -> saldo actual
	mu        sync.RWMutex
	storagePath string
}

// NewLedger crea un nuevo ledger local
func NewLedger(storagePath string) *Ledger {
	l := &Ledger{
		entries:     make([]LedgerEntry, 0),
		balances:    make(map[string]float64),
		storagePath: storagePath,
	}
	l.loadFromDisk()
	return l
}

// AddEntry agrega una entrada al ledger y actualiza el saldo correspondiente
func (l *Ledger) AddEntry(entry LedgerEntry) {
	l.mu.Lock()
	defer l.mu.Unlock()

	// Agregar entrada
	l.entries = append(l.entries, entry)
	log.Printf("[LEDGER] Entrada registrada: %s | %.6f MAIA | %s", entry.Type, entry.Amount, entry.ToDID)

	// Actualizar saldo (si tiene ToDID)
	if entry.ToDID != "" {
		l.balances[entry.ToDID] += entry.Amount
		log.Printf("[LEDGER] Nuevo saldo para %s: %.6f MAIA", entry.ToDID, l.balances[entry.ToDID])
	}

	// Persistir a disco (opcional, cada N operaciones)
	if len(l.entries)%10 == 0 {
		l.saveToDisk()
	}
}

// GetBalance devuelve el saldo actual de un DID
func (l *Ledger) GetBalance(did string) float64 {
	l.mu.RLock()
	defer l.mu.RUnlock()
	return l.balances[did]
}

// GetEntries devuelve todas las entradas (para auditoría)
func (l *Ledger) GetEntries() []LedgerEntry {
	l.mu.RLock()
	defer l.mu.RUnlock()
	return append([]LedgerEntry{}, l.entries...)
}

// saveToDisk guarda el ledger en un archivo local (para no perder el estado)
func (l *Ledger) saveToDisk() {
	data, err := json.MarshalIndent(l.entries, "", "  ")
	if err != nil {
		log.Printf("[LEDGER] Error serializando ledger: %v", err)
		return
	}
	if err := os.WriteFile(l.storagePath, data, 0644); err != nil {
		log.Printf("[LEDGER] Error guardando ledger en disco: %v", err)
		return
	}
	log.Printf("[LEDGER] Ledger guardado en %s (%d entradas)", l.storagePath, len(l.entries))
}

// loadFromDisk carga el ledger desde un archivo local (si existe)
func (l *Ledger) loadFromDisk() {
	data, err := os.ReadFile(l.storagePath)
	if err != nil {
		if !os.IsNotExist(err) {
			log.Printf("[LEDGER] Error leyendo ledger desde disco: %v", err)
		}
		return
	}
	var entries []LedgerEntry
	if err := json.Unmarshal(data, &entries); err != nil {
		log.Printf("[LEDGER] Error deserializando ledger: %v", err)
		return
	}
	l.entries = entries
	// Reconstruir balances
	for _, e := range entries {
		if e.ToDID != "" {
			l.balances[e.ToDID] += e.Amount
		}
	}
	log.Printf("[LEDGER] Ledger cargado desde disco: %d entradas", len(l.entries))
}
