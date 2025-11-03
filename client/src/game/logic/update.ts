import { type EntityData, EntityType } from "../../pb/entities";
import type { Event_DeltaEventData, Event_SnapshotEventData } from "../../pb/event";
import type Engine from "../Engine";
import Asteroid from "../entities/Asteroid";
import type { Entity, EntityMap } from "../entities/Entity";
import Player from "../entities/Player";
import Powerup from "../entities/Powerup";
import Projectile from "../entities/Projectile";
import { generatePlayerTrailAnimation } from "../graphics/animation";
import type { CanvasConfig } from "../graphics/game";

export function syncEntities(snapshot: Event_SnapshotEventData | null, game: Engine) {
  if (!snapshot) {
    return;
  }

  Object.values(snapshot.entities)
    .forEach(data => handleEntityData(data, game));
}

export function mergeDeltas(current: Event_DeltaEventData, next: Event_DeltaEventData): Event_DeltaEventData {
  const shouldOverwrite = current.timestamp < next.timestamp;
  next.updated
    .forEach(entity => {
      // TODO: is this necessary?
      if (!shouldOverwrite && current.updated.some(existing => existing.id === entity.id)) {
        return;
      }
      current.updated.push(entity);
    });

  current.removed = [...current.removed, ...next.removed];
  current.timestamp = Math.max(current.timestamp, next.timestamp);
  return current;
}

export function removeEntities(delta: Event_DeltaEventData, entities: EntityMap) {
  delta.removed
    .forEach(id => delete entities[id]);
}

export function updateEntities(game: Engine) {
  const { delta } = game;
  delta.updated
    .filter(entity => !delta.removed.includes(entity.id))
    .forEach(entity => handleEntityData(entity, game));
}

export function handleEntityData(data: EntityData, game: Engine) {
  const { id, type } = data;
  const { entities, canvasConfig } = game;

  if (entities[id]) {
    if (shouldUpdate(entities[id], canvasConfig)) {
      entities[id].update(data);
    }
    return;
  }

  switch (type) {
  case EntityType.ENTITY_TYPE_ASTEROID:
    entities[id] = new Asteroid(data);
    break;

  case EntityType.ENTITY_TYPE_PLAYER: {
    const player = new Player(data);
    entities[id] = player;
    const trailAnimation = generatePlayerTrailAnimation(new WeakRef(player));
    if (trailAnimation) {
      game.backgroundAnimations.push(trailAnimation);
    }
    break;
  }

  case EntityType.ENTITY_TYPE_POWERUP:
    entities[id] = new Powerup(data);
    break;

  case EntityType.ENTITY_TYPE_PROJECTILE:
    entities[id] = new Projectile(data);
    break;

  default:
    throw new Error(`unexpected entity ${data}`);
  }
}

function shouldUpdate(entity: Entity, canvasConfig: CanvasConfig) {
  return Math.abs(canvasConfig.x - entity.position.x) <= window.innerWidth
    && Math.abs(canvasConfig.y - entity.position.y) <= window.innerHeight;
}
