import type { PowerupAbility } from "../entities/Powerup";

export type EntityData = {
  type: EntityType,
  id: string,
  position: EntityPosition,
};

export type EntityType = "player" | "projectile" | "powerup";

export type EntityPosition = {
  x: number,
  y: number,
  theta: number,
};

export type PlayerEntityData = EntityData & {
  username: string;
};

export type PowerupEntityData = EntityData & {
  ability: PowerupAbility;
};
