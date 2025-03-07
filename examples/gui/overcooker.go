package main

import (
	"fmt"
	"log"

	ov "github.com/shanecandoit/go_overcooker/pkg/overcooker"

	"github.com/hajimehoshi/ebiten/v2"
)

// Game implements ebiten.Game interface.
type Game struct {
	Environment ov.Environment
	Step        int
	// Create policies for agents
	PolicyMap map[ov.Position]ov.Policy
}

// Update proceeds the game state.
// Update is called every tick (1/60 [s] by default).
func (g *Game) Update() error {

	env := g.Environment
	step := g.Step

	// Gather actions for each agent based on their policies
	actions_list := make([]int, len(env.Agents))
	prevPos := make([]ov.Position, len(env.Agents))
	for i, agent := range env.Agents {
		pos := ov.Position{X: agent.X, Y: agent.Y}
		prevPos[i] = pos
		policy := g.PolicyMap[pos]
		actions_list[i] = policy.GetActionProba()
	}

	// Apply actions
	rewards, done := env.Step(actions_list)
	fmt.Println("Rewards:", rewards)
	fmt.Println("Done:", done)
	// if done {
	// 	break
	// }

	// update policy based on rewards
	for i, agent := range env.Agents {
		agentAction := actions_list[i]
		agentReward := rewards[i]
		agentPos := ov.Position{X: agent.X, Y: agent.Y}

		// Update current position policy
		agentPolicy := g.PolicyMap[agentPos]
		agentPolicy = agentPolicy.Update(agentPos, agentAction, agentReward)
		g.PolicyMap[agentPos] = agentPolicy

		// Check if agent moved (position changed)
		if prevPos[i].X != agentPos.X || prevPos[i].Y != agentPos.Y {
			prevDeducedAction := whichAction(prevPos[i], agentPos)

			// Get previous position's policy
			prevPolicy := g.PolicyMap[prevPos[i]]

			// Calculate a discounted reward for the previous action
			// This creates a smooth reward gradient that flows backward
			discountFactor := float32(0.5) // Adjustable parameter
			discountedReward := agentReward * discountFactor

			// Only backpropagate positive rewards to encourage positive behavior chains
			// If we didn't move or the reward is negative, don't update the policy
			if prevDeducedAction != ov.Act_None && discountedReward > 0 {
				prevPolicy = prevPolicy.Update(prevPos[i], prevDeducedAction, discountedReward)
				g.PolicyMap[prevPos[i]] = prevPolicy
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

	g.Step++

	return nil
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

// Draw draws the game screen.
// Draw is called every frame (typically 1/60[s] for 60Hz display).
func (g *Game) Draw(screen *ebiten.Image) {
	// Write your game's rendering.
}

// Layout takes the outside size (e.g., the window size) and returns the (logical) screen size.
// If you don't have to adjust the screen size with the outside size, just return a fixed size.
func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 320, 240
}

func main() {
	game := &Game{}
	game.Environment = ov.SimpleEnvironment()
	// Specify the window size as you like. Here, a doubled size is specified.
	ebiten.SetWindowSize(800, 600)
	ebiten.SetWindowTitle("Overcooker")

	env := game.Environment
	game.Step = 1

	game.PolicyMap = ov.NewPolicyMap(env)

	// Call ebiten.RunGame to start your game loop.
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
