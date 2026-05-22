# Web5 Payments Gateway

**Sovereign Payment Gateway for IAON – AON‑CHP**

[![License](https://img.shields.io/badge/license-MIT%2BAnti--Corporate-blue)](LICENSE-TRINCHERA)
[![Go Version](https://img.shields.io/badge/go-1.23+-brightgreen)](https://golang.org/)
[![RAM Usage](https://img.shields.io/badge/ram-40~150MB-brightgreen)]()
[![Network](https://img.shields.io/badge/network-UDP%2FDHT-purple)]()

> **Decentralized, no fixed IP, no DNS. Each node is a cryptographic seat (DID) running on recycled hardware (TV boxes, Xeons, Poco F1).**

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
| **Zero‑IP / Zero‑DNS** | Peer location by XOR distance over Kademlia DHT; native UDP hole punching for CGNAT. |
| **Immutable Identity** | DID (`did:web5-mesh:secp256k1:<hash>`) – independent of physical IP. |
| **Self‑Sovereign** | Go 1.23 + secp256k1 + native DHT. No Web2 dependencies. |

---

## 🗂️ Repository Structure (planned)

```text
web5-payments-gateway/
├── cmd/
│   ├── gateway/
│   │   └── main.go           # Production entry point (UDP + DHT, no HTTP)
│   └── test-node/
│       └── main.go           # Validation node (includes --http-port for local testing)
├── internal/
│   ├── crypto/
│   │   └── identity.go       # DID + secp256k1 (adapted from web5-mesh)
│   ├── mesh/
│   │   └── transport.go      # UDP + DHT + packets + hole punching
│   └── payout/
│       ├── router.go         # 1% toll logic + event processing
│       ├── bridge.go         # MAIA network settlement bridge
│       └── ledger.go         # Local off‑chain reconciliation ledger
├── go.mod
└── README.md

🛫 Quick Start (Validation / Lab Mode)

⚠️ HTTP and /health are enabled ONLY in test-node (for local debugging).
Production (cmd/gateway) has HTTP disabled by default – zero public IP surface.

git clone https://github.com/mamanga1/web5-payments-gateway.git
cd web5-payments-gateway
go mod tidy
go build -o test-node ./cmd/test-node
./test-node --did "did:web5-mesh:test-01" --udp-port 9000 --http-port 8080 --run-test
curl http://localhost:8080/health   # → OK

⚠️ NOTICE
Protocol core owner: mamanga1 (IberaAON).
Negligent use or off‑chain ledger tampering by third parties does not exempt the operator from contractual responsibility to fund holders and mesh producers.

🧉 Credits
Built in the bunker of Corrientes, Argentina, on top of the web5-mesh IAON infrastructure.
Made with pride and endurance – without asking for permission.

