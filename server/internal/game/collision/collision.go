package collision

import (
	"math"
	"server/internal/game/entities"
	"server/internal/game/geometry"
	"slices"
)

// An Edge represents a horizontal boundary.
type Edge struct {
	id     *string // id of the entity it belongs to
	x      float64 // x-coordinate
	isLeft bool    // whether the edge is a left or right edge
}

// A CollisionHandler is a callback to handle a collision between entities with
// id1 and id2.
type CollisionHandler func(id1 *string, id2 *string)

// ResolveCollisionsLineSweep resolves collisions in entities in O(n log(n))
// time by using the line sweep algorithm. It maintains a window of entities
// with overlapping x-coordinates and only checks collisions within the window.
func ResolveCollisionsLineSweep(
	entities *map[string]entities.Entity,
	handleCollision CollisionHandler,
) {
	edges := getSortedEdges(entities)

	window := make(map[*string]bool)
	for _, edge := range edges {
		if edge.isLeft {
			e1 := (*entities)[*edge.id]

			for otherId := range window {
				e2 := (*entities)[*otherId]

				b1 := e1.GetBoundingBox()
				b2 := e2.GetBoundingBox()
				if !b1.DidCollide(b2) {
					continue
				}

				handleCollision(edge.id, otherId)
			}
			window[edge.id] = true
		} else {
			delete(window, edge.id)
		}
	}
}

// resolveCollisionsNaive resolves collisions in entities in O(n^2) time. It
// checks for collisions over all possible pairs of entities.
//
// It is only used for benchmarking.
func resolveCollisionsNaive(
	entities *map[string]entities.Entity,
	handleCollision CollisionHandler,
) {
	for id1 := range *entities {
		e1 := (*entities)[id1]
		for id2 := range *entities {
			if id1 == id2 {
				continue
			}

			e2 := (*entities)[id2]

			b1 := e1.GetBoundingBox()
			b2 := e2.GetBoundingBox()
			if !b1.DidCollide(b2) {
				continue
			}
			handleCollision(&id1, &id2)
		}
	}
}

// getSortedEdges returns all left and right edges in entities, ordered by
// x-coordinate. Left edges are ordered first in case of ties.
func getSortedEdges(entities *map[string]entities.Entity) []Edge {
	edges := make([]Edge, len(*entities)*2)

	i := 0
	for id, entity := range *entities {
		minX, maxX := entity.GetBoundingBox().HorizontalBounds()
		edges[i] = Edge{
			id:     &id,
			x:      minX,
			isLeft: true,
		}
		edges[i+1] = Edge{
			id:     &id,
			x:      maxX,
			isLeft: false,
		}
		i += 2
	}

	slices.SortFunc(edges, func(a Edge, b Edge) int {
		// If edges are effectively overlapping, order any left edges first
		// because left edges should only be ejected from the line sweep window
		// after checking for collisions with the other edge.
		if math.Abs(a.x-b.x) < geometry.EPSILON {
			if a.isLeft {
				return -1
			} else {
				return 1
			}
		}

		if math.Signbit(a.x - b.x) {
			return -1
		} else {
			return 1
		}
	})
	return edges
}
