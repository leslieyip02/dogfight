import type p5 from "p5";

import type { EntityMap } from "../entities/Entity";
import type Player from "../entities/Player";
import type Input from "../logic/input";

export type CanvasConfig = {
  x: number;
  y: number;
  zoom: number;
};

export interface GraphicsGameContext {
  instance: p5;
  entities: EntityMap;
  canvasConfig: CanvasConfig;
}

export interface GraphicsGUIContext {
  instance: p5;
  entities: EntityMap;
  canvasConfig: CanvasConfig;
  getClientPlayer: () => Player | null;
  getInput: () => Input;
}
