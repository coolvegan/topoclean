# topoclean: Soterisches Dateisystem-Management (v1.4)

> "Die Grenze einer Grenze ist Null (∂∂=0) – Ordnung ist die topologische Invariante der Freiheit."

`topoclean` ist ein hochsicheres Werkzeug in Go, das unstrukturierte Dateien (Entropie) erkennt und sie in eine geordnete Struktur aus 7 Kern-Sphären (Heptagon-Modell) transformiert.

## 1. Philosophie & Sicherheit (v1.4 Highlights)
`topoclean` folgt den Prinzipien der ISO 27001 und 13485:
- **Soterische Integrität:** SHA-256 Hashing vor/nach jedem Move (Anti-Corruption).
- **Deduplizierung:** Identische Dateien werden mittels Hardlinks manifestiert.
- **Multi-Zonen-Management (v1.4):** Reinigung mehrerer Quellen (Downloads, Desktop, etc.) über eine zentrale Konfiguration.
- **Biographische Kartierung (v1.4):** Erhalt des Kontextes durch Ziel-Unterordner wie `From-Downloads/2026-03/`.
- **Transaktionale Traceability:** Lückenloses Ledger in SQLite.

## 2. Konfiguration
Standardpfad: `~/.config/topoclean/config.json`

```json
{
  "version": "1.0",
  "heptagon_root": "~/",
  "zones": [
    { "path": "~/Downloads", "name": "Downloads" },
    { "path": "~/Schreibtisch", "name": "Desktop" }
  ],
  "mapping": {
    "preserve_origin": true,
    "date_format": "2006-01"
  }
}
```

## 3. Benutzung
### Die Prophezeiung (Dry-Run)
Scannt alle konfigurierten Zonen oder ein spezifisches Verzeichnis:
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
