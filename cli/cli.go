package cli

import (
	"alien-invasion/simulation"
	"errors"
	"flag"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"
)

const (
	// DefaultIterations used if number of iteration is not otherwise specified
	DefaultIterations int = 10000
	// DefaultNumberOfAliens used if number of Aliens is not otherwise specified
	DefaultNumberOfAliens int = 10
	// DefaultworldMapFile used if World file is not otherwise specified
	DefaultworldMapFile = "./data/map1.txt"
)

var (
	iterations, alienNumber int
	worldMapFile            string
)

// init cli flags
func init() {
	flag.IntVar(&iterations, "iterations", DefaultIterations, "number of iterations")
	flag.IntVar(&alienNumber, "aliens", DefaultNumberOfAliens, "number of aliens invading")
	flag.StringVar(&worldMapFile, "world", DefaultworldMapFile, "a file used as world map input")
	flag.Parse()
}

// validateInput validates input flags
func validateInput() error {
	if alienNumber <= 0 {
		return errors.New("Aliens number must be > 0")
	}
	if iterations <= 0 {
		return errors.New("Iterations number must be > 0")
	}
	if len(worldMapFile) == 0 {
		return errors.New("World map file path is empty")
	}
	return nil
}

// Execute command and set flags appropriately.
func Execute() {
	// Check input for errors
	if err := validateInput(); err != nil {
		fmt.Printf("Error while validating input : %s\n", err)
		flag.Usage()
		os.Exit(1)
	}
	// Read input file
	fmt.Printf("\nReading world map from file : %s\n", worldMapFile)
	world, in, err := simulation.ReadWorldMapFile(worldMapFile)
	if err != nil {
		fmt.Printf("Could not read world from map file \"%s\" with error: %s\n", worldMapFile, err)
		os.Exit(1)
	}

	// Build a simulator
	fmt.Printf("\n\nGenerating %d random Aliens...\n", alienNumber)

	r := buildRand()
	aliens := simulation.RandAliens(alienNumber, r, iterations)

	sim := simulation.NewSimulation(r, world, aliens)
	// Start the simulation and print any errors
	if err := sim.Start(); err != nil {
		fmt.Printf(formatMessage("Error while running simulation: %s"), err)
		os.Exit(1)
	}
	// Success
	fmt.Printf(formatMessage("Simulation Success"))

	fmt.Printf(formatMessage("Final World Map"))
	fmt.Print(in.FilterDestroyed(world))
}

// buildRand build a pseudorandom numbers generator from input flags
func buildRand() *rand.Rand {
	var seed = time.Now().UnixNano()
	var source = rand.NewSource(seed)
	return rand.New(source)
}

// formatMessage formats an important message ;)
func formatMessage(msg string) string {
	line := strings.Repeat("=", len(msg))
	out := fmt.Sprintf("\n\n")
	out += fmt.Sprintf("%s\n", line)
	out += fmt.Sprintf("%s\n", msg)
	out += fmt.Sprintf("%s\n", line)
	return out
}
