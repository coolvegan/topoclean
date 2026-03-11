# Konzept: Multi-Zonen-Management & Biographische Kartierung

## 1. Vision
`topoclean` soll nicht nur das Home-Verzeichnis, sondern alle "Entropie-Quellen" (Downloads, Desktop, Temp-Ordner) erfassen. Dabei muss die **biographische Herkunft** einer Datei gewahrt bleiben, damit der Nutzer versteht, warum eine Datei in einer bestimmten Sphäre gelandet ist.

## 2. Das Zonen-Modell (Entropy Zones)
Eine "Zone" ist ein Quellverzeichnis, das regelmäßig von Entropie gereinigt werden muss. Diese Zonen werden in einer zentralen Konfiguration definiert.

### 2.1. Standard-Zonen
- `~/Downloads` (Der Bahnhof der Unordnung)
- `~/Schreibtisch` oder `~/Desktop` (Das visuelle Rauschen)
- `~/tmp` oder `~/Temp` (Flüchtige Fragmente)

## 3. Biographische Kartierung (Source Mapping)
Um den Kontext zu wahren, werden Dateien innerhalb der 7 Ziel-Sphären in Unterordner verschoben, die ihre Herkunft widerspiegeln.

**Beispiel-Transformation:**
- `~/Downloads/rechnung.pdf` -> `~/02-Identity/From-Downloads/2026-03/rechnung.pdf`
- `~/Desktop/skizze.png` -> `~/05-Media/From-Desktop/2026-03/skizze.png`
- `~/tmp/script.py` -> `~/03-Creation/From-Tmp/2026-03/script.py`

## 4. Konfigurations-Architektur (`config.json`)
Die Steuerung erfolgt über eine soterische Konfigurationsdatei in `~/.config/topoclean/config.json`.

```json
{
  "version": "1.0",
  "heptagon_root": "~/",
  "zones": [
    {
      "path": "~/Downloads",
      "strategy": "aggressive",
      "keep_days": 30
    },
    {
      "path": "~/Desktop",
      "strategy": "immediate"
    }
  ],
  "mapping": {
    "preserve_origin": true,
    "date_format": "2006-01"
  }
}
```

## 5. Technische Umsetzung (Go)
- **`Config-Modul`**: Lädt und validiert die JSON-Konfiguration.
- **`Zone-Scanner`**: Iteriert über alle definierten Zonen statt nur über das Home-Verzeichnis.
- **`Path-Strategy`**: Eine neue Vektor-Strategie, die den Quellpfad auswertet und den Zielpfad um das `From-<Zone>` Präfix erweitert.

## 6. Soterischer Ausblick: Der Topo-Daemon
In einer künftigen Ausbaustufe kann `topoclean` als Hintergrunddienst (systemd-Unit oder Launchd) agieren, der mittels `inotify` (Dateisystem-Events) sofort reagiert, wenn eine Datei in einer Entropie-Zone landet und sie soterisch in die Ordnung überführt.

---
*Status: Konzept für Multi-Zonen-Management definiert. Bereit zur architektonischen Integration.*
