import p5 from "p5";

import { fetchSnapshot as fetchGameSnapshotData } from "../api/game";
import {
  type Event,
  type Event_DeltaEventData,
  Event_JoinEventData,
  Event_QuitEventData,
  EventType,
} from "../pb/event";
import Audiosheet from "./audio/audio";
import type { EntityMap } from "./entities/Entity";
import Player from "./entities/Player";
import { type AnimationStep } from "./graphics/animation";
import type { CanvasConfig, GraphicsGameContext, GraphicsGUIContext } from "./graphics/context";
import {
  centerCanvas,
  drawAnimations,
  drawBackground,
  drawEntities,
  updateCanvasConfig,
} from "./graphics/game";
import {
  drawHUD,
  drawMinimap,
  drawRespawnPrompt,
} from "./graphics/gui";
import Spritesheet from "./graphics/sprites";
import Input from "./logic/input";
import {
  mergeDeltas,
  removeEntities,
  syncEntities,
  type UpdateContext,
  updateEntities,
} from "./logic/update";

const FPS = 60;

class Engine implements GraphicsGameContext, GraphicsGUIContext, UpdateContext {
  instance: p5;

  clientId: string;
  entities: EntityMap;

  input: Input;
  delta: Event_DeltaEventData;

  canvasConfig: CanvasConfig;
  foregroundAnimations: AnimationStep[];
  backgroundAnimations: AnimationStep[];

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
    this.foregroundAnimations = [];
    this.backgroundAnimations = [];

    this.input = new Input(clientId, socket);
    this.delta = {
      timestamp: 0,
      updated: [],
      removed: [],
    };
  }

  setup = async () => {
    this.instance.createCanvas(window.innerWidth, window.innerHeight);
    this.instance.frameRate(FPS);
    await Audiosheet.loadAll();
    await Spritesheet.loadAll(this.instance);
  };

  draw = () => {
    this.handleUpdates();
    this.handleInput();

    this.instance.push();
    centerCanvas(this);
    drawBackground(this);
    drawAnimations(this, this.backgroundAnimations);
    drawEntities(this);
    drawAnimations(this, this.foregroundAnimations);
    this.instance.pop();

    drawMinimap(this);
    drawHUD(this);
    drawRespawnPrompt(this);
  };

  mousePressed = () => {
    this.input.handleMousePress();
  };

  init = async () => {
    await this.syncGameState();
  };

  receive = (gameEvent: Event) => {
    switch (gameEvent.type) {
    case EventType.EVENT_TYPE_JOIN:
      this.handleJoin(gameEvent.joinEventData!);
      break;

    case EventType.EVENT_TYPE_QUIT:
      this.handleQuit(gameEvent.quitEventData!);
      break;

    case EventType.EVENT_TYPE_DELTA:
      this.handleDelta(gameEvent.deltaEventData!);
      break;

    default:
      return;
    }
  };

  getClientPlayer = () => {
    return this.entities[this.clientId] as Player;
  };

  getInput = () => {
    return this.input;
  };

  addAnimation = (animation: AnimationStep, isForeground: boolean) => {
    const target = isForeground ? this.foregroundAnimations : this.backgroundAnimations;
    target.push(animation);
  };

  // eslint-disable-next-line @typescript-eslint/no-unused-vars, unused-imports/no-unused-vars
  private handleJoin = (_data: Event_JoinEventData) => {
    // TODO: maybe log a chat message
    return;
  };

  // eslint-disable-next-line @typescript-eslint/no-unused-vars, unused-imports/no-unused-vars
  private handleQuit = (_data: Event_QuitEventData) => {
    // TODO: maybe log a chat message
  };

  private handleDelta = (data: Event_DeltaEventData) => {
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
    removeEntities(this);
    updateEntities(this);

    this.delta.updated = [];
    this.delta.removed = [];

    const clientPlayer = this.entities[this.clientId] as Player;
    if (clientPlayer) {
      updateCanvasConfig(this.canvasConfig, clientPlayer);
    }
  };

  private syncGameState = async () => {
    await fetchGameSnapshotData()
      .then((snapshot) => {
        syncEntities(snapshot, this);
      });
  };
};

export default Engine;
