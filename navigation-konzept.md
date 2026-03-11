# Navigations-Axiom: Soterische Suche & Auffindbarkeit

## 1. Vision
Struktur darf niemals Freiheit (Erreichbarkeit) einschränken. `topoclean` nutzt sein transaktionales Gedächtnis (Ledger), um dem Nutzer einen sofortigen, tiefen Einblick in die neue Topologie zu gewähren, ohne dass dieser manuell durch Verzeichnisse navigieren muss.

## 2. Das Locate-Prinzip (The Oracle)
Das Ledger fungiert als Oracle. Jede Datei, die jemals durch `topoclean` transformiert wurde, ist über ihren Namen, ihren Hash oder ihre Herkunft (Zone) auffindbar.

### 2.1. Suche nach Mustern
- **Befehl:** `topoclean locate <MUSTER>`
- **Logik:** SQL `LIKE` Abfrage auf das Feld `dest_path` im Ledger.
- **Ergebnis:** Eine tabellarische Liste aller passenden Dateien mit ihrem aktuellen Aufenthaltsort.

### 2.2. Suche nach Identität (Hash)
- **Befehl:** `topoclean locate --hash <SHA256>`
- **Logik:** Findet alle Kopien/Hardlinks einer Datei im gesamten System.

## 3. Integration & Komfort
Um die Navigation noch flüssiger zu machen, kann `topoclean locate` so konfiguriert werden, dass es nur die Pfade ausgibt (z.B. für die Weitergabe an `fzf` oder `xargs`).

---
*Status: Navigations-Konzept definiert. Bereit zur Implementierung des Oracle-Moduls.*
