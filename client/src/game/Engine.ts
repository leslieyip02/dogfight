import p5 from "p5";

import { fetchSnapshot as fetchGameSnapshotData } from "../api/game";
import type { EntityMap } from "./entities/Entity";
import Player from "./entities/Player";
import type {
  DeltaEventData,
  Event,
  EventType,
  JoinEventData,
  QuitEventData,
} from "./types/event";
import {
  type CanvasConfig,
  drawBackground,
  drawEntities,
  drawRespawnPrompt,
  updateCanvasConfig,
} from "./utils/graphics";
import { drawHUD, drawMinimap } from "./utils/gui";
import Input from "./utils/input";
import { loadSpritesheet, type Spritesheet } from "./utils/sprites";
import { addAnimations, mergeDeltas, removeEntities, syncEntities, updateEntities } from "./utils/update";

const FPS = 60;

class Engine {
  instance: p5;

  clientId: string;
  entities: EntityMap;
  canvasConfig: CanvasConfig;
  spritesheet: Spritesheet;

  input: Input;
  delta: DeltaEventData;

  constructor(
    instance: p5,
    clientId: string,
    socket: WebSocket,
  ) {
    this.instance = instance;
    this.instance.setup = this.setup;
    this.instance.draw = this.draw;
    this.instance.mousePressed = this.mousePressed;

    this.clientId = clientId;
    this.entities = {};
    this.canvasConfig = { x: 0.0, y: 0.0, zoom: 1.0 };
    this.spritesheet = {};

    this.input = new Input(clientId, socket);
    this.delta = {
      "timestamp": 0,
      "updated": {},
      "removed": [],
    };
  }

  setup = async () => {
    this.instance.createCanvas(window.innerWidth, window.innerHeight);
    this.instance.frameRate(FPS);
    this.spritesheet = await loadSpritesheet(this.instance);
    Player.spritesheet = this.spritesheet;
  };

  draw = () => {
    this.handleUpdates();
    this.handleInput();

    drawBackground(this.canvasConfig, this.instance);

    const clientPlayer = this.entities[this.clientId] as Player;
    drawEntities(this.canvasConfig, this.entities, this.instance);
    drawMinimap(this.canvasConfig, clientPlayer, this.entities, this.instance);
    drawHUD(clientPlayer, this.input, this.instance);

    if (!clientPlayer) {
      drawRespawnPrompt(this.instance);
    }
  };

  mousePressed = () => {
    this.input.handleMousePress();
  };

  init = async () => {
    await this.syncGameState();
  };

  receive = (gameEvent: Event) => {
    switch (gameEvent["type"] as EventType) {
    case "join":
      this.handleJoin(gameEvent.data as JoinEventData);
      break;

    case "quit":
      this.handleQuit(gameEvent.data as QuitEventData);
      break;

    case "delta":
      this.handleDelta(gameEvent.data as DeltaEventData);
      break;

    default:
      return;
    }
  };

  // eslint-disable-next-line @typescript-eslint/no-unused-vars, unused-imports/no-unused-vars
  private handleJoin = (_data: JoinEventData) => {
    // TODO: maybe log a chat message
    return;
  };

  // eslint-disable-next-line @typescript-eslint/no-unused-vars, unused-imports/no-unused-vars
  private handleQuit = (_data: QuitEventData) => {
    // TODO: maybe log a chat message
  };

  private handleDelta = (data: DeltaEventData) => {
    this.delta = mergeDeltas(this.delta, data);
  };

  private handleInput = () => {
    const clientPlayer = this.entities[this.clientId] as Player;
    if (!clientPlayer) {
      this.input.handleRespawn();
      return;
    }
    this.input.handleInput(this.instance);
  };

  private handleUpdates = () => {
    addAnimations(this.delta, this.entities, this.spritesheet);
    removeEntities(this.delta, this.entities);
    updateEntities(this.delta, this.entities);

    this.delta.updated = {};
    this.delta.removed = [];

    const clientPlayer = this.entities[this.clientId] as Player;
    if (clientPlayer) {
      updateCanvasConfig(this.canvasConfig, clientPlayer);
    }
  };

  private syncGameState = async () => {
    await fetchGameSnapshotData()
      .then((snapshot) => syncEntities(snapshot, this.entities));
  };
};

export default Engine;
