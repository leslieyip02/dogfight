import Animation from "../entities/Animation";
import Asteroid from "../entities/Asteroid";
import type { EntityMap } from "../entities/Entity";
import Player from "../entities/Player";
import Powerup from "../entities/Powerup";
import Projectile from "../entities/Projectile";
import type { AsteroidEntityData, EntityData, PlayerEntityData, PowerupEntityData } from "../types/entity";
import type { DeltaEventData, SnapshotEventData } from "../types/event";
import type { Spritesheet } from "./sprites";

export function syncEntities(snapshot: SnapshotEventData | null, entities: EntityMap) {
  if (!snapshot) {
    return;
  }

  Object.values(snapshot.entities)
    .forEach(data => handleEntityData(data, entities));
}

export function mergeDeltas(current: DeltaEventData, next: DeltaEventData): DeltaEventData {
  const shouldOverwrite = current.timestamp < next.timestamp;
  Object.entries(next.updated)
    .forEach(entry => {
      const [id, data] = entry;
      if (!shouldOverwrite && current.updated[id]) {
        return;
      }
      current.updated[id] = data;
    });

  current.removed = [...current.removed, ...next.removed];
  current.timestamp = Math.max(current.timestamp, next.timestamp);
  return current;
}

export function removeEntities(delta: DeltaEventData, entities: EntityMap) {
  delta.removed
    .forEach(id => delete entities[id]);
}

export function updateEntities(delta: DeltaEventData, entities: EntityMap) {
  Object.entries(delta.updated)
    .filter(([id]) => !delta.removed.includes(id))
    .forEach(([,data]) => handleEntityData(data, entities));
}

export function handleEntityData(data: EntityData, entities: EntityMap) {
  const { id } = data;
  if (entities[id]) {
    entities[id].update(data);
    return;
  }

  switch (data.type) {
  case "asteroid":
    entities[id] = new Asteroid(data as AsteroidEntityData);
    break;

  case "player": {
    if (!Player.spritesheet) {
      // spritesheet hasn't loaded
      break;
    }

    entities[id] = new Player(data as PlayerEntityData);
    break;
  }

  case "projectile":
    entities[id] = new Projectile(data);
    break;

  case "powerup":
    entities[id] = new Powerup(data as PowerupEntityData);
    break;
  }
}

export function addAnimations(
  delta: DeltaEventData,
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
