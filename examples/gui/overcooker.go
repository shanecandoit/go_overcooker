package main

import (
	"fmt"
	"image"
	_ "image/png"
	"log"
	"os"
	"time"

	ov "github.com/shanecandoit/go_overcooker/pkg/overcooker"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

// Graphics assets:
// Items
// overcooker/examples/gui/art/onion_chopped_64x64.png
// overcooker/examples/gui/art/onion_raw_64x64.png
// overcooker/examples/gui/art/onion_soup_64x64.png
// Stations
// overcooker/examples/gui/art/station_onion_64x64.png
// overcooker/examples/gui/art/station_chop_64x64.png
// overcooker/examples/gui/art/station_cook_64x64.png
// overcooker/examples/gui/art/station_serve_64x64.png
// Agents
// overcooker/examples/gui/art/chef-1_64x64.png
// overcooker/examples/gui/art/chef-2_64x64.png

// Game implements ebiten.Game interface.
type Game struct {
	Environment ov.Environment
	Step        int
	// Create policies for agents
	PolicyMap map[ov.Position]ov.Policy
	Images    map[string]*ebiten.Image
	MaxSteps  int
}

// Update proceeds the game state.
// Update is called every tick (1/60 [s] by default).
func (g *Game) Update() error {

	// env := g.Environment

	// Gather actions for each agent based on their policies
	actions_list := make([]int, len(g.Environment.Agents))
	prevPos := make([]ov.Position, len(g.Environment.Agents))
	for i, agent := range g.Environment.Agents {
		pos := ov.Position{X: agent.X, Y: agent.Y}
		prevPos[i] = pos
		policy := g.PolicyMap[pos]
		actions_list[i] = policy.GetActionProba()
	}

	// Apply actions
	rewards, done := g.Environment.Step(actions_list)
	fmt.Println("Rewards:", rewards)
	fmt.Println("Done:", done)
	// if done {
	// 	break
	// }

	// update policy based on rewards
	for i, agent := range g.Environment.Agents {
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
	if g.Step < 1000 && g.Step%15 == 0 {
		g.Environment.EnvironmentSpawnRandomItemsForTraining()
	}

	// Display current state
	fmt.Printf("\nStep %d:\n", g.Step)
	g.Environment.Render()

	g.Step++

	// done?
	if g.MaxSteps > 0 && g.Step > g.MaxSteps {
		return fmt.Errorf("Max steps reached")
	}

	// limit the speed
	time.Sleep(100 * time.Millisecond)

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
	// Example: Draw a simple rectangle
	// screen.Fill(color.RGBA{0x80, 0x80, 0xc0, 0xff}) // light blue

	// Draw the environment
	for x := 0; x < g.Environment.Width+1; x++ {
		for y := 0; y < g.Environment.Height+1; y++ {
			// Draw stations
			station := g.Environment.GetStationAt(x, y)
			if station != nil {
				op := &ebiten.DrawImageOptions{}
				op.GeoM.Translate(float64(x*64), float64(y*64))
				// overcooker/examples/gui/art/station_chop_64x64.png
				// overcooker/examples/gui/art/station_cook_64x64.png
				// overcooker/examples/gui/art/station_serve_64x64.png
				switch station.Name[0:1] {
				case ov.StationOnion:
					screen.DrawImage(g.Images["station_onion_64x64.png"], op)
				case ov.StationChop:
					screen.DrawImage(g.Images["station_chop_64x64.png"], op)
				case ov.StationStove:
					screen.DrawImage(g.Images["station_cook_64x64.png"], op)
				case ov.StationDelivery:
					screen.DrawImage(g.Images["station_serve_64x64.png"], op)
				}
			}

			// Draw items
			item := g.Environment.GetItemAt(x, y)
			if item != nil {
				op := &ebiten.DrawImageOptions{}
				op.GeoM.Translate(float64(x*64), float64(y*64))
				switch item.Name {
				case ov.ItemOnionRaw:
					screen.DrawImage(g.Images["onion_raw_64x64.png"], op)
				case ov.ItemOnionChopped:
					screen.DrawImage(g.Images["onion_chopped_64x64.png"], op)
				case ov.ItemSoup:
					screen.DrawImage(g.Images["onion_soup_64x64.png"], op)
				}
			}
		}
	}

	// Draw the agents
	for i, agent := range g.Environment.Agents {
		x := agent.X
		y := agent.Y
		// ebitenutil.DebugPrintAt(screen, "A", x*64, y*64)
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(float64(x*64), float64(y*64))
		switch i % 2 {
		case 0:
			screen.DrawImage(g.Images["chef-1_64x64.png"], op)
		case 1:
			screen.DrawImage(g.Images["chef-2_64x64.png"], op)
		}
	}

	// Display total reward
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Total Reward: %.2f", g.Environment.TotalReward), 0, 0)

	// draw the step number
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Step: %d", g.Step), 0, 20)
	// draw number of things that happened
	eventCounts := g.Environment.EventCountsmap
	if eventCounts == nil {
		eventCounts = make(map[string]int)
	}
	onionGetCount := eventCounts["onion_get"]
	if onionGetCount > 0 {
		ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Onion Get: %d", onionGetCount), 0, 40)
	}
	onionChopCount := eventCounts["onion_chop"]
	if onionChopCount > 0 {
		ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Onion Chop: %d", onionChopCount), 0, 60)
	}
	onionCookCount := eventCounts["onion_cook"]
	if onionCookCount > 0 {
		ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Onion Cook: %d", onionCookCount), 0, 80)
	}
	soupDeliverCount := eventCounts["soup_deliver"]
	if soupDeliverCount > 0 {
		ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Soup Deliver: %d", soupDeliverCount), 0, 100)
	}

}

// Layout takes the outside size (e.g., the window size) and returns the (logical) screen size.
// If you don't have to adjust the screen size with the outside size, just return a fixed size.
func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 640, 480 // Increased screen size for better visibility
}

func (g *Game) loadImages() error {
	g.Images = make(map[string]*ebiten.Image)

	imagePaths := []string{
		// items
		"examples/gui/art/onion_chopped_64x64.png",
		"examples/gui/art/onion_raw_64x64.png",
		"examples/gui/art/onion_soup_64x64.png",
		// stations
		"examples/gui/art/station_chop_64x64.png",
		"examples/gui/art/station_cook_64x64.png",
		"examples/gui/art/station_serve_64x64.png",
		"examples/gui/art/station_onion_64x64.png",
		// agents
		"examples/gui/art/chef-1_64x64.png",
		"examples/gui/art/chef-2_64x64.png",
	}

	for _, path := range imagePaths {
		imgFile, err := os.Open(path)
		if err != nil {
			return fmt.Errorf("opening file %s: %w", path, err)
		}
		defer imgFile.Close()

		img, _, err := image.Decode(imgFile)
		if err != nil {
			return fmt.Errorf("decoding file %s: %w", path, err)
		}

		eImg := ebiten.NewImageFromImage(img)
		g.Images[path[len("examples/gui/art/"):]] = eImg // Store image with filename as key
	}

	return nil
}

func main() {
	game := &Game{}
	game.Environment = ov.SimpleEnvironment()
	// Specify the window size as you like. Here, a doubled size is specified.
	ebiten.SetWindowSize(800, 600)
	ebiten.SetWindowTitle("Overcooker")

	env := game.Environment
	game.Step = 1
	game.MaxSteps = 5000

	game.PolicyMap = ov.NewPolicyMap(env)

	// Load images
	if err := game.loadImages(); err != nil {
		log.Fatal("Error loading images:", err)
	}

	// Call ebiten.RunGame to start your game loop.
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
