package overcooker

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
