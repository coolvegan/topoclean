# topoclean: Soterisches Dateisystem-Management (v1.5)

> "Die Grenze einer Grenze ist Null (∂∂=0) – Ordnung ist die topologische Invariante der Freiheit."

## 1. Philosophie & Sicherheit (v1.5 Highlights)
`topoclean` folgt den Prinzipien der ISO 27001 und 13485:
- **Lethe-Prinzip (v1.5):** Sicheres Löschen durch Staging im Limbo (`04-Archive/.trash`).
- **Soterische Integrität:** SHA-256 Hashing vor/nach jedem Move.
- **Multi-Zonen-Management:** Reinigung mehrerer Quellen (Downloads, Desktop, etc.).
- **Transaktionale Traceability:** Jede Datei-Operation ist im Ledger verewigt.

## 2. Benutzung
### Die Prophezeiung (Dry-Run)
Scannt alle konfigurierten Zonen:
```bash
./topoclean
```

### Das Vergessen (Forget)
Überführt eine Datei soterisch in den Limbo:
```bash
./topoclean --forget <PFAD>
```

### Die Manifestation (Execute)
Führt den Plan physisch aus:
```bash
./topoclean --execute
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
