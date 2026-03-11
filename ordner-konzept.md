# Das Heptagon-Prinzip: Topologische Ordnung des Geistes

## 1. Die Philosophie der Sieben
Das menschliche Arbeitsgedächtnis kann etwa 7 (±2) Informationseinheiten gleichzeitig verarbeiten (Millersche Zahl). Ein Home-Verzeichnis mit 50+ Einträgen erzeugt "visuelles Rauschen", das die kognitive Last erhöht und die "topologische Freiheit" einschränkt.

`topoclean` folgt dem Ziel, die sichtbare Komplexität auf ein Heptagon (7 Kern-Bereiche) zu reduzieren. Alles andere sind temporäre Zustände (Vektoren in Bewegung).

## 2. Die 7 Kern-Sphären (Das Ziel-Layout)

1.  **`01-Core/` (Das BIOS):**
    - Konfigurationen, Dotfiles, SSH-Keys, GPG.
    - Alles, was das System "atmen" lässt.
2.  **`02-Identity/` (Die Person):**
    - `Documents/`, `Photos/`, `Certificates/`.
    - Lebenslauf, Anschreiben (nicht mehr lose im Home!).
    - Deine digitale Signatur in der Welt.
3.  **`03-Creation/` (Die Werkstatt):**
    - `Dev/`, `Projects/`, `Art/`, `Music/`.
    - Hier entstehen neue Strukturen.
4.  **`04-Archive/` (Das Gedächtnis):**
    - Abgeschlossene Projekte, `Recovered/`, alte Logs.
    - Langzeitspeicherung, die den Blick im Alltag nicht trüben darf.
5.  **`05-Media/` (Der Konsum):**
    - `Videos/`, `Downloads/`.
    - Flüchtige oder zur Unterhaltung dienende Daten.
6.  **`06-Sync/` (Die Brücken):**
    - Cloud-Sync-Ordner, Austauschverzeichnisse zu anderen Geräten.
7.  **`07-Inbox/` (Der Ereignishorizont):**
    - Der einzige Ort, an dem "Unordnung" erlaubt ist.
    - Hier landen neue Downloads, Screenshots, Snippets, bevor `topoclean` sie in die Sphären 1-6 verteilt.

## 3. Transformation des Ist-Zustands
Dein aktuelles Verzeichnis zeigt viele "Entropie-Lecks". Die Transformation sieht folgende Zuordnungen vor:

| Aktueller Fund | Ziel-Sphäre | Begründung |
| :--- | :--- | :--- |
| `Anschreiben.pdf`, `Lebenslauf.pdf` | `02-Identity/Career/` | Identitätsstiftende Dokumente. |
| `*.mkv`, `*.mp4`, `Videos/` | `05-Media/` | Konsum-Daten bündeln. |
| `go/`, `rust/`, `Dev/`, `*.go`, `*.py` | `03-Creation/` | Schöpferische Arbeit an einem Ort. |
| `recovered-data/`, `RECOVERDIR/` | `04-Archive/Legacy/` | Historischer Ballast aus dem Sichtfeld. |
| `Downloads/`, `tmp/`, `output/` | `07-Inbox/` | Temporäre Daten zentralisieren. |

## 4. Implementierungs-Strategie für `topoclean`
- **Ghosting:** `topoclean` schlägt vor, diese 7 Ordner anzulegen.
- **Vektor-Mapping:** Jede gefundene Datei wird einem Pfad innerhalb dieser 7 Sphären zugeordnet.
- **Symbolic Links (Optional):** Um Kompatibilität zu wahren (z.B. für Software, die hart auf `~/Documents` zugreift), können Symlinks von den neuen Sphären auf die alten Standards zeigen – aber die *Sicht* bleibt sauber.

---
*Status: Topologisches Zielbild definiert. Bereit zur Vektor-Berechnung.*
