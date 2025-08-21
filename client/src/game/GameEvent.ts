export type GameEventType = "join" | "quit" | "update";

export type GameEvent = {
  type: GameEventType
  data: GameJoinEventData | GameQuitEventData | GameUpdateEventData
};

export type GameJoinEventData = {
  clientId: string,
  username: string,
  x: number,
  y: number,
};

export type GameQuitEventData = {
  clientId: string,
};

export type GameUpdateEventData = {
  [clientId: string]: {
    x: number,
    y: number,
  }
};
