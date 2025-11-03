import { EntityData, EntityType } from "../../pb/entities";
import type { Event_DeltaEventData, Event_SnapshotEventData } from "../../pb/event";
import type { Vector } from "../../pb/vector";
import type Engine from "../Engine";
import Asteroid from "../entities/Asteroid";
import type { EntityMap } from "../entities/Entity";
import Player from "../entities/Player";
import Powerup from "../entities/Powerup";
import Projectile from "../entities/Projectile";
import { type AnimationStep, generatePlayerTrailAnimation } from "../graphics/animation";
import type { CanvasConfig } from "../graphics/context";

export interface UpdateContext {
  delta: Event_DeltaEventData;
  entities: EntityMap;
  canvasConfig: CanvasConfig;
  addAnimation: (animation: AnimationStep, isForeground: boolean) => void;
}

export function syncEntities(
  snapshot: Event_SnapshotEventData | null,
  game: Engine,
) {
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

export function removeEntities(context: UpdateContext) {
  const { delta, entities, addAnimation } = context;
  delta.removed
    .forEach(id => {
      const entity = entities[id];
      if (!entity) {
        return;
      }

      const animation = entity.remove();
      if (animation) {
        addAnimation(animation, true);
      }
      delete entities[id];
    });
}

export function updateEntities(context: UpdateContext) {
  const { delta } = context;
  delta.updated
    .filter(entityData => !delta.removed.includes(entityData.id))
    .forEach(entityData => handleEntityData(entityData, context));
}

export function handleEntityData(data: EntityData, context: UpdateContext) {
  const { id, type } = data;
  const { entities, canvasConfig, addAnimation } = context;

  if (entities[id]) {
    if (shouldCullEntity(entities[id].position, canvasConfig)) {
      return;
    }
    entities[id].update(data);
    return;
  }

  switch (type) {
  case EntityType.ENTITY_TYPE_ASTEROID:
    entities[id] = new Asteroid(data);
    break;

  case EntityType.ENTITY_TYPE_PLAYER: {
    const player = new Player(data);
    entities[id] = player;

    // TODO: move this somewhere else?
    const trailAnimation = generatePlayerTrailAnimation(new WeakRef(player));
    if (trailAnimation) {
      addAnimation(trailAnimation, false);
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

export function shouldCullEntity(position: Vector, canvasConfig: CanvasConfig): boolean {
  return Math.abs(canvasConfig.x - position.x) > window.innerWidth
    || Math.abs(canvasConfig.y - position.y) > window.innerHeight;
}
