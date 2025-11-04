import p5 from "p5";

import { fetchSnapshot as fetchGameSnapshotData, sendEvent } from "../api/game";
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
import { type CanvasConfig, type GraphicsGameContext, type GraphicsGUIContext, initCanvasConfig } from "./graphics/context";
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
import { convertInputToEvent, handleMouseMove, handleMousePress, initInput, type Input } from "./logic/input";
import {
  initDelta,
  mergeDeltas,
  removeEntities,
  syncEntities,
  type UpdateContext,
  updateEntities,
} from "./logic/update";

const FPS = 60;

/**
 * Represents the game.
 * Encapsulates updates and rendering for all entities.
 * Implements context interfaces so that they can be passed to other functions
 * easily without having to pass the entire class.
 */
class Engine implements GraphicsGameContext, GraphicsGUIContext, UpdateContext {
  instance: p5;
  clientId: string;
  host: string;
  socket: WebSocket;

  entities: EntityMap;
  delta: Event_DeltaEventData;
  input: Input;

  canvasConfig: CanvasConfig;
  foregroundAnimations: AnimationStep[];
  backgroundAnimations: AnimationStep[];

  constructor(
    instance: p5,
    clientId: string,
    host: string,
    socket: WebSocket,
  ) {
    this.instance = instance;
    this.instance.setup = this.setup;
    this.instance.draw = this.draw;
    this.instance.mousePressed = this.mousePressed;

    this.clientId = clientId;
    this.host = host;
    this.socket = socket;

    this.entities = {};
    this.delta = initDelta();
    this.input = initInput();

    this.canvasConfig = initCanvasConfig();
    this.foregroundAnimations = [];
    this.backgroundAnimations = [];
  }

  // ==========================================================================
  //  p5.js handlers
  // ==========================================================================

  /**
   * See https://p5js.org/reference/p5/setup/.
   */
  setup = async () => {
    this.instance.createCanvas(window.innerWidth, window.innerHeight);
    this.instance.frameRate(FPS);
    await Audiosheet.loadAll();
    await Spritesheet.loadAll(this.instance);
  };

  /**
   * See https://p5js.org/reference/p5/draw/.
   */
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

  /**
   * See https://p5js.org/reference/p5/mousePressed/.
   */
  mousePressed = () => {
    this.input = handleMousePress(this.input);

    if (this.getClientPlayer()) {
      Audiosheet.get("shoot")?.play();
    }
  };

  // ==========================================================================
  //  WebSocket handlers
  // ==========================================================================

  /**
   * Initializes game state by fetching the latest snapshot.
   * This should be called once the WebSocket connection is opened.
   */
  init = async () => {
    await this.syncGameState();
  };

  /**
   * Triages event handling to handler functions.
   * @param event game data from the server
   */
  receive = (event: Event) => {
    switch (event.type) {
    case EventType.EVENT_TYPE_JOIN:
      this.handleJoin(event.joinEventData!);
      break;

    case EventType.EVENT_TYPE_QUIT:
      this.handleQuit(event.quitEventData!);
      break;

    case EventType.EVENT_TYPE_DELTA:
      this.handleDelta(event.deltaEventData!);
      break;

    default:
      return;
    }
  };

  // ==========================================================================
  //  Getters (provides context to other modules)
  // ==========================================================================

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

  // ==========================================================================
  //  Private handlers
  // ==========================================================================

  // eslint-disable-next-line @typescript-eslint/no-unused-vars, unused-imports/no-unused-vars
  private handleJoin = (_data: Event_JoinEventData) => {
    // TODO: maybe log a chat message
    return;
  };

  // eslint-disable-next-line @typescript-eslint/no-unused-vars, unused-imports/no-unused-vars
  private handleQuit = (_data: Event_QuitEventData) => {
    // TODO: maybe log a chat message
  };

  /**
   * Reconciles the incoming delta with the current delta
   * @param data incoming data
   */
  private handleDelta = (data: Event_DeltaEventData) => {
    this.delta = mergeDeltas(this.delta, data);
  };

  /**
   * Processes user input, and sends an input event to the server.
   * @param data incoming data
   */
  private handleInput = () => {
    this.input = handleMouseMove(this.input, this.instance);

    const clientPlayer = this.getClientPlayer();
    const [input, event] = convertInputToEvent(this.input, this.clientId, !clientPlayer);
    if (event) {
      sendEvent(this.socket, event);
    }
    this.input = input;
  };

  /**
   * Updates all entities.
   * Remove entities marked for removal.
   * If the client player exists, center the canvas around the player.
   */
  private handleUpdates = () => {
    removeEntities(this);
    updateEntities(this);

    this.delta.updated = [];
    this.delta.removed = [];

    const clientPlayer = this.getClientPlayer();
    if (clientPlayer) {
      updateCanvasConfig(this.canvasConfig, clientPlayer);
    }
  };

  /**
   * Fetches a snapshot of the game's state from the server and sets the local
   * state to match.
   */
  private syncGameState = async () => {
    await fetchGameSnapshotData(this.host)
      .then((snapshot) => {
        syncEntities(snapshot, this);
      });
  };
};

export default Engine;
