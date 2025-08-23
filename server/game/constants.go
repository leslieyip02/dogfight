package game

import "time"

const WIDTH = 10000.0
const HEIGHT = 10000.0

const FPS = 60
const FRAME_DURATION = time.Second / FPS

const ACCELERATION_DECAY = 2.0
const MAX_PLAYER_SPEED = 12.0
const PROJECTILE_SPEED = 16.0

const TURN_RATE_DECAY = 8.0
const MAX_TURN_RATE = 0.3

const PLAYER_RADIUS = 40.0
const PROJECTILE_RADIUS = 10.0
