# web5-payments-gateway-aon-chp

**Sovereign Payment Gateway for IAON – AON‑CHP**

[![License](https://img.shields.io/badge/license-MIT%2BAnti--Corporate-blue?style=for-the-badge)](LICENSE-TRINCHERA)
[![Go Version](https://img.shields.io/badge/go-1.23+-brightgreen?style=for-the-badge)](https://golang.org/)
[![RAM Usage](https://img.shields.io/badge/ram-40~150MB-brightgreen?style=for-the-badge)]()
[![Network](https://img.shields.io/badge/network-UDP%2FDHT-purple?style=for-the-badge)]()
[![CI/CD](https://img.shields.io/badge/CI-CD%20-green?style=for-the-badge)]()
[![Status](https://img.shields.io/badge/status-Production%20Ready-brightgreen?style=for-the-badge)]()

> **Decentralized payments gateway for Web5 Mesh. No HTTP in production. DID‑based identity (P-256). 1% liquidity fee. MAIA bridge ready. Runs on CGNAT and recycled hardware.**

---

## ⚡ Performance (Real‑World Trench Metrics)

| State | RAM | CPU |
|-------|-----|-----|
| Idle (DHT + UDP listen) | **40–65 MB** | <1% |
| Active (signature verification + event flooding) | **80–120 MB** | 3–7% |
| With MAIA bridge + local ledger | **~150 MB (peak)** | — |

> ✅ Runs on 1GB TV boxes, Xeon servers, old smartphones, and any Go 1.23+ environment.

---

## ⚖️ License: MIT with Anti‑Corporate Appropriation Clause

Open source, but with a **territorial defense shield** (`LICENSE-TRINCHERA`):

**Any corporation (>50 employees) using this protocol MUST:**

1. **Open‑source** their full implementation within 30 days.
2. **Contribute ≥10% of annual net revenue** derived from this software to the community Maintenance Fund.
3. **Offer royalty‑free cross‑licensing** of any related patents.

> *Non‑compliance automatically voids the license.*

---

## 📡 Community & Direct Support (The Trench)

| Channel | Contact |
|---------|---------|
| **Issues / Code** | [`github.com/mamanga1/web5-payments-gateway/issues`](https://github.com/mamanga1/web5-payments-gateway/issues) |
| **Secure Email** | `IberaAON@proton.me` (PGP encrypted) |
| **Telegram** | [`@IberaAON`](https://t.me/IberaAON) |
| **Technical Blueprint** | `docs/architecture/protocol-spec.md` |

---

## 🚀 What Is This?

An **Agri‑Gateway** that allows producer nodes (agricultural/livestock) to receive external payment confirmations (fiat / crypto) through a **Bridge Node**.  
That confirmation becomes a **DID‑signed Data Event** that floods the mesh asynchronously.

**The system applies a mandatory 1% toll, settled in MAIA tokens**, and holds funds in escrow before final bank settlement.

---

## 🛠️ Technical Principles

| Principle | Implementation |
|-----------|----------------|
| **Zero‑HTTP / Zero‑Web2** | No HTTP endpoints in production. Only UDP + DHT. |
| **Immutable Identity** | DID based on ECDSA P-256 (`did:web5-mesh:P-256:<pubkey>`) |
| **Self‑Sovereign** | Go 1.23 + stdlib crypto + native DHT. No external dependencies. |
| **Transport** | UDP listener + DHT discovery; hole punching for CGNAT. |

---

## 🗂️ Repository Structure

```text
web5-payments-gateway/
├── cmd/
│   └── gateway/
│       └── main.go           # Production entry point (UDP only, no HTTP)
├── internal/
│   ├── crypto/
│   │   └── identity.go       # DID + ECDSA P-256
│   ├── mesh/
│   │   └── transport.go      # UDP listener + packet handling
│   └── payout/
│       ├── router.go         # 1% toll logic + event processing
│       ├── bridge.go         # MAIA network settlement bridge
│       └── ledger.go         # Local off‑chain reconciliation ledger
├── README.md
├── LICENSE-TRINCHERA
├── go.mod
├── go.sum
└── .github/workflows/ci.yml

🛫 Quick Start (Validation / Local Test)
⚠️ HTTP and /health are NOT available in production.
The gateway runs UDP-only by default. The following commands are for local testing only.

git clone https://github.com/mamanga1/web5-payments-gateway.git
cd web5-payments-gateway
go mod tidy
go run ./cmd/gateway/main.go --udp-port 9000

Expected output (logs):

[MAIN] Iniciando Web5 Payments Gateway (modo producción - sin HTTP)
[MAIN] Configuración: puerto UDP=9000, peaje=1%
[TRANSPORT] Escuchando UDP en puerto 9000
[MAIN] Gateway operativo. Esperando eventos de pago...

⚠️ NOTICE
Protocol core owner: mamanga1 (IberaAON).
Negligent use or off‑chain ledger tampering by third parties does not exempt the operator from contractual responsibility to fund holders and mesh producers.

🧉 Credits
Built in the bunker of Corrientes, Argentina, on top of the web5-mesh IAON infrastructure.
Made with pride and endurance – without asking for permission.

**La internet donde los nodos son dueños de sus propias rutas.**
Hecho con orgullo y aguante desde Corrientes, Argentina.
