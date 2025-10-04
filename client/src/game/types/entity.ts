import type { PowerupAbility } from "../entities/Powerup";
import type { Vector } from "./geometry";

export type EntityType = "asteroid" | "player" | "projectile" | "powerup";

export type EntityData = {
  type: EntityType,
  id: string,
  position: Vector,
  velocity: Vector,
  rotation: number,
};

export type AsteroidEntityData = EntityData & {
  points: Vector[],
}

export type PlayerEntityData = EntityData & {
  username: string,
};

export type PowerupEntityData = EntityData & {
  ability: PowerupAbility,
};
