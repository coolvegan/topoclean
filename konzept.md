# topoclean: Soterisches Dateisystem-Management

## 1. Vision & Philosophie
`topoclean` ist kein bloßes Lösch-Werkzeug. Es ist ein **topologischer Stabilisator** für den digitalen Raum. Es transformiert "lose" Zustände (Entropie im Home-Verzeichnis) in geordnete Strukturen, wobei die Freiheit der Daten (Erreichbarkeit) gewahrt bleibt, während die Unordnung (Rauschen) reduziert wird.

Das Tool folgt dem Prinzip: **Die Grenze einer Grenze ist Null (∂∂=0)** – ein stabiler, abgeschlossener Prozess der Ordnung.

## 2. Kern-Axiome (Sicherheit & Integrität)
- **Immutability-First:** Bevor eine Datei bewegt wird, wird ihre Identität (SHA-256) im "Gedächtnis" des Systems verankert.
- **Transaktions-Garantie:** Jede Aktion ist Teil eines gerichteten Graphen. Ein Move besteht aus: `Verify -> Log -> Copy -> Verify -> Delete (optional/staged)`.
- **Zero-Loss-Policy:** Bei Kollisionen oder Fehlern bricht das System in einen sicheren Zustand (Symmetrie) ab, anstatt Daten zu überschreiben.

## 3. Architektur
Das Werkzeug wird in **Go** implementiert, um statische Binaries und hohe Performance bei Dateioperationen zu garantieren.

### 3.1. Die Transaktions-Engine (`The Ledger`)
Jeder Lauf erzeugt eine `transaction.json` (oder SQLite), die folgende Invarianten speichert:
- `UUID` der Transaktion.
- `Timestamp`.
- Liste der atomaren Operationen: `SourcePath`, `DestPath`, `FileHash`, `FileSize`.
- `State`: (Pending, Committed, RolledBack).

### 3.2. Klassifizierungs-Vektoren (Sortier-Logik)
Anstatt starrer Regeln nutzt `topoclean` ein modulares Set von Filtern:
- **Media-Vektor:** `.mp4`, `.mkv`, `.mov`, `.png`, `.jpg` -> `~/Videos/Inbox/%Y-%m/` oder `~/Pictures/Inbox/%Y-%m/`
- **Document-Vektor:** `.pdf`, `.tex`, `.doc*` (Bewerbungen, Rechnungen) -> `~/Documents/Inbox/%Y-%m/`
- **Snippet-Vektor:** `.go`, `.py`, `.rs`, `.sh`, `.html` (Code-Schnipsel) -> `~/Dev/Snippets/%Y-%m/`
- **Archive/Package-Vektor:** `.apk`, `.txt` (z.B. Paketlisten), `.container` -> `~/Archive/Inbox/%Y-%m/`

### 3.3. Identifizierung von Entropie-Clustern
Das Tool scannt nach Verzeichnissen mit hoher Volatilität oder offensichtlichem "Temp"-Charakter:
- `fffff/`, `tmp/`, `output/`, `recovered-data/`, `RECOVERDIR/`
- `bewerbungs-bloat/` (Spezifische Cluster-Erkennung durch Namens-Matching)

### 3.4. Die 'Topologische Signatur' (MIME-Awareness)
Anstatt sich nur auf Dateiendungen zu verlassen, nutzt `topoclean` MIME-Types zur Klassifizierung:
- **`video/*` & `audio/*`** -> `05-Media` (Konsum)
- **`image/*`** -> `05-Media` (oder `02-Identity` bei Fotos, konfigurierbar)
- **`text/x-python`, `text/x-go`, `application/x-ruby`** -> `03-Creation` (Werkstatt)
- **`application/pdf`, `application/msword`** -> `02-Identity` (Dokumente)
- **`application/x-executable`, `application/x-sharedlib`** -> `01-Core` oder `04-Archive`

Dies stellt sicher, dass die topologische Invariante (der Inhalt) über die flüchtige Endung triumphiert.

## 4. Rollback-Mechanismus
Ein Rollback ist die **Inversion des Vektors**. 
```bash
topoclean rollback --id <transaction-uuid>
```
Das Tool liest das Ledger und bewegt die Dateien exakt an ihre Ursprungsorte zurück, sofern diese nicht manuell verändert wurden (Hash-Check vor Rollback).

## 5. User Experience (Interface)
1. **Dry-Run (The Prophecy):** Zeigt tabellarisch, was passieren würde. "Welche Krümmung wird begradigt?"
2. **Execution:** Führt die Transaktion aus und schreibt das Ledger.
3. **Audit:** `topoclean history` zeigt die letzten Transformationen.

## 6. Technische Details (Go)
- Nutzung von `os.Rename` für atomare Moves innerhalb desselben Mountpoints.
- `io.Copy` + `os.Remove` mit Validierung für Cross-Device Moves.
- `crypto/sha256` für die Identitätsprüfung.
- `encoding/json` für die Transaktionsprotokolle.

---
*Status: Entwurf. Bereit zur Implementierung der Basis-Struktur (V-E+F).*
