package types

import "fmt"

const (
	// FlagDestroyed is a flag used to mark destroyed Cities
	FlagDestroyed string = "destroyed"
)

// City has a name and is connected to other Cities via Links
type City struct {
	Name  string
	Links map[string]*City
	Flags map[string]bool
}

// NewCity creates a City with a name and default flags
func NewCity(name string) City {
	// FlagDestroyed default is false
	return City{
		Name:  name,
		Links: make(map[string]*City),
		Flags: make(map[string]bool),
	}
}

func (c *City) AddLink(roadName string, connectingCity *City) {
	c.Links[roadName] = connectingCity
}

// IsDestroyed checks if City is destroyed
func (c *City) IsDestroyed() bool {
	return c.Flags[FlagDestroyed]
}

// Destroy City makes City burn in flames
func (c *City) Destroy() {
	c.Flags[FlagDestroyed] = true
}

// String representation for a City does not print destroyed linked Cities
func (c *City) String() string {
	var links string
	for roadName, connCity := range c.Links {
		// If other City destroyed print nothing
		if connCity.IsDestroyed() {
			continue
		}
		// If other City survived print Link
		links += fmt.Sprintf("%s=%s ", roadName, connCity.Name)
	}
	if len(links) == 0 {
		return c.Name
	}
	return fmt.Sprintf("%s %s", c.Name, links[:len(links)-1])
}
