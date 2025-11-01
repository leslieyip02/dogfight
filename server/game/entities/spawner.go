package entities

import (
	"fmt"
	"math"
	"math/rand"
	"server/game/constants"
	"server/game/geometry"
	"server/id"
)

const (
	SPAWN_AREA_WIDTH  = 10000.0
	SPAWN_AREA_HEIGHT = 10000.0

	INITIAL_ASTEROID_COUNT = 32
	INITIAL_POWERUP_COUNT  = 3

	ASTEROID_SPAWN_INTERVAL = 60 * constants.FPS
	POWERUP_SPAWN_INTERVAL  = 20 * constants.FPS
	RESET_INTERVAL          = 5 * 60 * constants.FPS
)

// A Spawner is responsible for spawning new entities into the game. It used
// the current frame count to decide when to spawn a new entity.
type Spawner struct {
	counter int // frame count
}

func NewSpawner() Spawner {
	return Spawner{
		counter: 0,
	}
}

func (s *Spawner) SpawnPlayer(id string, username string) (*Player, error) {
	position := *geometry.NewRandomVector(
		0,
		0,
		SPAWN_AREA_WIDTH,
		SPAWN_AREA_HEIGHT,
	)
	velocity := *geometry.NewVector(0, 0)
	rotation := 0.0
	return newPlayer(id, position, velocity, rotation, username), nil
}

func (s *Spawner) spawnRandomAsteroid() (*Asteroid, error) {
	id, err := id.NewShortId()
	if err != nil {
		return nil, err
	}

	points := geometry.NewRandomConvexHull(
		ASTEROID_MIN_NUM_POINTS,
		ASTEROID_MAX_NUM_POINTS,
		ASTEROID_MIN_RADIUS,
		ASTEROID_MAX_RADIUS,
	)
	if geometry.HullArea(points) < ASTEROID_MIN_AREA {
		return nil, fmt.Errorf("too small")
	}

	position := *geometry.NewRandomVector(
		0,
		0,
		SPAWN_AREA_WIDTH,
		SPAWN_AREA_HEIGHT,
	)
	velocity := *geometry.NewRandomVector(
		0,
		0,
		ASTEROID_MAX_SPEED,
		ASTEROID_MAX_SPEED,
	)
	rotation := rand.Float64() * math.Pi * 2
	spin := rand.Float64()*ASTEROID_MAX_SPIN*2 - ASTEROID_MAX_SPIN
	return newAsteroid(id, position, velocity, rotation, &points, spin), nil
}

func (s *Spawner) spawnPowerup() (*Powerup, error) {
	id, err := id.NewShortId()
	if err != nil {
		return nil, err
	}

	position := *geometry.NewRandomVector(
		0,
		0,
		SPAWN_AREA_WIDTH,
		SPAWN_AREA_HEIGHT,
	)
	ability := newRandomAbility()
	return newPowerup(id, position, ability), nil
}

func (s *Spawner) InitEntities() []Entity {
	entities := []Entity{}

	for range INITIAL_ASTEROID_COUNT {
		asteroid, err := s.spawnRandomAsteroid()
		if err == nil {
			entities = append(entities, asteroid)
		}
	}
	for range INITIAL_POWERUP_COUNT {
		powerup, err := s.spawnPowerup()
		if err == nil {
			entities = append(entities, powerup)
		}
	}

	s.counter = (s.counter + 1) % RESET_INTERVAL
	return entities
}

func (s *Spawner) PollNewEntities() []Entity {
	entities := []Entity{}

	if s.counter%ASTEROID_SPAWN_INTERVAL == 0 {
		asteroid, err := s.spawnRandomAsteroid()
		if err == nil {
			entities = append(entities, asteroid)
		}
	}
	if s.counter%POWERUP_SPAWN_INTERVAL == 0 {
		powerup, err := s.spawnPowerup()
		if err == nil {
			entities = append(entities, powerup)
		}
	}

	s.counter = (s.counter + 1) % RESET_INTERVAL
	return entities
}
