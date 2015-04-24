package astar

import "container/heap"

func Path2(from, to Pather, trylimit int, lenmax int) (path []Pather, trycount int) {
	nm := nodeMap{}
	nq := &priorityQueue{}
	heap.Init(nq)
	fromNode := nm.get(from)
	fromNode.open = true
	heap.Push(nq, fromNode)
	for {
		if nq.Len() == 0 {
			// There's no path, return found false.
			return nil, trycount
		}
		current := heap.Pop(nq).(*node)
		current.open = false
		current.closed = true
		for _, neighbor := range current.pather.PathNeighbors() {
			trycount++
			if trycount > trylimit {
				return nil, trycount
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
				return p, trycount
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
