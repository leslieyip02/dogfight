import type { EntityMap } from "../entities/Entity";
import Player from "../entities/Player";
import Powerup from "../entities/Powerup";
import Projectile from "../entities/Projectile";
import type { EntityData, PlayerEntityData, PowerupEntityData } from "../types/entity";
import type { DeltaEventData, SnapshotEventData } from "../types/event";

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
    .forEach(id => {
      entities[id].remove();
      delete entities[id];
    });
}

export function updateEntities(delta: DeltaEventData, entities: EntityMap) {
  Object.entries(delta.updated)
    .filter(([id]) => !delta.removed.includes(id))
    .forEach(([,data]) => handleEntityData(data, entities));
}

export function handleEntityData(data: EntityData, entities: EntityMap) {
  const { id, position } = data;
  if (entities[id]) {
    entities[id].update(position);
    return;
  }

  switch (data.type) {
  case "player":
    entities[id] = new Player(position, (data as PlayerEntityData).username);
    break;

  case "projectile":
    entities[id] = new Projectile(position);
    break;

  case "powerup":
    entities[id] = new Powerup(position, (data as PowerupEntityData).ability);
    break;

  default:
    break;
  }
}
