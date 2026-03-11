# Technische Dokumentation: topoclean Architektur (v1.3)

## 1. Systemübersicht
`topoclean` ist modular nach SOLID-Prinzipien aufgebaut. Die Kernlogik nutzt das Strategy-Pattern für maximale Flexibilität.

## 2. Die Module

### 2.1. Ledger (internal/ledger)
- **SQLite-Backend:** CGO-freie Persistenz der Transaktionen.
- **Traceability:** Speichert Source, Dest, Hash und Size jeder Operation.
- **Status-Management:** `Pending` -> `Committed` | `RolledBack`.

### 2.2. Scanner (internal/scanner)
- **Magic-Byte Analyse:** Nutzt `http.DetectContentType` (512 Bytes) für MIME-Ermittlung.
- **No-Fly-Zone:** Schützt Dotfiles und das Ledger-System.

### 2.3. Vector (internal/vector) - Strategie-Pipeline
Nutzt eine Kette von spezialisierten Strategien (`Strategy` Interface):
1. **MIME-Strategy:** Priorisierte Inhaltsanalyse.
2. **Extension-Strategy:** Robuster Fallback für bekannte Endungen (v1.3 erweitert).
3. **Substring-Strategy:** Semantische Analyse von Dateinamen (z.B. `vault`, `inkasso`).

### 2.4. App Orchestrator (internal/app)
- **Deduplizierung:** Implementiert via Hardlinks (`os.Link`). Prüft Sitzungs-Cache und Ledger auf identische Hashes.
- **Soterische Operation:** `Hash -> Copy -> Verify -> Log -> Delete`.
- **Inversion:** Ermöglicht bit-genauen Rollback aller Operationen einer Transaktion.

## 3. Sicherheits-Invarianten
- **Atomic Moves:** Dateien werden erst gelöscht, wenn das Ziel kryptographisch verifiziert wurde.
- **Interactive-Confirmation:** Physische Änderungen erfordern menschliche Validierung.

---
*Status: v1.3 Dokumentation synchronisiert mit Codebase.*
