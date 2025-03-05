
# over cooker

overcooked rl ascii style

multiple agent rl

## Agent

Each agent has a name which is a unique number "1", "2"

### Actions

Things an agent can do. In this game that means:

- move_up
- move_right
- move_down
- move_left
- interact
- nop

### Inventory

An agent can hold only a single thing at a time

- empty
- onion_raw
- onion_chopped
- soup_cooked

## Items

Items exist usually on a station, in the player inventory, maybe on the floor.

Things like:

- Onions
- Tomatos

## Stations

Places where tasks are performed

- OnionBox
- Chopping
- Deliver
- Stove

Interact with a station from any adjacent point

## Policy Map

A Policy Map is a spatially organized representation of an agent's policy, where each location in a discrete space is associated with a set of actions and their corresponding probabilities.
In essence, it's a grid (or tilemap) that dictates what action an agent should take when it occupies a particular cell.
Think of arrows on the floor, every square has at least one.

### Supervisor

A Supervisor dreams up policies. They can generate a random policy map.

## Plan

Right now, we only implement onion soup setups with multiple agents.

A basic environment looks like:

    . . . . . . . . . .
    . a1. . O1. . . . C1
    . . . . o2. . . . .
    . . . . . . . . . .
    . a2o1. . . . . . .
    . . . . . D1. . . S1

We see

- two agents "a1" and "a2"
- two onions on the floor "o1" and "o2"
- an onion box O1
- a chopping station C1
- a delivery point D1

## Current Implementation (v1)

In the current version, agents perform random actions without learning mechanisms.
The environment spawns random items every 5th step to create varied scenarios.
This random item spawning serves as a simple training curriculum, exposing agents to different situations they might encounter.
