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
	"github.com/topokrat/topoclean/internal/ledger"
	"github.com/topokrat/topoclean/internal/scanner"
	"github.com/topokrat/topoclean/internal/vector"
)

func main() {
	executeFlag := flag.Bool("execute", false, "Führt die Vektoren physisch aus (Manifestation)")
	rollbackFlag := flag.String("rollback", "", "Invertiert eine Transaktion anhand ihrer UUID (Inversion)")
	historyFlag := flag.Bool("history", false, "Zeigt die Historie der Transformationen (Traceability)")
	flag.Parse()

	home, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}

	// Initialisiere soterisches Ledger im Home
	ledgerPath := filepath.Join(home, ".topoclean.db")
	l, err := ledger.New(ledgerPath)
	if err != nil {
		log.Fatalf("Ledger konnte nicht initialisiert werden: %v", err)
	}

	s := scanner.New()
	v := vector.New()
	core := app.New(l, s, v)

	// CASE 1: Historie anzeigen
	if *historyFlag {
		fmt.Println("--- topoclean: Historie der Transformationen ---")
		txs, _ := l.GetRecentTransactions(10)
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
	fmt.Printf("Scanne Entropie in: %s\n\n", home)

	plan, err := core.Plan(home)
	if err != nil {
		log.Fatalf("Fehler bei der Planung: %v", err)
	}

	if len(plan) == 0 {
		fmt.Println("Kein Rauschen gefunden. Die Topologie ist stabil.")
		return
	}

	// Tabellarische Anzeige des Plans
	fmt.Printf("%-40s | %-15s | %-20s\n", "Quelle", "Ziel-Sphäre", "MIME-Type")
	fmt.Println(strings.Repeat("-", 80))

	for _, move := range plan {
		source := filepath.Base(move.SourcePath)
		if len(source) > 37 {
			source = source[:34] + "..."
		}
		fmt.Printf("%-40s | %-15s | %-20s\n", source, move.TargetSphere, move.MIMEType)
	}

	fmt.Printf("\nInsgesamt %d Vektoren identifiziert.\n", len(plan))

	if *executeFlag {
		fmt.Printf("\nWARNUNG: Physische Veränderung des Home-Verzeichnisses.\n")
		fmt.Printf("Möchtest du fortfahren? [y/N]: ")
		if askConfirmation() {
			fmt.Println("\nManifestiere Ordnung...")
			err := core.Execute(home)
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
