package main

import (
	"fmt"

	// ov "github.com/shanecandoit/go_overcooker"
	ov "github.com/shanecandoit/go_overcooker/pkg/overcooker"
)

func main() {
	fmt.Println("Start")

	// Create a new environment
	env := ov.SimpleEnvironment()

	// Print the environment
	fmt.Println("Environment:", env)
	env.Render()

	// Create policies for agents
	policyMap := ov.NewPolicyMap(env)

	// Run simulation for N steps
	numSteps := 1000 * 5
	for step := 1; step <= numSteps; step++ {
		// Gather actions for each agent based on their policies
		actions_list := make([]int, len(env.Agents))
		prevPos := make([]ov.Position, len(env.Agents))
		for i, agent := range env.Agents {
			pos := ov.Position{X: agent.X, Y: agent.Y}
			prevPos[i] = pos
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
			agentPos := ov.Position{X: agent.X, Y: agent.Y}

			// Update current position policy
			agentPolicy := policyMap[agentPos]
			agentPolicy = agentPolicy.Update(agentPos, agentAction, agentReward)
			policyMap[agentPos] = agentPolicy

			// Check if agent moved (position changed)
			if prevPos[i].X != agentPos.X || prevPos[i].Y != agentPos.Y {
				prevDeducedAction := whichAction(prevPos[i], agentPos)

				// Get previous position's policy
				prevPolicy := policyMap[prevPos[i]]

				// Calculate a discounted reward for the previous action
				// This creates a smooth reward gradient that flows backward
				discountFactor := float32(0.5) // Adjustable parameter
				discountedReward := agentReward * discountFactor

				// Only backpropagate positive rewards to encourage positive behavior chains
				// If we didn't move or the reward is negative, don't update the policy
				if prevDeducedAction != ov.Act_None && discountedReward > 0 {
					prevPolicy = prevPolicy.Update(prevPos[i], prevDeducedAction, discountedReward)
					policyMap[prevPos[i]] = prevPolicy
				}
			}
		}

		// early on, spawn some items
		// every 15th step, spawn some items
		if step < 1000 && step%15 == 0 {
			env.EnvironmentSpawnRandomItemsForTraining()
		}

		// Display current state
		fmt.Printf("\nStep %d:\n", step)
		env.Render()

		// Optional: add a small delay to see each step
		// time.Sleep(500 * time.Millisecond)
	}

	// Print the environment
	fmt.Println("Number of steps:", numSteps)
	fmt.Println("Final Environment:")
	env.Render()
	fmt.Println("EventCountsmap:", env.EventCountsmap)

	// Print the policy map
	fmt.Println("Policy Map:")
	for y := 0; y < env.Height+1; y++ {
		for x := 0; x < env.Width+1; x++ {
			pos := ov.Position{X: x, Y: y}
			policy := policyMap[pos]
			fmt.Printf("Policy at %v: %2.2v\n", pos, policy)
		}
	}

}

func whichAction(prevPos, newPos ov.Position) int {
	if prevPos.X == newPos.X {
		if prevPos.Y > newPos.Y {
			return ov.Act_North
		} else {
			return ov.Act_South
		}
	} else {
		if prevPos.X > newPos.X {
			return ov.Act_West
		} else {
			return ov.Act_East
		}
	}
	return ov.Act_None
}

// Policy is a map of (discrete) actions to probabilities
// type Policy map[int]float32

// func NewPolicy() Policy {
// 	equalProb := float32(1.0 / (ov.Act_Interact + 1))
// 	// 16.67% for each action
// 	return Policy{
// 		ov.Act_None:     equalProb,
// 		ov.Act_North:    equalProb,
// 		ov.Act_South:    equalProb,
// 		ov.Act_East:     equalProb,
// 		ov.Act_West:     equalProb,
// 		ov.Act_Interact: equalProb,
// 	}
// }

// // PolicyMap a map of actions at every location
// type PolicyMap map[ov.Position]Policy

// // NewPolicyMap creates a new policy map
// func NewPolicyMap(env ov.Environment) PolicyMap {
// 	policyMap := PolicyMap{}
// 	for y := 0; y < env.Height+1; y++ {
// 		for x := 0; x < env.Width+1; x++ {
// 			policyMap[ov.Position{X: x, Y: y}] = NewPolicy()
// 		}
// 	}
// 	return policyMap
// }

// // GetActionBest returns an action based on the policy
// func (p ov.Policy) GetActionBest() int {
// 	// Implementation that selects an action based on probabilities
// 	// For now, just return the most probable action
// 	bestAction := ov.Act_None

// 	bestProb := float32(0.0)
// 	for action, prob := range p {
// 		if prob > bestProb {
// 			bestProb = prob
// 			bestAction = action
// 		}
// 	}
// 	return bestAction
// }
