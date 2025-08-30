import type { PowerupType } from "./entities/Powerup";

export type GameEventType = "join" | "quit" | "position" | "powerup" | "input";

export type GameEvent = {
  type: GameEventType
  data: GameJoinEventData
    | GameQuitEventData
    | GameUpdatePositionEventData
    | GameUpdatePowerupEventData
    | GameInputEventData
};

export type GameJoinEventData = {
  id: string,
  username: string,
  position: EntityPosition
};

export type GameQuitEventData = {
  id: string,
};

export type EntityPosition = {
  x: number,
  y: number,
  theta: number,
};

export type GameUpdatePositionEventData = {
  players: {
    [id: string]: EntityPosition,
  },
  projectiles: {
    [id: string]: EntityPosition,
  },
};

export type GameUpdatePowerupEventData = {
  id: string,
  type: PowerupType
  active: boolean,
  position?: EntityPosition,
};

export type GameInputEventData = {
  clientId: string,
  mouseX: number,
  mouseY: number,
  mousePressed: boolean,
};
