package simulation

import (
	"bufio"
	"errors"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"strings"

	"alien-invasion/simulation/types"
)

const (
	// RandAlienNameLen is a constant used to normalize number chosen as Alien name
	RandAlienNameLen = 10
)

//  Four corners of the World (cardinal directions)
var compass = map[string]string{
	"north": "south",
	"south": "north",
	"west":  "east",
	"east":  "west",
}

// WorldMapFile is a list of cities as read from the input file
type WorldMapFile []*City

// FilterDestroyed Cities from WorldMapFile
func (in WorldMapFile) FilterDestroyed(world types.World) WorldMapFile {
	out := make(WorldMapFile, 0, len(in))
	processed := make(map[string]bool)
	for _, city := range in {
		// If processed continue
		if processed[city.Name] {
			continue
		}
		// If not destroyed process
		if !city.IsDestroyed() {
			out = append(out, city)
			processed[city.Name] = true
			continue
		}
		// If destroyed process links (separated graph)
		for _, other := range city.Links {
			// Some links are already processed or destroyed
			if processed[other.Name] || other.IsDestroyed() {
				continue
			}
			// Other links we process
			out = append(out, other)
			processed[other.Name] = true
		}
	}
	return out
}

// String representation of a WorldMapFile, used to display output in same format as input
func (in WorldMapFile) String() string {
	var out string
	for _, city := range in {
		// If destroyed print nothing
		if city.IsDestroyed() {
			continue
		}
		// If survived print city name with links
		out += fmt.Sprintf("%s\n", city)
	}
	return out
}

// ReadWorldMapFile takes in a file and constructs a World map
func ReadWorldMapFile(file string) (World, WorldMapFile, error) {
	// Open and close file
	f, err := os.Open(file)
	if err != nil {
		return nil, nil, err
	}
	defer f.Close()

	// Init scanner to scan lines
	scanner := bufio.NewScanner(f)
	scanner.Split(bufio.ScanLines)

	// Prepare data structures
	w := make(World)
	input := make(WorldMapFile, 0)
	for scanner.Scan() {
		sections := strings.Split(scanner.Text(), " ")
		// Add new City to the world map
		city := w.AddNewCity(sections[0])
		// Connect nearby Cities
		for _, s := range sections[1:] {
			roadName, cityName, err := extractLink(s)
			if err != nil {
				return nil, nil, err
			}
			// Add new linked City if unknown
			other, exists := w[cityName]
			if !exists {
				other = w.AddNewCity(cityName)
			}
			city.AddLink(roadName, other)
			other.AddLink(compass[roadName], city)
			//fmt.Printf("RoadName : %s City : %s\n", roadName, other.Name)
		}
		input = append(input, city)
		fmt.Printf("Reading Line: %s\n", city)
	}

	return w, input, nil
}

// extractLink extracts a road and city name from input string or returns an error
func extractLink(input string) (string, string, error) {
	link := strings.Split(input, "=")
	if len(link) != 2 {
		return "", "", errors.New("Wrong link format")
	}
	roadName, cityName := link[0], link[1]
	return roadName, cityName, nil
}

// RandAliens creates N new Aliens with random names
func RandAliens(n int, r *rand.Rand, iterations int) []*Alien {
	aliens := make([]*Alien, n)
	for i := 0; i < n; i++ {
		name := "Alien " + strconv.Itoa(i+1)
		alien := types.NewAlien(name, iterations)
		aliens[i] = &alien
	}
	return aliens
}
