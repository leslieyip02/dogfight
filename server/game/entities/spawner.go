package entities

const (
	ASTEROID_SPAWN_INTERVAL = 60 * FPS
	POWERUP_SPAWN_INTERVAL  = 20 * FPS
	INITIAL_ASTEROID_COUNT  = 32
	INITIAL_POWERUP_COUNT   = 3
	RESET_INTERVAL          = 5 * 60 * FPS
)

type Spawner struct {
	counter int
}

func NewSpawner() Spawner {
	return Spawner{
		counter: 0,
	}
}

func (s *Spawner) InitEntities() []Entity {
	entities := []Entity{}

	for range INITIAL_ASTEROID_COUNT {
		asteroid, err := newRandomAsteroid()
		if err == nil {
			entities = append(entities, asteroid)
		}
	}
	for range INITIAL_POWERUP_COUNT {
		powerup, err := newRandomPowerup()
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
		asteroid, err := newRandomAsteroid()
		if err == nil {
			entities = append(entities, asteroid)
		}
	}
	if s.counter%POWERUP_SPAWN_INTERVAL == 0 {
		powerup, err := newRandomPowerup()
		if err == nil {
			entities = append(entities, powerup)
		}
	}

	s.counter = (s.counter + 1) % RESET_INTERVAL
	return entities
}
