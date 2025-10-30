package collision

import (
	"math"
	"server/game/entities"
	"slices"
)

type Edge struct {
	id     *string
	x      float64
	isLeft bool
}

type CollisionHandler func(id1 *string, id2 *string)

func ResolveCollisions(
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
		if math.Signbit(a.x - b.x) {
			return -1
		} else {
			return 0
		}
	})
	return edges
}
