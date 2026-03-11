# Lethe-Prinzip: Soterisches Löschen & Daten-Lebenszyklus

## 1. Vision
Löschen ist in der Topokratie kein destruktiver Akt, sondern die bewusste Entscheidung zur Reduktion von Komplexität. `topoclean` stellt sicher, dass dieses Vergessen kontrolliert, nachvollziehbar und reversibel (bis zum finalen Purge) erfolgt.

## 2. Der soterische Lösch-Prozess (The Limbo)
Anstatt Dateien sofort mittels `os.Remove` zu vernichten, durchlaufen sie zwei Phasen:

### Phase 1: Die Staging-Löschung (`forget`)
- **Aktion:** Die Datei wird aus ihrer Sphäre nach `04-Archive/.trash/` verschoben.
- **Ledger:** Die Operation wird als `State: StagedForDeletion` markiert.
- **Ziel:** Visuelle Bereinigung bei gleichzeitiger physikalischer Erhaltung für eine Sicherheitsfrist.

### Phase 2: Das endgültige Vergessen (`purge`)
- **Aktion:** Dateien im `.trash`, die die konfigurierte Haltefrist (z.B. 30 Tage) überschritten haben, werden physikalisch gelöscht.
- **Ledger:** Der Hash der Datei bleibt als "Epitaph" erhalten, markiert als `State: Purged`.

## 3. Automatisierte Exspiration
Zonen können in der `config.json` eine `keep_days` Richtlinie erhalten. `topoclean` erkennt abgelaufene Dateien in temporären Sphären (z.B. `07-Inbox/From-Temp`) und schlägt deren Überführung in den Limbo (`.trash`) vor.

## 4. Befehle
- `topoclean forget <pfad>`: Verschiebt eine Datei/Verzeichnis soterisch in den Limbo.
- `topoclean purge`: Bereinigt den Limbo von abgelaufenen Fragmenten.
- `topoclean --cleanup`: Integrierter Lauf, der Exspirationen prüft.

---
*Status: Lethe-Konzept definiert. Bereit zur Implementierung der Limbo-Logik.*
