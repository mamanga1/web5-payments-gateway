# web5-payments-gateway-aon-chp

**UDP + DHT payments gateway for Web5 Mesh.**  
Zero HTTP by default, DID-based identity (P-256), 1% liquidity fee, MAIA bridge ready.  
Designed for field nodes running on CGNAT and weak connectivity.

---

## ⚖️ License: MIT with Anti-Corporate Appropriation Clause

Base: MIT License.

**Anti-appropriation clause:**

Any corporation (>50 employees) using this protocol must:

- open-source their implementation within 30 days;
- contribute ≥10% of net revenue to the maintenance fund;
- offer patent cross-licensing where applicable.

See `LICENSE-TRINCHERA` for full text.

---

## 📞 Community & Direct Support

- **Issues and code:** [`github.com/mamanga1/web5-payments-gateway/issues`](https://github.com/mamanga1/web5-payments-gateway/issues)
- **Bunker mail (secure, E2E):** `IberaAON@proton.me`
- **Telegram trinchera:** [@IberaAON](https://t.me/IberaAON)
- **Protocol spec:** `docs/architecture/protocol-spec.md`

---

## 🚀 What it is

Minimal payments gateway that:

- listens UDP on the Web5 Mesh DHT;
- applies a **1% liquidity fee** per validated payment event;
- offloads settlement to your MAIA network via `bridge.go`.

---

## 🛠️ Technical principles

- **No Web2 dependencies** by default.
- **Identity:** DID based on ECDSA P-256 (`did:web5-mesh:P-256:<pubkey>`).
- **Transport:** UDP + DHT; no public IPs required; discovery via mesh layer.
- **Code layout:** clean `cmd/`/`internal/` split, strict CI with `go.sum`.

---

## 🗂️ Repo structure

web5-payments-gateway/
├── cmd/gateway/main.go # entry point (UDP-only, production-ready)
├── internal/
│ ├── crypto/identity.go # DID parsing and P-256 ECDSA keys
│ ├── mesh/transport.go # UDP listener, packet framing, DHT hooks
│ └── payout/
│ ├── router.go # 1% fee logic + local routing rules
│ ├── bridge.go # settlement bridge to MAIA (UDP/RPC-ready)
│ └── ledger.go # off-chain ledger for saldos and reconciliation
├── README.md
├── LICENSE-TRINCHERA
├── go.mod
└── go.sum


---

## 🛫 Quick start (local test)

```bash
git clone https://github.com/mamanga1/web5-payments-gateway.git
cd web5-payments-gateway
go mod tidy
go run ./cmd/gateway/main.go --udp-port 9000

⚠️ HTTP is disabled by default in production.
Only UDP listeners and internal payout logic are active.

Logs will show: PAYMENT_CONFIRMED, fee application, MAIA settlement attempts.

NOTICE
Owner of the core protocol: mamanga1 (IberaAON).
Negligent use of this implementation does not exempt third parties from contractual liability with fund holders.

