package main

import (
	"fmt"
	"math/rand/v2"
)

func main() {
	fmt.Println("Start")

	// Create a new environment
	env := SimpleEnvironment()

	// Print the environment
	fmt.Println("Environment:", env)
	env.render()

	// Create policies for agents
	policyMap := NewPolicyMap(env)

	// Run simulation for N steps
	numSteps := 1000
	for step := 1; step <= numSteps; step++ {
		// Gather actions for each agent based on their policies
		actions_list := make([]int, len(env.Agents))
		for i, agent := range env.Agents {
			pos := Position{X: agent.X, Y: agent.Y}
			policy := policyMap[pos]
			actions_list[i] = policy.GetActionProba()
		}

		// Apply actions
		rewards, done := env.Step(actions_list)
		fmt.Println("Rewards:", rewards)
		fmt.Println("Done:", done)
		if done {
			break
		}

		// every 5th step, spawn some items
		if step%5 == 0 {
			env.environmentSpawnRandomItemsForTraining()
		}

		// Display current state
		fmt.Printf("\nStep %d:\n", step)
		env.render()

		// Optional: add a small delay to see each step
		// time.Sleep(500 * time.Millisecond)
	}

	// Print the environment
	fmt.Println("Final Environment:")
	env.render()
	fmt.Println("EventCountsmap:", env.EventCountsmap)

}

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

// Policy is a map of (discrete) actions to probabilities
type Policy map[int]float32

func NewPolicy() Policy {
	equalProb := float32(1.0 / (Act_Interact + 1))
	// 16.67% for each action
	return Policy{
		Act_None:     equalProb,
		Act_North:    equalProb,
		Act_South:    equalProb,
		Act_East:     equalProb,
		Act_West:     equalProb,
		Act_Interact: equalProb,
	}
}

// PolicyMap a map of actions at every location
type PolicyMap map[Position]Policy

// NewPolicyMap creates a new policy map
func NewPolicyMap(env Environment) PolicyMap {
	policyMap := PolicyMap{}
	for y := 0; y < env.Height+1; y++ {
		for x := 0; x < env.Width+1; x++ {
			policyMap[Position{X: x, Y: y}] = NewPolicy()
		}
	}
	return policyMap
}

// GetActionBest returns an action based on the policy
func (p Policy) GetActionBest() int {
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

// GetActionProba returns an action based on the policy
func (p Policy) GetActionProba() int {
	r := rand.Float32() // Random number between 0 and 1
	cumulative := float32(0.0)

	for action, prob := range p {
		cumulative += prob
		if r <= cumulative {
			return action
		}
	}

	// Fallback (should not reach here if probabilities sum to 1)
	return Act_None
}
