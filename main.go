package main

import (
	"fmt"
	"math/rand"
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

		// update policy based on rewards
		for i, agent := range env.Agents {
			agentAction := actions_list[i]
			agentReward := rewards[i]
			agentPos := Position{X: agent.X, Y: agent.Y}
			agentPolicy := policyMap[agentPos]
			agentPolicy = agentPolicy.Update(agentPos, agentAction, agentReward)

			// update the policy map
			policyMap[agentPos] = agentPolicy
		}

		// every 15th step, spawn some items
		if step%15 == 0 {
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

	// Print the policy map
	fmt.Println("Policy Map:")
	for y := 0; y < env.Height+1; y++ {
		for x := 0; x < env.Width+1; x++ {
			pos := Position{X: x, Y: y}
			policy := policyMap[pos]
			fmt.Printf("Policy at %v: %2.2v\n", pos, policy)
		}
	}

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

func (p Policy) Update(position Position, action int, reward float32) Policy {
	// Update the policy based on the reward
	// agentPolicy = agentPolicy.Update(agentPos, agentAction, agentReward)

	// Create a copy of the policy to avoid modifying the original
	newPolicy := make(Policy)
	for a, prob := range p {
		newPolicy[a] = prob
	}

	// Learning parameters
	learningRate := float32(0.1) // How quickly to adapt to new rewards

	// Skip update for zero rewards to avoid reinforcing neutral actions
	if reward == 0 {
		return newPolicy
	}

	// For positive rewards: increase probability of the taken action
	// For negative rewards: decrease probability of the taken action

	// Get current probability for the action
	currentProb := newPolicy[action]

	// Calculate the adjustment (scaled by learning rate)
	// - For positive rewards: move probability toward 1.0
	// - For negative rewards: move probability toward 0.0
	var adjustment float32
	if reward > 0 {
		adjustment = learningRate * reward * (1.0 - currentProb)
	} else {
		adjustment = learningRate * reward * currentProb
		// Ensure we don't go below zero
		if currentProb+adjustment < 0.05 {
			adjustment = 0.05 - currentProb
		}
	}

	// Apply the adjustment
	newPolicy[action] += adjustment

	// Distribute the remaining probability mass among other actions
	// to ensure probabilities still sum to 1.0
	remainingProb := float32(1.0) - newPolicy[action]

	// Calculate total probability of other actions before normalization
	totalOtherProb := float32(0.0)
	for a := range newPolicy {
		if a != action {
			totalOtherProb += newPolicy[a]
		}
	}

	// Normalize other probabilities
	if totalOtherProb > 0 {
		for a := range newPolicy {
			if a != action {
				newPolicy[a] = newPolicy[a] * remainingProb / totalOtherProb
			}
		}
	} else {
		// If all other probabilities were 0, distribute evenly
		equalShare := remainingProb / float32(len(newPolicy)-1)
		for a := range newPolicy {
			if a != action {
				newPolicy[a] = equalShare
			}
		}
	}

	return newPolicy
}
