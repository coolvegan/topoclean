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
	dryRunFlag := flag.Bool("dry-run", true, "Zeigt nur die Prophezeiung der Ordnung (Standard)")
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

	// Wenn --execute gesetzt ist UND --dry-run explizit auf false (oder default)
	if *executeFlag {
		fmt.Printf("\nWARNUNG: Du bist dabei, die Topologie deines Home-Verzeichnisses physisch zu verändern.\n")
		fmt.Printf("Möchtest du fortfahren? [y/N]: ")
		
		reader := bufio.NewReader(os.Stdin)
		response, _ := reader.ReadString('\n')
		response = strings.ToLower(strings.TrimSpace(response))

		if response == "y" || response == "yes" {
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
