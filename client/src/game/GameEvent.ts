export type GameEventType = "join" | "quit" | "update" | "input";

export type GameEvent = {
  type: GameEventType
  data: GameJoinEventData | GameQuitEventData | GameUpdatePositionEventData | GameInputEventData
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
}

export type GameUpdatePositionEventData = {
  players: {
    [id: string]: EntityPosition,
  },
  projectiles: {
    [id: string]: EntityPosition,
  }
};

export type EntityStatus = {
  health: number,
}

export type GameUpdateStatusEventData = {
  [id: string]: EntityStatus,
}

export type GameInputEventData = {
  clientId: string,
  mouseX: number,
  mouseY: number,
  mousePressed: boolean,
}
