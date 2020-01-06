package simulation

import (
	"fmt"
	"math/rand"
	"sort"

	"alien-invasion/simulation/types"
	"alien-invasion/util"
)

type (
	// Alien attacking City
	Alien = types.Alien
	// City connected to other Cities
	City = types.City
	// World map of Cities
	World = types.World
)

// Aliens is a collection of all Aliens
type Aliens []*Alien

// AlienOccupation maps all Aliens by name
type AlienOccupation map[string]*Alien

// CityDefense maps Aliens by City
type CityDefense map[string]AlienOccupation

// Simulation struct represents a running simulation
type Simulation struct {
	// Simulation config
	R *rand.Rand
	// World state
	World
	Aliens
	CityDefense
}

// NoOpReason why Alien did not make a move
type NoOpReason uint8

// NoOpError when next move can not be made
type NoOpError struct {
	reason NoOpReason
}

const (
	// NoOpAlienDead when Alien is Dead
	NoOpAlienDead NoOpReason = iota
	// NoOpAlienTrapped when Alien is Trapped
	NoOpAlienTrapped
	// NoOpWorldDestroyed when World is destroyed
	NoOpWorldDestroyed
	// NoOpIterationsOver when iterations are over
	NoOpIterationsOver
	// NoOpMessage when no-op
	NoOpMessage = " || NO move! %s\n"
)

// Error string representation
func (err *NoOpError) Error() string {
	return fmt.Sprintf("Simulator no-op with reason: %d", err.reason)
}

// NewSimulation inits a new Simulation instance
func NewSimulation(r *rand.Rand, world World, aliens Aliens) Simulation {
	return Simulation{
		R:           r,
		World:       world,
		Aliens:      aliens,
		CityDefense: make(CityDefense),
	}
}

// Start the simulation
func (s *Simulation) Start() error {
	fmt.Printf("\nSimulation Started")
	fmt.Printf("\n\n")
	fmt.Print("----------\n")

	for s.shouldStop() {
		// Shuffle cards every iteration
		picks := util.ShuffleLen(len(s.Aliens), s.R)
		// Aliens make their moves
		noOpRound := true
		for _, p := range picks {
			if err := s.MoveAlien(s.Aliens[p]); err != nil {
				if _, ok := err.(*NoOpError); ok {
					// Alien made no move
					continue
				}
				return err
			}
			// If just one move is made, we continue the simulation
			noOpRound = false
		}
		// Check if last iteration was empty (no moves)
		if noOpRound {
			fmt.Printf("\n")
			fmt.Printf("Simulation ended early : No Available Moves")
			return nil
		}
	}
	// Game Over
	return nil
}

// MoveAlien moves the Alien position in the simulation
func (s *Simulation) MoveAlien(alien *Alien) error {
	from, to, err := s.pickMove(alien)
	fmt.Printf("\n\nMoving Alien: %s\n", alien.Name)
	fmt.Printf(" => From: %s\n", from)
	fmt.Printf(" => To: %s\n", to)
	if err != nil {
		// no-op error
		if noop, ok := err.(*NoOpError); ok {
			switch noop.reason {
			case NoOpWorldDestroyed:
				fmt.Printf(NoOpMessage, "World is destroyed.")
			case NoOpAlienDead:
				fmt.Printf(NoOpMessage, "Alien Dead.")
			case NoOpAlienTrapped:
				fmt.Printf(NoOpMessage, "Alien Trapped.")
			case NoOpIterationsOver:
				fmt.Printf(NoOpMessage, "Alien Iterations Over.")
			}
		}
		return err
	}
	// Move
	alien.RemainingIteration -= 1
	alien.InvadeCity(to)
	if from != nil {
		// Move from City
		delete(s.CityDefense[from.Name], alien.Name)
	}
	// Init city defense
	if s.CityDefense[to.Name] == nil {
		s.CityDefense[to.Name] = make(AlienOccupation)
	}
	// Move to City
	s.CityDefense[to.Name][alien.Name] = alien
	if len(s.CityDefense[to.Name]) > 1 {
		to.Destroy()
		// Kill Aliens and notify
		out := fmt.Sprintf(" || %s has been destroyed by ", to.Name)
		for _, a := range s.CityDefense[to.Name] {
			out += fmt.Sprintf(" %s and ", a.Name)
			a.Kill()
		}
		out = out[:len(out)-5] + "!\n"
		fmt.Print(out)
	}
	// Done
	return nil
}

// pickMove returns Alien move from City to City
func (s *Simulation) pickMove(alien *Alien) (*City, *City, error) {
	// Check if dead or trapped
	from := alien.City()
	if err := checkAlien(alien); err != nil {
		return from, nil, err
	}
	// At the beginning
	if from == nil {
		to := s.pickAnyCity()
		if to == nil {
			return from, to, &NoOpError{reason: NoOpWorldDestroyed}
		}
		return from, to, nil
	}
	// Move to next City
	to := s.pickConnectedCity(alien)
	return from, to, nil
}

// checkAlien returns NoOpError if Alien dead or trapped
func checkAlien(alien *Alien) *NoOpError {
	if alien.IsDead() {
		return &NoOpError{NoOpAlienDead}
	}
	if alien.IsTrapped() {
		return &NoOpError{NoOpAlienTrapped}
	}
	if alien.IsIterationOver() {
		return &NoOpError{NoOpIterationsOver}
	}
	return nil
}

// pickConnectedCity picks a random road to undestroyed City
func (s *Simulation) pickConnectedCity(alien *Alien) *City {
	// Nil if still not invading
	if !alien.IsInvading() {
		return nil
	}

	// Any undestroyed connected city
	for _, c := range alien.City().Links {
		if !c.IsDestroyed() {
			return c
		}
	}
	// No connected undestroyed City
	return nil
}

// pickAnyCity picks any undestroyed City in the World
func (s *Simulation) pickAnyCity() *City {
	// Any undestroyed city, pick deterministically
	var keys []string
	for k := range s.World {
		if c := s.World[k]; !c.IsDestroyed() {
			keys = append(keys, k)
		}
	}
	// If all Cities destroyed
	if len(keys) == 0 {
		return nil
	}
	// Sort keys for a deterministic pick
	sort.Strings(keys)
	pick := s.R.Intn(len(keys))
	return s.World[keys[pick]]
}

func (s *Simulation) shouldStop() bool {
	var allCityDestroyed = true
	for k := range s.World {
		if c := s.World[k]; !c.IsDestroyed() {
			allCityDestroyed = false
			break
		}
	}
	if allCityDestroyed {
		fmt.Printf("All cities are destroyed. Exiting !!\n")
	}
	var allAliensItrOver = true
	for _, a := range s.Aliens {
		if !a.IsIterationOver() {
			allAliensItrOver = false
			break
		}
	}
	if allAliensItrOver {
		fmt.Printf("Iterations of all aliens are over. Exiting !!\n")
	}

	return !(allCityDestroyed || allAliensItrOver)
}
