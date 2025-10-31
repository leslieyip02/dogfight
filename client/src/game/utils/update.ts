import { type Entity as EntityData,EntityType } from "../../pb/entities";
import type { Event_DeltaEventData, Event_SnapshotEventData } from "../../pb/event";
import Animation from "../entities/Animation";
import Asteroid from "../entities/Asteroid";
import type { EntityMap } from "../entities/Entity";
import Player from "../entities/Player";
import Powerup from "../entities/Powerup";
import Projectile from "../entities/Projectile";
import { SOUNDS } from "./sounds";
import type { Spritesheet } from "./sprites";

export function syncEntities(snapshot: Event_SnapshotEventData | null, entities: EntityMap) {
  if (!snapshot) {
    return;
  }

  Object.values(snapshot.entities)
    .forEach(data => handleEntityData(data, entities));
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

export function updateEntities(delta: Event_DeltaEventData, entities: EntityMap) {
  delta.updated
    .filter(({ id }) => !delta.removed.includes(id))
    .forEach((entity) => handleEntityData(entity, entities));
}

export function handleEntityData(data: EntityData, entities: EntityMap) {
  const { id, type } = data;
  if (entities[id]) {
    entities[id].update(data);
    return;
  }

  switch (type) {
  case EntityType.ENTITY_TYPE_ASTEROID:
    entities[id] = new Asteroid(data);
    break;

  case EntityType.ENTITY_TYPE_PLAYER: {
    if (!Player.spritesheet) {
      // spritesheet hasn't loaded
      break;
    }

    entities[id] = new Player(data);
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

export function addAnimations(
  delta: Event_DeltaEventData,
  entities: EntityMap,
  spritesheet: Spritesheet,
) {
  // TODO: refactor
  delta.removed
    .forEach(id => {
      const animationName = entities[id]?.removalAnimationName();
      if (!animationName || !(animationName in spritesheet)) {
        return;
      }

      if (animationName in SOUNDS) {
        SOUNDS[animationName].play();
      }

      const animationId = `${id}-animation`;
      const animation = new Animation(
        entities[id].position,
        spritesheet[animationName],
        () => {
          delete entities[animationId];
        },
      );
      entities[animationId] = animation;
    });
}
