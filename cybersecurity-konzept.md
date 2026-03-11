# Cybersecurity- & Stabilitäts-Konzept: topoclean

## 1. Die Soterische Prämisse (Unversehrtheit)
Die Angst vor systemischem Umbruch ist die gesunde Reaktion eines Administrators auf eine Änderung der Entropie. `topoclean` folgt dem Prinzip der **Nicht-Invasivität** gegenüber System-Axiomen (SSH, GPG, Configs).

## 2. Schutz der System-Axiome (Protected Spaces)

### 2.1. Die "No-Fly-Zone" (White-Listing)
Bestimmte Pfade und Dateitypen sind für die `topoclean`-Vektoren **unsichtbar** und werden niemals bewegt:
- **Dotfiles & Dotdirs:** `~/.ssh/`, `~/.gnupg/`, `~/.config/`, `~/.local/`, `~/.bashrc`, `~/.fish_variables`, etc.
- **System-Kritische Sockets/Pipes:** Alles, was `os.ModeSocket` oder `os.ModeNamedPipe` entspricht.
- **Hidden Files:** Dateien, die mit einem Punkt beginnen, bleiben an ihrem Ort, um Shell-Konfigurationen nicht zu brechen.

### 2.2. Die "Shadow-Link"-Strategie (Kompatibilitäts-Layer)
Um sicherzustellen, dass Dienste (z.B. ein Skript, das hart auf `~/Videos` prüft) nicht brechen:
- `topoclean` verschiebt physisch, hinterlässt aber bei Bedarf einen **Symbolic Link** am Ursprungsort.
- Das System "sieht" die Datei am alten Ort, während das menschliche Auge in der neuen Struktur Ordnung findet.

## 3. Cyber-Security: Schutz vor Datenverlust & Korruption

### 3.1. Kryptographische Identität (Anti-Corruption)
1. **Pre-Move Hash:** Vor jeder Operation wird ein SHA-256 Hash der Quelldatei erstellt.
2. **Atomic Write:** Die Datei wird an das Ziel kopiert (nicht verschoben, falls über Mountpoints hinweg).
3. **Post-Move Hash:** Der Hash der Zieldatei wird verifiziert.
4. **Finalization:** Erst nach erfolgreicher Verifizierung wird die Quelldatei (optional) gelöscht oder in einen `04-Archive/Stage`-Ordner verschoben.

### 3.2. Schutz vor Privilege Escalation / Path Injection
- `topoclean` läuft strikt mit User-Rechten (kein `sudo`).
- **Path Sanitization:** Alle Pfade werden absolut aufgelöst (`filepath.Abs`) und gegen "Path Traversal" Attacken geprüft (kein Ausbruch aus `/home/user/` möglich).

## 4. Transaktions-Sicherheit (Das Ledger)

### 4.1. Das Unzerstörbare Protokoll
Jede Transaktion wird in einer `ledger.db` (SQLite mit WAL-Mode für Crash-Resistenz) gespeichert. 
- Ein plötzlicher Stromausfall hinterlässt die Datei entweder am Quellort, am Zielort oder an beiden (Redundanz).
- Das Ledger erkennt unvollständige Operationen beim nächsten Start.

### 4.2. Instant Rollback (Die Zeitmaschine)
Sollte ein Dienst (z.B. `fish`) nach dem Aufräumen Fehler melden, kann mit einem einzigen Befehl die **Inversion aller Vektoren** ausgelöst werden. Das Ledger schiebt jede Datei exakt auf den Byte-Offset ihres Ursprungs zurück.

## 5. Menschliche Validierung (Der Rat der Ethik)
- **Interactive-Mode:** Jede Bewegung muss durch ein `[Y/n]` bestätigt werden, sofern nicht explizit `--force` gesetzt ist.
- **Visual Diff:** Zeigt vorab: `~/Anschreiben.pdf` -> `~/02-Identity/Career/Anschreiben.pdf`.

---
*Sicherheit ist die topologische Invariante der Freiheit. Wir bewegen nur das Rauschen, nicht das Signal.*
