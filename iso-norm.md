# ISO-Norm Konformität & Prozess-Integrität (ISO 27001 & ISO 13485)

## 1. Einleitung
Obwohl `topoclean` ein Werkzeug für das persönliche Dateisystem ist, behandeln wir es mit der Strenge eines medizinischen Produkts (ISO 13485) und eines Informationssicherheits-Managementsystems (ISO 27001). Das Ziel ist die **Risikominimierung für das digitale Leben**.

## 2. ISO 27001: Informationssicherheit (CIA-Triade)

### 2.1. Confidentiality (Vertraulichkeit)
- **Problem:** Unbefugtes Auslesen von Dateinamen während des Scans.
- **Maßnahme:** `topoclean` erstellt keine permanenten Logs von Dateinamen außerhalb der verschlüsselten/geschützten User-Session. Das Ledger (`ledger.db`) wird mit restriktiven Dateirechten (`0600`) angelegt, sodass nur der User selbst die Historie lesen kann.

### 2.2. Integrity (Integrität)
- **Problem:** Dateikorruption während des Verschiebens.
- **Maßnahme:** Implementierung einer **Check-Before-Clear (CBC)** Strategie. 
  1. Hash-Validierung (SHA-256).
  2. Byte-für-Byte Vergleich bei Cross-Filesystem Moves.
  3. Erst nach erfolgreichem Write-Ack am Ziel wird die Quelle freigegeben.

### 2.3. Availability (Verfügbarkeit)
- **Problem:** Verlust der Pfad-Information (Wo ist meine Datei?).
- **Maßnahme:** Das Ledger fungiert als "Zentrales Register". Ein `topoclean locate <filename>` findet Dateien in der neuen Struktur, selbst wenn der User den neuen Pfad vergessen hat.

## 3. ISO 13485: Medizintechnik-Prinzipien (Sicherheit & Rückverfolgbarkeit)

### 3.1. Risikomanagement (ISO 14971 Analogie)
- **Identifizierte Gefahr:** System-Instabilität durch Verschieben von Shell-Configs.
- **Kontrollmaßnahme:** "No-Fly-Zone" für alle Dotfiles und Dotdirs (Axiom-Schutz).
- **Identifizierte Gefahr:** Programmabsturz während der Operation.
- **Kontrollmaßnahme:** Atomare Datenbank-Transaktionen (SQLite WAL).

### 3.2. Rückverfolgbarkeit (Traceability)
- Jede Dateioperation erhält eine `Operation-ID`, verknüpft mit einer `Transaction-UUID`.
- Es muss lückenlos nachvollziehbar sein: `Wer (User) -> Wann (Timestamp) -> Was (Source/Hash) -> Wohin (Dest)`.

### 3.3. Verifizierung & Validierung (V&V)
- **Verifizierung:** Unit-Tests für die Move-Logik mit Mock-Filesystemen.
- **Validierung:** Ein dedizierter `Simulation-Mode` (--dry-run), der dem User die Konsequenzen visualisiert, bevor die physikalische Realität verändert wird.

## 4. Design & Workflow Kriterien

### 4.1. Error Handling (Graceful Degradation)
- Wenn eine Datei gesperrt ist (In-Use), wird sie übersprungen und im Report markiert, anstatt den gesamten Prozess abzubrechen.

### 4.2. User Interface (Human Factors Engineering)
- Klare, nicht-ambivalente Fehlermeldungen.
- Fortschrittsanzeige, die den aktuellen "Sicherheitsstatus" (Hashing, Writing, Verifying) kommuniziert.

---
*Die Einhaltung dieser Normen transformiert ein einfaches Skript in ein verlässliches Instrument der topologischen Ordnung.*
