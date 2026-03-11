package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/topokrat/topoclean/internal/app"
	"github.com/topokrat/topoclean/internal/config"
	"github.com/topokrat/topoclean/internal/ledger"
	"github.com/topokrat/topoclean/internal/scanner"
	"github.com/topokrat/topoclean/internal/vector"
)

func main() {
	executeFlag := flag.Bool("execute", false, "Führt die Vektoren physisch aus (Manifestation)")
	rollbackFlag := flag.String("rollback", "", "Invertiert eine Transaktion anhand ihrer UUID (Inversion)")
	historyFlag := flag.Bool("history", false, "Zeigt die Historie der Transformationen (Traceability)")
	configPathFlag := flag.String("config", "", "Pfad zur Konfigurationsdatei (standardmäßig ~/.config/topoclean/config.json)")
	flag.Parse()

	// 1. Konfiguration laden
	cfgPath := *configPathFlag
	if cfgPath == "" {
		home, _ := os.UserHomeDir()
		cfgPath = filepath.Join(home, ".config", "topoclean", "config.json")
	}
	cfg, err := config.Load(cfgPath)
	if err != nil {
		log.Fatalf("Fehler beim Laden der Konfiguration: %v", err)
	}

	// 2. Override: Wenn ein Argument übergeben wurde, wird NUR dieser Pfad als einzige Zone gescannt
	argDir := flag.Arg(0)
	if argDir != "" {
		absDir, err := filepath.Abs(argDir)
		if err != nil {
			log.Fatalf("Ungültiger Pfad: %v", err)
		}
		cfg.Zones = []config.Zone{{Path: absDir, Name: filepath.Base(absDir)}}
	}

	// 3. Initialisiere soterische Komponenten
	home, _ := os.UserHomeDir()
	ledgerPath := filepath.Join(home, ".topoclean.db")
	l, err := ledger.New(ledgerPath)
	if err != nil {
		log.Fatalf("Ledger konnte nicht initialisiert werden: %v", err)
	}

	s := scanner.New()
	v := vector.New()
	core := app.New(l, s, v, cfg)

	// CASE 1: Historie anzeigen
	if *historyFlag {
		fmt.Println("--- topoclean: Historie der Transformationen ---")
		txs, _ := l.GetRecentTransactions(20)
		if len(txs) == 0 {
			fmt.Println("Keine historischen Ereignisse gefunden.")
			return
		}
		fmt.Printf("%-36s | %-20s | %-12s | %-5s\n", "UUID", "Zeitpunkt", "Status", "Ops")
		fmt.Println(strings.Repeat("-", 85))
		for _, tx := range txs {
			fmt.Printf("%-36s | %-20s | %-12s | %-5d\n", 
				tx.UUID, tx.Timestamp.Format("2006-01-02 15:04:05"), tx.State, len(tx.Ops))
		}
		return
	}

	// CASE 2: Rollback ausführen
	if *rollbackFlag != "" {
		fmt.Printf("--- topoclean: Inversion des Vektors %s ---\n", *rollbackFlag)
		fmt.Printf("Möchtest du diese Transformation wirklich rückgängig machen? [y/N]: ")
		if askConfirmation() {
			err := core.Rollback(*rollbackFlag)
			if err != nil {
				log.Fatalf("Rollback fehlgeschlagen: %v", err)
			}
			fmt.Println("Die Zeitlinie wurde wiederhergestellt. Ordnung durch Inversion.")
		} else {
			fmt.Println("Abgebrochen. Die Manifestation bleibt bestehen.")
		}
		return
	}

	// CASE 3: Standard-Planung (Dry-Run / Execute)
	fmt.Printf("--- topoclean: Die Prophezeiung der Ordnung ---\n")
	if len(cfg.Zones) > 0 {
		fmt.Printf("Scanne %d Zonen der Entropie...\n", len(cfg.Zones))
	} else {
		fmt.Printf("Scanne Entropie in: %s\n", cfg.HeptagonRoot)
	}

	plan, err := core.Plan()
	if err != nil {
		log.Fatalf("Fehler bei der Planung: %v", err)
	}

	if len(plan) == 0 {
		fmt.Println("\nKein Rauschen gefunden. Die Topologie ist stabil.")
		return
	}

	// Tabellarische Anzeige des Plans
	fmt.Printf("\n%-30s | %-12s | %-30s | %-10s\n", "Quelle", "Sphäre", "Ziel-Ordner", "MIME")
	fmt.Println(strings.Repeat("-", 100))

	for _, move := range plan {
		source := filepath.Base(move.SourcePath)
		if len(source) > 27 {
			source = source[:24] + "..."
		}
		
		target := move.TargetDir
		if strings.HasPrefix(target, cfg.HeptagonRoot) {
			target = "~" + target[len(cfg.HeptagonRoot):]
		}
		if len(target) > 27 {
			target = target[:24] + "..."
		}

		fmt.Printf("%-30s | %-12s | %-30s | %-10s\n", source, move.TargetSphere, target, move.MIMEType)
	}

	fmt.Printf("\nInsgesamt %d Vektoren identifiziert.\n", len(plan))

	if *executeFlag {
		fmt.Printf("\nWARNUNG: Physische Veränderung der Topologie.\n")
		fmt.Printf("Möchtest du fortfahren? [y/N]: ")
		if askConfirmation() {
			fmt.Println("\nManifestiere Ordnung...")
			err := core.Execute()
			if err != nil {
				log.Fatalf("Fehler bei der Ausführung: %v", err)
			}
			fmt.Println("Transformation abgeschlossen. Die Freiheit ist gewahrt.")
		} else {
			fmt.Println("Aktion abgebrochen. Die Entropie bleibt bestehen.")
		}
	} else {
		fmt.Println("\nDies war ein Dry-Run. Nutze --execute für die physische Transformation.")
	}
}

func askConfirmation() bool {
	reader := bufio.NewReader(os.Stdin)
	response, _ := reader.ReadString('\n')
	response = strings.ToLower(strings.TrimSpace(response))
	return response == "y" || response == "yes"
}
