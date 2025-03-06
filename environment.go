package main

import (
	"fmt"
	"math/rand"
)

// Environment is the world where agents interact
type Environment struct {
	Name   string
	Agents []Agent

	Items    []Item
	Stations []Station

	Width, Height int

	// map of event counts, like achievements
	// maybe just for debugging
	EventCountsmap map[string]int
}

// Item is a generic object in the environment
type Item struct {
	Name string
	X, Y int
}

// Items
// o: onion, p: chopped onion, s: soup
// for v1 a soup is made from a single cooked, chopped onion
// for v1 no multiple ingredients, no burnt soup, no timed cooking
const ItemOnionRaw = "o"
const ItemOnionChopped = "p"
const ItemSoup = "s"

// Station is a place where agents can interact with items
type Station struct {
	Name string
	X, Y int
}

type Position struct {
	X, Y int
}

// Stations
const StationOnion = "O"    // Station for getting onions
const StationChop = "C"     // Station for chopping onions
const StationStove = "S"    // Station for stove
const StationDelivery = "D" // Station for delivery

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
func (env *Environment) Step(actions []int) (rewards []float32, done bool) {

	rewards = make([]float32, len(env.Agents))
	// right not no done?
	done = false

	// sanity check
	if len(actions) != len(env.Agents) {
		panic("Number of actions must match number of agents")
	}

	// Apply actions for each agent
	for i, action := range actions {
		agent := &env.Agents[i]
		reward := RewardStalling
		// default to a small negative reward w RewardStalling

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
			reward = env.handleInteraction(agent)
		}

		// Check if movement is valid
		if newX >= 0 && newX < env.Width+1 &&
			newY >= 0 && newY < env.Height+1 &&
			env.getAgentAt(newX, newY) == nil {
			agent.X, agent.Y = newX, newY
		} else {
			reward = RewardInvalidAction
		}

		// set rewards
		rewards[i] = float32(reward)

	}
	return rewards, done
}

func (env *Environment) CheckEventCountsmap() {
	if env.EventCountsmap == nil {
		env.EventCountsmap = make(map[string]int)
	}
}

// handleInteraction processes an agent's attempt to interact
func (env *Environment) handleInteraction(agent *Agent) float64 {

	// check env.EventCountsmap
	env.CheckEventCountsmap()

	// Check if agent is at a station
	station := env.getStationAt(agent.X, agent.Y)
	reward := RewardInvalidAction
	if station != nil {
		switch station.Name[0:1] {
		case StationOnion:
			// If agent doesn't have an item, give them an onion
			if agent.Inventory.Name == "" {
				agent.Inventory = Item{Name: ItemOnionRaw, X: -1, Y: -1} // -1 indicates in inventory
				reward = RewardOnionGet
				env.EventCountsmap["onion_get"]++
			}
		case StationChop:
			// If agent has an onion, chop it
			if agent.Inventory.Name == ItemOnionRaw {
				agent.Inventory.Name = ItemOnionChopped
				reward = RewardOnionChop
				env.EventCountsmap["onion_chop"]++
			}
		case StationStove:
			// If agent has a chopped onion, cook it
			if agent.Inventory.Name == ItemOnionChopped {
				agent.Inventory.Name = ItemSoup
				reward = RewardOnionCook
				env.EventCountsmap["onion_cook"]++
			}
		case StationDelivery:
			// If agent has a cooked soup, deliver it
			if agent.Inventory.Name == ItemSoup {
				agent.Inventory = Item{} // Reset inventory
				reward = RewardDeliverSoup
				env.EventCountsmap["soup_deliver"]++
			}
		}
		return reward
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
		reward = RewardPickup
	} else if item == nil && agent.Inventory.Name != "" {
		// NOTE: Dropping items is not allowed in v1
		// 	// Drop the item
		// 	droppedItem := agent.Inventory
		// 	droppedItem.X, droppedItem.Y = agent.X, agent.Y
		// 	env.Items = append(env.Items, droppedItem)
		// 	agent.Inventory = Item{} // Reset inventory
	}
	return reward
}

func (env *Environment) environmentSpawnRandomItemsForTraining() {

	// clean up junk
	env.Items = []Item{}
	rawOnionCount := 0
	rawOnionMax := 3
	choppedOnionCount := 0
	choppedOnionMax := 3
	soupCount := 0
	soupMax := 3
	for i := range env.Items {
		env.Items[i] = Item{}
		if env.Items[i].Name == ItemOnionRaw {
			rawOnionCount++
			if rawOnionCount > rawOnionMax {
				env.Items = append(env.Items[:i], env.Items[i+1:]...)
				fmt.Println("deleting raw onion")
			}
		}
		if env.Items[i].Name == ItemOnionChopped {
			choppedOnionCount++
			if choppedOnionCount > choppedOnionMax {
				env.Items = append(env.Items[:i], env.Items[i+1:]...)
				fmt.Println("deleting chopped onion")
			}
		}
		if env.Items[i].Name == ItemSoup {
			soupCount++
			if soupCount > soupMax {
				env.Items = append(env.Items[:i], env.Items[i+1:]...)
				fmt.Println("deleting soup")
			}
		}
	}

	listOfEmptyPositions := []Position{}
	for y := 0; y < env.Height+1; y++ {
		for x := 0; x < env.Width+1; x++ {
			if env.getAgentAt(x, y) == nil && env.getItemAt(x, y) == nil && env.getStationAt(x, y) == nil {
				listOfEmptyPositions = append(listOfEmptyPositions, Position{X: x, Y: y})
			}
		}
	}

	// shuffle the listOfEmptyPositions
	rand.Shuffle(len(listOfEmptyPositions), func(i, j int) {
		listOfEmptyPositions[i], listOfEmptyPositions[j] = listOfEmptyPositions[j], listOfEmptyPositions[i]
	})

	countEmptyPositions := len(listOfEmptyPositions)
	fmt.Println("countEmptyPositions:", countEmptyPositions)
	if countEmptyPositions < 3 {
		fmt.Println("Not enough empty positions")
		return
	}

	// print the first 5 positions
	fmt.Println("listOfEmptyPositions:", listOfEmptyPositions[:5])

	// we spawn some items, not stations
	// we want to learn to interact with things
	//
	// spawn a cooked soup
	// spawn a burnt soup
	// spawn a raw onion
	// spawn a chopped onion
	env.Items = append(env.Items, Item{Name: ItemSoup, X: listOfEmptyPositions[0].X, Y: listOfEmptyPositions[0].Y})
	env.Items = append(env.Items, Item{Name: ItemOnionRaw, X: listOfEmptyPositions[1].X, Y: listOfEmptyPositions[1].Y})
	env.Items = append(env.Items, Item{Name: ItemOnionChopped, X: listOfEmptyPositions[2].X, Y: listOfEmptyPositions[2].Y})

}

// RewardValues defines the point values for different actions
const (
	RewardPickup        = 0.1  // Small reward for picking up items
	RewardOnionGet      = 0.2  // Getting an onion from station
	RewardOnionChop     = 0.5  // Successfully chopping an onion
	RewardOnionCook     = 0.7  // Successfully cooking a chopped onion
	RewardDeliverSoup   = 1.0  // Delivering a finished soup
	RewardInvalidAction = -0.1 // Small penalty for invalid actions
	RewardStalling      = -0.1 // Small penalty for stalling

	// Maybe not:
	// RewardDrop       = 0.0  // Neutral for dropping items
	// no dropping in v1
)
