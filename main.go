package main

import "fmt"

func main() {
	fmt.Println("Start")

	// Create a new environment
	env := SimpleEnvironment()

	// Print the environment
	fmt.Println("Environment:", env)
	env.render()
}

// Agent is a person or a robot in the environment
type Agent struct {
	Name string
	X, Y int

	// agent inventory, only 1 object at a time
	Inventory Item
	// onion, tomato, lettuce, cheese, bread, patty
}

// Agent Actions
// move_up, move_right, move_down, move_left
// interact, nop
// N S E W I _
// N: North, S: South, E: East, W: West
// I: Interact, _: Wait No operation
const Act_None = 0
const Act_North = 1
const Act_South = 2
const Act_East = 3
const Act_West = 4
const Act_Interact = 5

// Environment is the world where agents interact
type Environment struct {
	Name   string
	Agents []Agent

	Items    []Item
	Stations []Station

	Width, Height int
}

// Item is a generic object in the environment
type Item struct {
	Name string
	X, Y int
}

const ItemOnionRaw = "o"
const ItemOnionChopped = "p"

// Station is a place where agents can interact with items
type Station struct {
	Name string
	X, Y int
}

type Position struct {
	X, Y int
}

// Stations {{
// Station for getting onions
const StationOnion = "O"

// Station for chopping
const StationChop = "C"

// Station for stove
const StationStove = "S"

// Station for delivery
const StationDelivery = "D"

// }}

func SimpleEnvironment() Environment {
	// a simple environment with 2 agents
	env := Environment{
		Name: "env-1",
		Agents: []Agent{
			Agent{Name: "a1", X: 1, Y: 1},
			Agent{Name: "a2", X: 1, Y: 4},
		},
	}

	// a box of onions
	env.Items = append(env.Items, Item{Name: ItemOnionRaw + "1", X: 2, Y: 4})
	env.Items = append(env.Items, Item{Name: ItemOnionRaw + "2", X: 4, Y: 2})

	// a box of onions
	env.Stations = append(env.Stations, Station{Name: StationOnion + "1", X: 4, Y: 1})
	// a chopping station
	env.Stations = append(env.Stations, Station{Name: StationChop + "1", X: 9, Y: 1})
	// a cooking station
	env.Stations = append(env.Stations, Station{Name: StationStove + "1", X: 9, Y: 5})

	// a delivery serving station
	env.Stations = append(env.Stations, Station{Name: StationDelivery + "1", X: 5, Y: 5})

	// set the width and height of the environment
	env.Width = 5  // min width
	env.Height = 5 // min height
	for _, agent := range env.Agents {
		if agent.X > env.Width {
			env.Width = agent.X
		}
		if agent.Y > env.Height {
			env.Height = agent.Y
		}
	}
	for _, ressource := range env.Items {
		if ressource.X > env.Width {
			env.Width = ressource.X
		}
		if ressource.Y > env.Height {
			env.Height = ressource.Y
		}
	}
	for _, station := range env.Stations {
		if station.X > env.Width {
			env.Width = station.X
		}
		if station.Y > env.Height {
			env.Height = station.Y
		}
	}

	return env
}

func (env *Environment) getAgentAt(x, y int) *Agent {
	for i := range env.Agents {
		if env.Agents[i].X == x && env.Agents[i].Y == y {
			return &env.Agents[i] // Returns pointer to actual agent
		}
	}
	return nil
}

func (env *Environment) getItemAt(x, y int) *Item {
	for i := range env.Items {
		if env.Items[i].X == x && env.Items[i].Y == y {
			return &env.Items[i]
		}
	}
	return nil
}

func (env *Environment) getStationAt(x, y int) *Station {
	for i := range env.Stations {
		if env.Stations[i].X == x && env.Stations[i].Y == y {
			return &env.Stations[i]
		}
	}
	return nil
}

func (env *Environment) render() {
	// each object may take 2 characters

	fmt.Println("Environment:", env)
	maxY := env.Height + 1
	maxX := env.Width + 1
	for y := 0; y < maxY; y++ {
		for x := 0; x < maxX; x++ {
			agent := env.getAgentAt(x, y)
			resource := env.getItemAt(x, y)
			station := env.getStationAt(x, y)
			if agent != nil {
				fmt.Print(agent.Name)
			} else if resource != nil {
				fmt.Print(resource.Name)
			} else if station != nil {
				fmt.Print(station.Name)
			} else {
				fmt.Print(". ")
			}
		}
		fmt.Println()
	}
}

// Step moves the environment forward by applying the given actions
func (env *Environment) Step(actions []int) {
	if len(actions) != len(env.Agents) {
		panic("Number of actions must match number of agents")
	}

	// Apply actions for each agent
	for i, action := range actions {
		agent := &env.Agents[i]

		// Handle movement
		newX, newY := agent.X, agent.Y
		switch action {
		case Act_North:
			newY--
		case Act_South:
			newY++
		case Act_East:
			newX++
		case Act_West:
			newX--
		case Act_Interact:
			env.handleInteraction(agent)
		}

		// Check if movement is valid
		if newX >= 0 && newX < env.Width+1 &&
			newY >= 0 && newY < env.Height+1 &&
			env.getAgentAt(newX, newY) == nil {
			agent.X, agent.Y = newX, newY
		}
	}
}

// handleInteraction processes an agent's attempt to interact
func (env *Environment) handleInteraction(agent *Agent) {
	// Check if agent is at a station
	station := env.getStationAt(agent.X, agent.Y)
	if station != nil {
		switch station.Name[0:1] {
		case StationOnion:
			// If agent doesn't have an item, give them an onion
			if agent.Inventory.Name == "" {
				agent.Inventory = Item{Name: ItemOnionRaw, X: -1, Y: -1} // -1 indicates in inventory
			}
		case StationChop:
			// If agent has an onion, chop it
			if agent.Inventory.Name == ItemOnionRaw {
				agent.Inventory.Name = ItemOnionChopped
			}
		}
		return
	}

	// Check if there's an item to pick up
	item := env.getItemAt(agent.X, agent.Y)
	if item != nil && agent.Inventory.Name == "" {
		agent.Inventory = *item
		// Remove the item from the environment
		for i, it := range env.Items {
			if it.X == item.X && it.Y == item.Y {
				env.Items = append(env.Items[:i], env.Items[i+1:]...)
				break
			}
		}
	} else if item == nil && agent.Inventory.Name != "" {
		// Drop the item
		droppedItem := agent.Inventory
		droppedItem.X, droppedItem.Y = agent.X, agent.Y
		env.Items = append(env.Items, droppedItem)
		agent.Inventory = Item{} // Reset inventory
	}
}

// Policy is a map of (discrete) actions to probabilities
type Policy map[int]float32

// PolicyMap a map of actions at every location
type PolicyMap map[Position]Policy

// GetAction returns an action based on the policy
func (p Policy) GetAction() int {
	// Implementation that selects an action based on probabilities
	// For now, just return the most probable action
	bestAction := Act_None
	bestProb := float32(0.0)
	for action, prob := range p {
		if prob > bestProb {
			bestProb = prob
			bestAction = action
		}
	}
	return bestAction
}
