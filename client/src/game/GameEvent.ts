export type GameEventType = "join" | "quit" | "update" | "input";

export type GameEvent = {
  type: GameEventType
  data: GameJoinEventData | GameQuitEventData | GameUpdateEventData
};

export type GameJoinEventData = {
  clientId: string,
  username: string,
  position: EntityPosition
};

export type GameQuitEventData = {
  clientId: string,
};

export type EntityPosition = {
  x: number,
  y: number,
  theta: number,
}

export type GameUpdateEventData = {
  [id: string]: EntityPosition | null,
};

export type GameInputEventData = {
  clientId: string,
  mouseX: number,
  mouseY: number,
  mousePressed: boolean,
}
