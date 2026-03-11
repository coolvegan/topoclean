# Technische Dokumentation: topoclean Architektur

## 1. Systemübersicht
`topoclean` ist modular aufgebaut, um SRP (Single Responsibility Principle) und IoC (Inversion of Control) zu gewährleisten. Die Kernmodule sind in `internal/` gekapselt.

## 2. Die Module

### 2.1. Ledger (internal/ledger)
Das Ledger ist das "Zentrale Gedächtnis". Es nutzt eine SQLite-Datenbank (transpiliert aus C, daher CGO-frei).
- **Tabellen:** `transactions` (Metadaten) & `operations` (konkrete Dateibewegungen).
- **Invarianten:** Jede Transaktion hat eine UUID und einen Status (`Pending`, `Committed`, `RolledBack`).
- **Traceability:** Lückenlose Erfassung von `Source`, `Dest`, `Hash` und `Size`.

### 2.2. Scanner (internal/scanner)
Der Scanner ist das "Auge" des Systems.
- **Magic Bytes:** Nutzt `http.DetectContentType` (512 Byte Probe), um das wahre Wesen einer Datei (MIME) zu erkennen.
- **No-Fly-Zone:** Ignoriert Dotfiles (`.*`) und das Ledger selbst (`*.db`).
- **FileInfo:** Liefert Pfad, Größe, Endung und MIME an das System.

### 2.3. Vector (internal/vector)
Der Vector ist das "Gehirn" (Klassifizierung).
- **Heptagon-Regeln:** Mappt MIME-Typen und Dateiendungen auf die 7 Sphären.
- **Hierarchie:** Inhaltsanalyse (MIME) triumphiert über die Dateiendung.

### 2.4. App Orchestrator (internal/app)
Das "Herz", das alle Komponenten verbindet.
- **Execute-Zyklus:** `Scan` -> `Plan` -> `Begin Tx` -> `For each: (Hash -> Copy -> Verify -> Log -> Delete)` -> `Commit Tx`.
- **Rollback-Inversion:** `Load Tx` -> `For each op: (Verify Hash -> Restore -> Delete Target)` -> `Update Tx State`.

## 3. Sicherheitsmechanismen (Soterik)
- **Anti-Corruption:** SHA-256 Hashing vor und nach jedem Schreibvorgang.
- **Graceful Degradation:** Fehler bei einzelnen Dateien führen nicht zum Abbruch der gesamten Transaktion, sondern werden protokolliert.
- **Interactive-Mode:** Explizite Bestätigung für alle destruktiven Operationen.

## 4. Wartung & Tests (TDD)
- **Unit Tests:** Jede funktionale Änderung muss durch einen Test in `*_test.go` gedeckt sein.
- **Integration Tests:** Der vollständige Zyklus von Scan bis Rollback wird auf `staging` validiert.

---
*Status: Architektur verifiziert. Prozess-Sicherheit gewährleistet.*
