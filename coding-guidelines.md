# Engineering Standards & Coding Guidelines: topoclean

## 1. Architektur-Prinzipien (SOLID & mehr)

Wir folgen den **SOLID** Prinzipien, um eine robuste und erweiterbare Codebasis zu garantieren:
- **S - Single Responsibility Principle (SRP):** Jedes Paket/Jede Funktion hat genau eine Aufgabe (z.B. nur Hashing, nur Ledger-Schreiben, nur Dateiscan).
- **O - Open/Closed Principle (OCP):** Neue Sortier-Vektoren können hinzugefügt werden, ohne den Kern-Scanner zu verändern.
- **L - Liskov Substitution Principle (LSP):** Schnittstellen (Interfaces) müssen so stabil sein, dass Implementierungen (z.B. verschiedene Ledger-Backends) austauschbar sind.
- **I - Interface Segregation Principle (ISP):** Kleine, spezifische Interfaces statt "Monster-Interfaces".
- **D - Dependency Inversion Principle (DIP):** High-Level Module (Business Logic) hängen nicht von Low-Level Modulen (Dateisystem) ab, sondern beide von Abstraktionen (Interfaces). Dies ermöglicht einfaches Mocking für Tests.

**Weitere Prinzipien:**
- **DRY (Don't Repeat Yourself):** Logik zur Pfad-Validierung oder Hash-Berechnung wird zentralisiert.
- **IoC (Inversion of Control):** Abhängigkeiten werden injiziert (Dependency Injection), nicht im Konstruktor hart kodiert.

## 2. Entwicklungs-Workflow (Git-Flow & TDD)

### 2.1. Branch-Strategie
1.  **`feat/<feature-name>` (Feature Branches):** Hier findet die eigentliche Entwicklung statt.
2.  **`staging`:** Integration aller Features. Hier laufen die **Smoke-Tests** und **großen Integrationstests**.
3.  **`main`:** Der "Stable" Release. Nur Code, der das Staging erfolgreich durchlaufen hat.

### 2.2. Test-Driven Development (TDD)
Jeder Zyklus folgt dem **Red-Green-Refactor** Prinzip:
1.  **Red:** Schreibe einen fehlerhaften Test für die neue Anforderung.
2.  **Green:** Schreibe den minimalen Code, um den Test zu bestehen.
3.  **Refactor:** Optimiere den Code unter Beibehaltung der Test-Grün-Phase.

### 2.3. Test-Hierarchie
- **Unit Tests:** Testen isolierte Funktionen (z.B. Hash-Algorithmus).
- **Integration Tests:** Testen das Zusammenspiel (z.B. Scanner + Vektor-Logik).
- **Smoke Tests:** Kurze End-to-End Tests auf dem Staging-Branch (z.B. "Bewegt eine Datei und rollt sie erfolgreich zurück").

## 3. Dokumentations-Pflichten

### 3.1. Developer Documentation (Dev-Doku)
Für jedes neue Feature muss eine technische Dokumentation (`docs/dev/<feature>.md`) erstellt werden. Sie muss so klar sein, dass ein menschlicher Entwickler das System ohne KI-Hilfe verstehen und warten kann.
- Architektur-Diagramme (ASCII/Mermaid).
- Erläuterung der Invarianten.
- Bekannte Edge-Cases.

### 3.2. User Documentation (User-Doku)
Jeder Befehl und jede Konfigurationsmöglichkeit muss in der `README.md` oder in `docs/user/` für den Endanwender erklärt werden. Fokus: Sicherheit und "Wie mache ich einen Rollback?".

## 4. Go-spezifische Richtlinien
- Nutzung von `go fmt` und `go vet`.
- Explizites Error-Handling (keine ignorierten Fehler!).
- Kommentare an allen exportierten Typen und Funktionen.
- Strukturierung in `cmd/` (Binaries) und `internal/` (Core-Logik, geschützt vor externem Import).

---
*Status: Engineering-Regelwerk aktiv. Jede Zeile Code muss diesen Standard atmen.*
