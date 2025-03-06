
# Agent Memory System Proposal

## Overview

We propose adding a memory system to agents in the Overcooker simulation that allows them to record, recall, and learn from interactions with the environment. This will enable agents to develop an understanding of game mechanics through experience.

## Core Concepts

1. Transformation Triples

    Each memory entry will be stored as a "transformation triple":

    `(Input Item, Station, Output Item)`

    Examples:

    - (Raw Onion, Chopping Station, Chopped Onion)
    - (Chopped Onion, Stove, Soup)
    - (Empty Hands, Onion Station, Raw Onion)

2. Success Tracking

    For each transformation attempt, we'll record:

    - Whether it succeeded
    - How many times it has been attempted
    - Success rate

3. Memory Structure

    ```go
    type Transformation struct {
        InputItem  string    // What the agent had in inventory
        Station    string    // What station was used
        OutputItem string    // What resulted from the interaction
        Success    bool      // Whether the transformation worked
    }

    type AgentMemory struct {
        Transformations []Transformation
        SuccessCount    map[string]int
        FailureCount    map[string]int
    }
    ```

## Usage Scenarios

### Learning Game Mechanics

An agent with raw onion interacts with different stations
After experimentation, it learns that:
Raw onion + Chopping Station → Chopped onion (100% success)
Raw onion + Stove → No change (0% success)

### Collaborative Learning

Agents can share discoveries
New agents can learn from experienced agents
Knowledge transfer accelerates learning

### Decision Making

Agents can use memories to:

- Choose optimal stations for desired transformations
- Plan multi-step recipes efficiently
- Predict outcomes of actions

Implementation Approach

1. Extend Agent struct to include memory
2. Modify handleInteraction to record transformations
3. Add query methods for memory retrieval
4. Implement memory visualization for debugging
5. Create mechanisms for memory sharing between agents

## Benefits

Emergent Behavior: Agents will naturally discover efficient workflows
Adaptation: Agents can adapt to changing environments
Visualization: Memories provide insight into agent learning process
Foundation for learning: This lays groundwork for reinforcement learning

## Next Steps

1. Implement basic memory tracking in Agent struct
2. Add recording logic in interaction handling
3. Develop simple visualization of agent knowledge
4. Test with different environment configurations
5. This can be expanded into something like a [Tuple Space](https://en.wikipedia.org/wiki/Tuple_space) or a [Blackboard](https://en.wikipedia.org/wiki/Blackboard_system)
6. Imagine showing blank slate agents a "training video" like sequence of Memories.

This memory system will provide a foundation for more sophisticated agent behaviors while maintaining the simplicity and clarity of the codebase.
