// ============================================================================
// cmd/gateway/main.go - Entry Point de Producción (UDP + DHT)
// ============================================================================
// Este es el binario que corre en el nodo productor (TV box, Xeon, Poco F1).
// No expone HTTP. Solo escucha UDP, aplica peaje del 1% y liquida en MAIA.
// ============================================================================

package main

import (
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/mamanga1/web5-payments-gateway/internal/crypto"
	"github.com/mamanga1/web5-payments-gateway/internal/mesh"
	"github.com/mamanga1/web5-payments-gateway/internal/payout"
)

func main() {
	// Flags de configuración
	nodeName := flag.String("name", "gateway-node", "Nombre del nodo")
	udpPort := flag.Int("udp-port", 4242, "Puerto UDP para escuchar eventos")
	tollRate := flag.Float64("toll-rate", 0.01, "Porcentaje de peaje (0.01 = 1%)")
	maiaEndpoint := flag.String("maia-endpoint", "", "Endpoint UDP de la red MAIA (ej: xeon.maia.net:4242)")
	flag.Parse()

	if *maiaEndpoint == "" {
		log.Fatal("[MAIN] Se requiere --maia-endpoint (dirección UDP de la red MAIA)")
	}

	log.Println("[MAIN] Iniciando Web5 Payments Gateway (modo producción - sin HTTP)")
	log.Printf("[MAIN] Configuración: puerto UDP=%d, peaje=%.0f%%, MAIA=%s", *udpPort, *tollRate*100, *maiaEndpoint)

	// 1. Generar identidad criptográfica
	identity, err := crypto.NewIdentity(*nodeName)
	if err != nil {
		log.Fatalf("[MAIN] Error generando identidad: %v", err)
	}
	log.Printf("[MAIN] Identidad generada: DID=%s", identity.GetDIDString())

	// 2. Crear el router (peaje + ledger)
	router := payout.NewRouter(identity.GetDIDString(), *tollRate)

	// 3. Crear el ledger local (persistencia)
	ledger := payout.NewLedger("ledger.json")

	// (Opcional) Conectar ledger con router para registrar operaciones
	_ = ledger

	// 4. Crear transporte UDP (escucha eventos de la malla)
	transport, err := mesh.NewTransport(*udpPort, identity.GetDIDString(), router)
	if err != nil {
		log.Fatalf("[MAIN] Error creando transporte: %v", err)
	}

	// 5. Iniciar listener UDP en una goroutine
	go transport.Start()

	// 6. Crear bridge MAIA para liquidar peajes
	bridge := payout.NewMAIABridge(*maiaEndpoint)

	// 7. Loop periódico para liquidar peajes acumulados
	go func() {
		ticker := time.NewTicker(1 * time.Minute)
		defer ticker.Stop()
		for range ticker.C {
			balance := router.GetBalance(identity.GetDIDString())
			if balance > 0 {
				log.Printf("[MAIN] Liquidando %.6f MAIA en red %s", balance, *maiaEndpoint)
				if err := bridge.SettleToll(identity.GetDIDString(), balance); err != nil {
					log.Printf("[MAIN] Error liquidando peaje: %v", err)
				}
			}
		}
	}()

	log.Println("[MAIN] Gateway operativo. Esperando eventos de pago...")

	// 8. Esperar señal de terminación (Ctrl+C)
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	<-sigCh

	log.Println("[MAIN] Apagando gateway...")
	transport.Stop()
	log.Println("[MAIN] Gateway detenido.")
}
