import { type EntityData, EntityType } from "../../pb/entities";
import type { Event_DeltaEventData, Event_SnapshotEventData } from "../../pb/event";
import Audiosheet from "../audio/audio";
import Animation from "../entities/Animation";
import Asteroid from "../entities/Asteroid";
import type { Entity, EntityMap } from "../entities/Entity";
import Player from "../entities/Player";
import Powerup from "../entities/Powerup";
import Projectile from "../entities/Projectile";
import type { CanvasConfig } from "../graphics/game";
import Spritesheet from "../graphics/sprites";

export function syncEntities(snapshot: Event_SnapshotEventData | null, entities: EntityMap, canvasConfig: CanvasConfig) {
  if (!snapshot) {
    return;
  }

  Object.values(snapshot.entities)
    .forEach(data => handleEntityData(data, entities, canvasConfig));
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

export function updateEntities(delta: Event_DeltaEventData, entities: EntityMap, canvasConfig: CanvasConfig) {
  delta.updated
    .filter(entity => !delta.removed.includes(entity.id))
    .forEach(entity => handleEntityData(entity, entities, canvasConfig));
}

export function handleEntityData(data: EntityData, entities: EntityMap, canvasConfig: CanvasConfig) {
  const { id, type } = data;

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

  case EntityType.ENTITY_TYPE_PLAYER:
    entities[id] = new Player(data);
    break;

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

export function addAnimations(
  delta: Event_DeltaEventData,
  entities: EntityMap,
) {
  // TODO: refactor
  delta.removed
    .forEach(id => {
      const animationName = entities[id]?.removalAnimationName();
      if (!animationName) {
        return;
      }

      Audiosheet.get(animationName)?.play();

      const animationId = `${id}-animation`;
      const frames = Spritesheet.getAnimationFrames(animationName);
      if (!frames) {
        return;
      }

      const animation = new Animation(
        entities[id].position,
        frames,
        () => {
          delete entities[animationId];
        },
      );
      entities[animationId] = animation;
    });
}

function shouldUpdate(entity: Entity, canvasConfig: CanvasConfig) {
  return Math.abs(canvasConfig.x - entity.position.x) <= window.innerWidth
    && Math.abs(canvasConfig.y - entity.position.y) <= window.innerHeight;
}
