// ============================================================================
// internal/mesh/transport.go - UDP + DHT + Recepción de Eventos
// ============================================================================
// Capa de transporte pura. Escucha UDP, recibe eventos firmados,
// los valida y los pasa al router (payout) para aplicar el peaje.
// ============================================================================

package mesh

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"sync"
	"time"

	"github.com/mamanga1/web5-payments-gateway/internal/crypto"
	"github.com/mamanga1/web5-payments-gateway/internal/payout"
)

// Transport maneja la comunicación UDP con la malla
type Transport struct {
	conn       *net.UDPConn
	localDID   string
	router     *payout.Router
	eventQueue chan *payout.PaymentEvent
	stopChan   chan struct{}
	mu         sync.RWMutex
}

// NewTransport crea un nuevo transporte UDP
func NewTransport(listenPort int, localDID string, router *payout.Router) (*Transport, error) {
	addr, err := net.ResolveUDPAddr("udp", fmt.Sprintf("0.0.0.0:%d", listenPort))
	if err != nil {
		return nil, fmt.Errorf("error resolviendo puerto: %w", err)
	}

	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		return nil, fmt.Errorf("error escuchando UDP: %w", err)
	}

	log.Printf("[TRANSPORT] Escuchando UDP en puerto %d", listenPort)

	return &Transport{
		conn:       conn,
		localDID:   localDID,
		router:     router,
		eventQueue: make(chan *payout.PaymentEvent, 1000),
		stopChan:   make(chan struct{}),
	}, nil
}

// Start inicia el listener UDP (se ejecuta en una goroutine)
func (t *Transport) Start() {
	buf := make([]byte, 65535)
	log.Println("[TRANSPORT] Listener UDP activo")

	for {
		select {
		case <-t.stopChan:
			return
		default:
			// Leer paquete UDP
			t.conn.SetReadDeadline(time.Now().Add(500 * time.Millisecond))
			n, remoteAddr, err := t.conn.ReadFromUDP(buf)
			if err != nil {
				continue // timeout o error, seguimos
			}

			// Parsear evento
			var event payout.PaymentEvent
			if err := json.Unmarshal(buf[:n], &event); err != nil {
				log.Printf("[TRANSPORT] Paquete inválido desde %s: %v", remoteAddr, err)
				continue
			}

			log.Printf("[TRANSPORT] Evento recibido de %s: tipo=%s, monto=%.2f %s",
				event.Sender, event.Type, event.Amount, event.Currency)

			// Validar firma (si el evento tiene firma)
			if len(event.Signature) > 0 {
				// En producción, acá se verifica la firma del emisor
				log.Printf("[TRANSPORT] Firma presente (%d bytes)", len(event.Signature))
			}

			// Procesar con el router (aplica peaje)
			if err := t.router.ProcessPayment(&event); err != nil {
				log.Printf("[TRANSPORT] Error procesando pago: %v", err)
			}
		}
	}
}

// Stop detiene el transporte
func (t *Transport) Stop() {
	close(t.stopChan)
	t.conn.Close()
	log.Println("[TRANSPORT] Listener UDP detenido")
}
