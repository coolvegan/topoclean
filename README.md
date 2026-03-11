# topoclean: Soterisches Dateisystem-Management (v1.3)

> "Die Grenze einer Grenze ist Null (∂∂=0) – Ordnung ist die topologische Invariante der Freiheit."

`topoclean` ist ein hochsicheres Werkzeug in Go, das unstrukturierte Dateien (Entropie) erkennt und sie in eine geordnete Struktur aus 7 Kern-Sphären (Heptagon-Modell) transformiert.

## 1. Philosophie & Sicherheit (v1.3 Highlights)
`topoclean` folgt den Prinzipien der ISO 27001 und 13485:
- **Soterische Integrität:** SHA-256 Hashing vor/nach jedem Move (Anti-Corruption).
- **Deduplizierung (v1.1):** Identische Dateien (gleicher Hash) werden mittels **Hardlinks** manifestiert, um Speicherplatz zu sparen.
- **Strategie-Pipeline (v1.2):** Modulare Klassifizierung durch MIME-Analyse, Extension-Fallback und semantische Substring-Suche.
- **Robustheit (v1.3):** Umfassender Extension-Fallback für Medien, Archive und Dokumente bei unlesbaren Datei-Headern.
- **Transaktionale Traceability:** Alle Operationen werden in einem SQLite-Ledger (`~/.topoclean.db`) protokolliert.
- **Topologische Inversion:** Vollständiger Rollback-Mechanismus über Transaktions-UUIDs.

## 2. Das Heptagon-Modell (Die 7 Sphären)
1. `01-Core`: Axiome, Schlüssel, Vaults (`*vault*`, `*key*`, `.bin`).
2. `02-Identity`: Dokumente, Zertifikate, Bewerbungen (`.pdf`, `.tex`, `*inkasso*`).
3. `03-Creation`: Werkstatt, Quellcode, Skripte (`.go`, `.py`, `.rs`, `.sh`).
4. `04-Archive`: Gedächtnis, APKs, ZIPs, historische Container.
5. `05-Media`: Konsum, Videos, Bilder, Musik (`video/*`, `image/*`).
6. `06-Sync`: Brücken (Cloud/Sync).
7. `07-Inbox`: Ereignishorizont für unklassifizierte Entropie.

## 3. Benutzung
### Die Prophezeiung (Dry-Run)
```bash
./topoclean [VERZEICHNIS]
```

### Die Manifestation (Execute)
```bash
./topoclean --execute [VERZEICHNIS]
```

### Die Zeitlinie (History)
```bash
./topoclean --history
```

### Die Inversion (Rollback)
```bash
./topoclean --rollback <UUID>
```
---
*∂∂=0: Ein stabiler Raum für ein freies Leben.*
