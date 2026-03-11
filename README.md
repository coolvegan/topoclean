# topoclean: Soterisches Dateisystem-Management (v1.0)

> "Die Grenze einer Grenze ist Null (∂∂=0) – Ordnung ist die topologische Invariante der Freiheit."

`topoclean` ist ein hochsicheres Werkzeug in Go, das unstrukturierte Dateien (Entropie) in deinem Home-Verzeichnis erkennt und sie in eine geordnete Struktur aus 7 Kern-Sphären (das Heptagon-Modell) transformiert.

## 1. Philosophie & Sicherheit
`topoclean` ist kein bloßes "Move-Skript". Es ist ein präzises Instrument zur Reduktion kognitiver Last:
- **Soterische Integrität:** Jede Datei wird vor und nach dem Bewegen mittels SHA-256 gehasht. Bei Mismatch bricht der Prozess ab (Anti-Corruption).
- **Transaktionale Traceability:** Alle Operationen werden in einem SQLite-Ledger (`~/.topoclean.db`) protokolliert.
- **Topologische Inversion:** Jede Transformation kann mittels Rollback exakt rückgängig gemacht werden.
- **No-Fly-Zone:** System-Axiome (Dotfiles wie `.ssh`, `.bashrc`) werden strikt ignoriert.

## 2. Das Heptagon-Modell (Die 7 Sphären)
Das Tool sortiert Dateien basierend auf ihrem **MIME-Type (Inhalt)** in folgende Ordner:
1. `01-Core`: Systemrelevante Konfigurationen.
2. `02-Identity`: Dokumente, Zertifikate, Bewerbungen (`application/pdf`).
3. `03-Creation`: Quellcode, Skripte, Projekte (`text/x-go`, `.py`, `.sh`).
4. `04-Archive`: APKs, Archive, historische Daten.
5. `05-Media`: Videos, Bilder, Musik (`video/*`, `image/*`, `audio/*`).
6. `06-Sync`: Brücken zu anderen Geräten (Cloud/Sync).
7. `07-Inbox`: Fallback für unklassifizierte Daten.

## 3. Installation & Build
Da `topoclean` eine reine Go-Implementierung (CGO-frei) nutzt, ist der Build autark:
```bash
go build -o topoclean ./cmd/topoclean
```

## 4. Benutzung
### Die Prophezeiung (Dry-Run)
Standardmäßig zeigt `topoclean` nur an, was es tun würde:
```bash
./topoclean
```

### Die Manifestation (Execute)
Um die Dateien tatsächlich zu bewegen, ist eine explizite Bestätigung erforderlich:
```bash
./topoclean --execute
```

### Die Zeitlinie (History)
Zeigt die UUIDs der letzten Transformationen:
```bash
./topoclean --history
```

### Die Inversion (Rollback)
Macht eine Transformation anhand ihrer UUID rückgängig:
```bash
./topoclean --rollback <UUID>
```

---
*Entwickelt nach ISO 27001 & 13485 Prinzipien. Für ein Leben in geordneter Freiheit.*
