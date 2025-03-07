
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

## Status

Some success
After 5000 rounds of asking the policy and updating it based on rewards.

We see a high of 0.24

Number of steps: 5000
Final Environment:
Environment: &{env-1 [{a1 3 3 {s 5 4}} {a2 8 0 {o 1 1}}] [{s 5 4} {o 8 5} {p 6 4}] [{O1 4 1} {C1 9 1} {S1 9 5} {D1 5 5}] 9 5 map[onion_chop:4 onion_cook:4 onion_get:3 soup_deliver:6]}
. . . . . . . . a2.
. . . . O1. . . . C1
. . . . . . . . . .
. . . a1. . . . . .
. . . . . sp. . .
. . . . . D1. . oS1
EventCountsmap: map[onion_chop:4 onion_cook:4 onion_get:3 soup_deliver:6]
Policy Map:
Policy at {0 0}: map[00:0.16 01:0.13 02:0.21 03:0.21 04:0.13 05:0.17]
Policy at {1 0}: map[00:0.16 01:0.13 02:0.21 03:0.16 04:0.16 05:0.18]
Policy at {2 0}: map[00:0.15 01:0.13 02:0.21 03:0.16 04:0.17 05:0.17]
Policy at {3 0}: map[00:0.17 01:0.13 02:0.22 03:0.17 04:0.16 05:0.16]
Policy at {4 0}: map[00:0.16 01:0.13 02:0.22 03:0.17 04:0.16 05:0.17]
Policy at {5 0}: map[00:0.16 01:0.11 02:0.23 03:0.17 04:0.16 05:0.17]
Policy at {6 0}: map[00:0.14 01:0.12 02:0.24 03:0.16 04:0.18 05:0.16]
Policy at {7 0}: map[00:0.17 01:0.12 02:0.22 03:0.17 04:0.15 05:0.17]
Policy at {8 0}: map[00:0.17 01:0.12 02:0.24 03:0.15 04:0.16 05:0.17]
Policy at {9 0}: map[00:0.15 01:0.13 02:0.21 03:0.13 04:0.21 05:0.16]
Policy at {0 1}: map[00:0.16 01:0.17 02:0.16 03:0.21 04:0.13 05:0.17]
Policy at {1 1}: map[00:0.15 01:0.17 02:0.16 03:0.16 04:0.17 05:0.18]
Policy at {2 1}: map[00:0.16 01:0.17 02:0.17 03:0.18 04:0.15 05:0.16]
Policy at {3 1}: map[00:0.18 01:0.15 02:0.17 03:0.16 04:0.18 05:0.17]
Policy at {4 1}: map[00:0.17 01:0.17 02:0.17 03:0.16 04:0.16 05:0.17]
Policy at {5 1}: map[00:0.16 01:0.17 02:0.17 03:0.16 04:0.17 05:0.17]
Policy at {6 1}: map[00:0.17 01:0.17 02:0.16 03:0.17 04:0.15 05:0.17]
Policy at {7 1}: map[00:0.16 01:0.18 02:0.17 03:0.17 04:0.17 05:0.15]
Policy at {8 1}: map[00:0.17 01:0.17 02:0.15 03:0.18 04:0.17 05:0.16]
Policy at {9 1}: map[00:0.17 01:0.16 02:0.18 03:0.1 04:0.23 05:0.16]
Policy at {0 2}: map[00:0.17 01:0.17 02:0.18 03:0.21 04:0.13 05:0.15]
Policy at {1 2}: map[00:0.17 01:0.17 02:0.16 03:0.18 04:0.14 05:0.17]
Policy at {2 2}: map[00:0.16 01:0.16 02:0.19 03:0.15 04:0.18 05:0.17]
Policy at {3 2}: map[00:0.15 01:0.17 02:0.17 03:0.15 04:0.18 05:0.18]
Policy at {4 2}: map[00:0.15 01:0.17 02:0.17 03:0.17 04:0.17 05:0.17]
Policy at {5 2}: map[00:0.16 01:0.19 02:0.17 03:0.17 04:0.15 05:0.17]
Policy at {6 2}: map[00:0.15 01:0.18 02:0.17 03:0.16 04:0.17 05:0.18]
Policy at {7 2}: map[00:0.16 01:0.15 02:0.17 03:0.17 04:0.18 05:0.16]
Policy at {8 2}: map[00:0.18 01:0.17 02:0.17 03:0.16 04:0.16 05:0.16]
Policy at {9 2}: map[00:0.16 01:0.17 02:0.15 03:0.12 04:0.23 05:0.17]
Policy at {0 3}: map[00:0.17 01:0.16 02:0.16 03:0.21 04:0.13 05:0.17]
Policy at {1 3}: map[00:0.17 01:0.16 02:0.17 03:0.17 04:0.17 05:0.16]
Policy at {2 3}: map[00:0.16 01:0.17 02:0.16 03:0.16 04:0.17 05:0.17]
Policy at {3 3}: map[00:0.17 01:0.18 02:0.15 03:0.17 04:0.14 05:0.18]
Policy at {4 3}: map[00:0.17 01:0.18 02:0.16 03:0.15 04:0.17 05:0.17]
Policy at {5 3}: map[00:0.16 01:0.15 02:0.17 03:0.2 04:0.15 05:0.17]
Policy at {6 3}: map[00:0.18 01:0.16 02:0.19 03:0.14 04:0.17 05:0.17]
Policy at {7 3}: map[00:0.16 01:0.18 02:0.18 03:0.17 04:0.15 05:0.17]
Policy at {8 3}: map[00:0.18 01:0.16 02:0.17 03:0.16 04:0.18 05:0.16]
Policy at {9 3}: map[00:0.16 01:0.17 02:0.15 03:0.14 04:0.23 05:0.15]
Policy at {0 4}: map[00:0.17 01:0.16 02:0.16 03:0.21 04:0.13 05:0.17]
Policy at {1 4}: map[00:0.18 01:0.17 02:0.15 03:0.17 04:0.16 05:0.16]
Policy at {2 4}: map[00:0.17 01:0.16 02:0.17 03:0.16 04:0.18 05:0.17]
Policy at {3 4}: map[00:0.16 01:0.16 02:0.17 03:0.18 04:0.16 05:0.18]
Policy at {4 4}: map[00:0.17 01:0.16 02:0.16 03:0.16 04:0.18 05:0.16]
Policy at {5 4}: map[00:0.16 01:0.18 02:0.17 03:0.17 04:0.15 05:0.16]
Policy at {6 4}: map[00:0.18 01:0.17 02:0.15 03:0.17 04:0.17 05:0.17]
Policy at {7 4}: map[00:0.17 01:0.16 02:0.17 03:0.16 04:0.17 05:0.16]
Policy at {8 4}: map[00:0.16 01:0.15 02:0.18 03:0.18 04:0.17 05:0.16]
Policy at {9 4}: map[00:0.16 01:0.18 02:0.15 03:0.12 04:0.22 05:0.16]
Policy at {0 5}: map[00:0.16 01:0.21 02:0.13 03:0.21 04:0.13 05:0.16]
Policy at {1 5}: map[00:0.15 01:0.22 02:0.12 03:0.17 04:0.17 05:0.18]
Policy at {2 5}: map[00:0.16 01:0.22 02:0.12 03:0.17 04:0.18 05:0.16]
Policy at {3 5}: map[00:0.16 01:0.22 02:0.12 03:0.18 04:0.17 05:0.15]
Policy at {4 5}: map[00:0.17 01:0.22 02:0.11 03:0.17 04:0.17 05:0.16]
Policy at {5 5}: map[00:0.16 01:0.22 02:0.12 03:0.17 04:0.17 05:0.17]
Policy at {6 5}: map[00:0.17 01:0.22 02:0.12 03:0.16 04:0.16 05:0.17]
Policy at {7 5}: map[00:0.16 01:0.24 02:0.12 03:0.17 04:0.17 05:0.15]
Policy at {8 5}: map[00:0.17 01:0.22 02:0.12 03:0.15 04:0.16 05:0.18]
Policy at {9 5}: map[00:0.17 01:0.2 02:0.12 03:0.14 04:0.2 05:0.16]

## TODO

adding a gui to visualize the environment and agents
