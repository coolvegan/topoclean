# Technische Dokumentation: topoclean Architektur (v1.4)

## 1. Systemübersicht
`topoclean` ist modular nach SOLID-Prinzipien aufgebaut. Die Kernlogik nutzt das Strategy-Pattern für maximale Flexibilität und ein zentrales Konfigurations-Modul für Multi-Zonen-Management.

## 2. Die Module

### 2.1. Config (internal/config) - Neu in v1.4
- **Zuständigkeit:** Lädt die `config.json`, validiert Zonen und löst Pfade (Tilde-Expansion) auf.
- **Biographische Regeln:** Definiert, ob die Herkunft (`From-<Zone>`) gewahrt werden soll.

### 2.2. Ledger (internal/ledger)
- **SQLite-Backend:** CGO-freie Persistenz der Transaktionen.
- **Traceability:** Speichert Source, Dest, Hash und Size jeder Operation.

### 2.3. Scanner (internal/scanner)
- **Multi-Zone-fähig:** Scans können nun mit einem `ZoneName` assoziiert werden.
- **Magic-Byte Analyse:** Nutzt `http.DetectContentType` (512 Bytes).

### 2.4. Vector (internal/vector) - Strategie-Pipeline
- **MIME-Strategy:** Priorisierte Inhaltsanalyse.
- **Extension-Strategy:** Robuster Fallback (v1.3).
- **Substring-Strategy:** Semantische Analyse (v1.2).

### 2.5. App Orchestrator (internal/app)
- **Multi-Zone Planer:** Iteriert über alle konfigurierten Zonen.
- **Biographisches Mapping:** Erzeugt dynamische Zielpfade (`<Sphere>/From-<Zone>/<Date>/`).
- **Deduplizierung:** Implementiert via Hardlinks.

## 3. Sicherheits-Invarianten
- **Atomic Moves:** Dateien werden erst gelöscht, wenn das Ziel kryptographisch verifiziert wurde.
- **Interactive-Confirmation:** Physische Änderungen erfordern menschliche Validierung.

---
*Status: v1.3 Dokumentation synchronisiert mit Codebase.*
