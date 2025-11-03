import { EntityData, EntityType } from "../../pb/entities";
import type {
  Event_DeltaEventData,
  Event_SnapshotEventData,
} from "../../pb/event";
import type { Vector } from "../../pb/vector";
import Asteroid from "../entities/Asteroid";
import type { EntityMap } from "../entities/Entity";
import Player from "../entities/Player";
import Powerup from "../entities/Powerup";
import Projectile from "../entities/Projectile";
import {
  type AnimationStep,
  generatePlayerTrailAnimation,
} from "../graphics/animation";
import type { CanvasConfig } from "../graphics/context";

/**
 * The context to update.
 * Contains a subset of fields from Engine to avoid passing the whole Engine.
 */
export interface UpdateContext {
  delta: Event_DeltaEventData;
  entities: EntityMap;
  canvasConfig: CanvasConfig;
  addAnimation: (animation: AnimationStep, isForeground: boolean) => void;
}

/**
 * Forces the context to match the given snapshot.
 * @param context the context to update
 */
export function syncEntities(
  snapshot: Event_SnapshotEventData | null,
  context: UpdateContext,
) {
  if (!snapshot) {
    return;
  }

  Object.values(snapshot.entities)
    .forEach(data => handleEntityData(data, context));
}

export function initDelta(): Event_DeltaEventData {
  return {
    timestamp: 0,
    updated: [],
    removed: [],
  };
}

/**
 * Combines two deltas. If next is more recent, it will overwrite current.
 * @param current existing delta
 * @param next incoming delta
 */
export function mergeDeltas(
  current: Event_DeltaEventData,
  next: Event_DeltaEventData,
): Event_DeltaEventData {
  const shouldOverwrite = current.timestamp < next.timestamp;
  next.updated
    .forEach(entity => {
      // TODO: maybe the membership check can be rewritten
      const shouldAdd = shouldOverwrite
        || !current.updated.some(existing => existing.id === entity.id);
      if (!shouldAdd) {
        return;
      }
      current.updated.push(entity);
    });

  current.removed = [...current.removed, ...next.removed];
  current.timestamp = Math.max(current.timestamp, next.timestamp);
  return current;
}

/**
 * Removes entities that are marked for removal within the given context.
 * Will add a removal animation to the context if needed.
 * @param context the context to update
 */
export function removeEntities(context: UpdateContext) {
  const { delta, entities, canvasConfig, addAnimation } = context;
  delta.removed
    .forEach(id => {
      const entity = entities[id];
      if (!entity) {
        return;
      }

      if (!shouldCullEntity(entity.position, canvasConfig)) {
        const animation = entity.onRemove();
        if (animation) {
          addAnimation(animation, true);
        }
      }
      delete entities[id];
    });
}

/**
 * Updates all entities within the given context.
 * @param context the context to update
 */
export function updateEntities(context: UpdateContext) {
  const { delta } = context;
  delta.updated
    .filter(entityData => !delta.removed.includes(entityData.id))
    .forEach(entityData => handleEntityData(entityData, context));
}

/**
 * Updates an entity based on the given data.
 * Creates the entity if it doesn't exist in the given context.
 * @param data new data
 * @param context the context to update
 */
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

/**
 * Check if an entity should be updated/drawn.
 * Uses the window height and screen as boundaries.
 * @param position location of the entity
 * @param canvasConfig canvas origin
 */
export function shouldCullEntity(
  position: Vector,
  canvasConfig: CanvasConfig,
): boolean {
  return Math.abs(canvasConfig.x - position.x) > window.innerWidth
    || Math.abs(canvasConfig.y - position.y) > window.innerHeight;
}
