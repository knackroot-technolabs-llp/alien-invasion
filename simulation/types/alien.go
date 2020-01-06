package types

import (
	"fmt"
	//"alien-invasion/types"
)

const (
	// FlagDead is a flag used to mark dead Aliens
	FlagDead string = "dead"
)

// Alien can be dead or alive and occupying a City
type Alien struct {
	Name string
	city *City
	Flags map[string]bool
	RemainingIteration int
}

// NewAlien creates an Alien with a name and default flags
func NewAlien(name string, iterations int) Alien {
	// Flags FlagDead default is false
	return Alien{
		Name: name,
		city: nil,
		Flags: make(map[string]bool),
		RemainingIteration: iterations,
	}
}

// InvadeCity change City this Alien is occupying
func (a *Alien) InvadeCity(city *City) {
	a.city = city
}

// City returns City this Alien is occupying
func (a *Alien) City() *City {
	return a.city
}

// IsDead checks if Alien died
func (a *Alien) IsDead() bool {
	return a.Flags[FlagDead]
}

// Kill Alien makes it dead
func (a *Alien) Kill() {
	a.Flags[FlagDead] = true
}

// IsInvading checks if Alien is currently invading a City
func (a *Alien) IsInvading() bool {
	return a.city != nil
}

// IsTrapped checks if Alien is trapped in a City with no roads out
func (a *Alien) IsTrapped() bool {
	if !a.IsInvading() {
		return false
	}
	
	for _, c := range a.City().Links {
		if !c.IsDestroyed() {
			return false
		}
	}
	
	return true
}

// Are iteration over?
func (a *Alien) IsIterationOver() bool {
	return (a.RemainingIteration == 0)
}

// String representation for an Alien
func (a *Alien) String() string {
	return fmt.Sprintf("name=%s city={%s}\n", a.Name, a.city)
}
