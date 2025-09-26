import type { EntityData } from "./entity";

export type EventType = "join" | "quit" | "input" | "snapshot" | "delta";

export type Event = {
  type: EventType
  data: JoinEventData
    | QuitEventData
    | InputEventData
    | SnapshotEventData
    | DeltaEventData;
};

export type JoinEventData = {
  id: string,
  username: string,
};

export type QuitEventData = {
  id: string,
};

export type InputEventData = {
  id: string,
  mouseX: number,
  mouseY: number,
  mousePressed: boolean,
};

export type SnapshotEventData = {
  timestamp: number,
  entities: {
    [id: string]: EntityData,
  },
};

export type DeltaEventData = {
  timestamp: number,
  updated: {
    [id: string]: EntityData,
  },
  removed: string[],
};
