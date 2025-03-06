
# Policy Map

## Overview

A Policy Map is a spatial representation of decision-making strategies in a grid-based environment. It assigns action probabilities to each location, guiding agent behavior throughout the space.

## Core Concepts

### Structure

Each grid cell contains a policy (a probability distribution over possible actions)
Actions include: Move North, South, East, West; Interact; Wait
Policies determine what agents should do when standing in a particular location

### Visualization

Think of a floor with directional arrows showing the recommended action
Higher probability actions are displayed more prominently
The complete map reveals the overall strategy for navigating the environment

## Types of Policy Maps

1. Agent-specific Maps

    - Each agent maintains its own policy map
    - Allows for specialized roles and coordination
    - Maps evolve based on individual learning experiences

2. Environment Maps

    - A shared policy map for all agents in an environment
    - Provides baseline behaviors appropriate to the layout
    - Can be refined through collective experience

3. User-defined Maps

    - Designer-created policy maps for specific behaviors
    - Can be used to guide agents toward desired strategies
    - Serves as a starting point for further learning

## Applications

### Learning

- Policy maps evolve as agents learn from rewards and penalties
- High-reward actions gradually receive higher probabilities
- The map becomes a visual representation of learned knowledge

### Collaboration

- Maps can be shared between agents to transfer knowledge
- Maps can be blended to create consensus policies
- Different maps can assign complementary roles to agents

### Supervision

- A "Supervisor" can generate or modify policy maps
- In interactive settings, this could be the human player
- Maps provide an intuitive interface for guiding agent behavior

## Implementation

Policy maps are implemented as nested dictionaries mapping:

- Positions to policies
- Policies mapping actions to probabilities

This structure enables efficient lookup and modification while agents explore the environment.

    ```go
    // PolicyMap a map of actions at every location
    type PolicyMap map[Position]Policy

    type Position struct {
        X, Y int
    }

    // Policy is a map of (discrete) actions to probabilities
    type Policy map[int]float32
    ```

## Vague thoughts

- A policy map is a fixed size for a given env size, doesnt grow or get deeper as long as actions are fixed.
- Genetic Algorithms may work here
- PolicyMaps can be compared using naive agents
- Grow backward in time from rewards?
- probabilities are floats now but change to parts?
