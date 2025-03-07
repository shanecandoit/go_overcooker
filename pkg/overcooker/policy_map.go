package overcooker

import "math/rand"

// Policy is a map of (discrete) actions to probabilities
type Policy map[int]float32

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

// Update updates the policy based on the reward
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
