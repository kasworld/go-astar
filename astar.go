package astar

import "container/heap"

// astar is an A* pathfinding implementation.

// Pather is an interface which allows A* searching on arbitrary objects which
// can represent a weighted graph.
type Pather interface {
	// PathNeighbors returns the direct neighboring nodes of this node which
	// can be pathed to.
	PathNeighbors() []Pather
	// PathNeighbourCost calculates the exact movement cost to neighbor nodes.
	PathNeighborCost(to Pather) float64
	// PathEstimatedCost is a heuristic method for estimating movement costs
	// between non-adjacent nodes.
	PathEstimatedCost(to Pather) float64
}

// node is a wrapper to store A* data for a Pather node.
type node struct {
	pather Pather
	cost   float64
	rank   float64
	parent *node
	open   bool
	closed bool
	index  int
}

// nodeMap is a collection of nodes keyed by Pather nodes for quick reference.
type nodeMap map[Pather]*node

// get gets the Pather object wrapped in a node, instantiating if required.
func (nm nodeMap) get(p Pather) *node {
	n, ok := nm[p]
	if !ok {
		n = &node{
			pather: p,
		}
		nm[p] = n
	}
	return n
}

// Path calculates a short path and the distance between the two Pather nodes.
//
// If no path is found, found will be false.
func Path(from, to Pather) (path []Pather, distance float64, found bool) {
	nm := nodeMap{}
	nq := &priorityQueue{}
	heap.Init(nq)
	fromNode := nm.get(from)
	fromNode.open = true
	heap.Push(nq, fromNode)
	for {
		if nq.Len() == 0 {
			// There's no path, return found false.
			return
		}
		current := heap.Pop(nq).(*node)
		current.open = false
		current.closed = true
		for _, neighbor := range current.pather.PathNeighbors() {
			cost := current.cost + current.pather.PathNeighborCost(neighbor)
			neighborNode := nm.get(neighbor)
			if neighbor == to {
				// Found a path to the goal.
				p := []Pather{}
				curr := neighborNode
				curr.parent = current
				for curr != nil {
					p = append(p, curr.pather)
					curr = curr.parent
				}
				return p, cost, true
			}
			if cost < neighborNode.cost {
				if neighborNode.open {
					heap.Remove(nq, neighborNode.index)
				}
				neighborNode.open = false
				neighborNode.closed = false
			}
			if !neighborNode.open && !neighborNode.closed {
				neighborNode.cost = cost
				neighborNode.open = true
				neighborNode.rank = cost + neighbor.PathEstimatedCost(to)
				neighborNode.parent = current
				heap.Push(nq, neighborNode)
			}
		}
	}
}

func Path2(from, to Pather, limit int, lenmax int) (path []Pather, count int) {
	nm := nodeMap{}
	nq := &priorityQueue{}
	heap.Init(nq)
	fromNode := nm.get(from)
	fromNode.open = true
	heap.Push(nq, fromNode)
	for {
		if nq.Len() == 0 {
			// There's no path, return found false.
			return nil, count
		}
		current := heap.Pop(nq).(*node)
		current.open = false
		current.closed = true
		for _, neighbor := range current.pather.PathNeighbors() {
			count++
			if count > limit {
				return nil, count
			}
			cost := current.cost + current.pather.PathNeighborCost(neighbor)
			neighborNode := nm.get(neighbor)
			if neighbor == to {
				// Found a path to the goal.
				p := []Pather{}
				curr := neighborNode
				curr.parent = current
				for plen := 0; curr != nil && plen < lenmax; plen++ {
					p = append(p, curr.pather)
					curr = curr.parent
				}
				return p, count
			}
			if cost < neighborNode.cost {
				if neighborNode.open {
					heap.Remove(nq, neighborNode.index)
				}
				neighborNode.open = false
				neighborNode.closed = false
			}
			if !neighborNode.open && !neighborNode.closed {
				neighborNode.cost = cost
				neighborNode.open = true
				neighborNode.rank = cost + neighbor.PathEstimatedCost(to)
				neighborNode.parent = current
				heap.Push(nq, neighborNode)
			}
		}
	}
}
