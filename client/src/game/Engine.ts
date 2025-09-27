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
import { drawBackground, drawEntities, drawMinimap, loadSpritesheet, type Spritesheet } from "./utils/graphics";
import Input from "./utils/input";
import { addAnimations, mergeDeltas, removeEntities, syncEntities, updateEntities } from "./utils/update";

const FPS = 60;

class Engine {
  instance: p5;

  clientId: string;
  entities: EntityMap;

  input: Input;
  delta: DeltaEventData;

  spritesheet: Spritesheet;

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

    this.input = new Input(clientId, socket);
    this.delta = {
      "timestamp": 0,
      "updated": {},
      "removed": [],
    };

    this.spritesheet = {};
  }

  setup = async () => {
    this.instance.createCanvas(window.innerWidth, window.innerHeight);
    this.instance.frameRate(FPS);
    this.spritesheet = await loadSpritesheet(this.instance);
  };

  draw = () => {
    this.handleInput();
    this.handleUpdates();

    const clientPlayer = this.entities[this.clientId] as Player;
    if (!clientPlayer) {
      return;
    }

    const zoom = this.input.calculateZoom();
    drawBackground(clientPlayer, zoom, this.instance);
    drawEntities(clientPlayer, this.entities, zoom, this.instance);
    drawMinimap(clientPlayer, this.entities, this.instance);
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

  private handleQuit = (data: QuitEventData) => {
    delete this.entities[data.id];
  };

  private handleDelta = (data: DeltaEventData) => {
    this.delta = mergeDeltas(this.delta, data);
  };

  private handleInput = () => {
    this.input.handleInput(this.instance);
    this.input.calculateZoom();
  };

  private handleUpdates = () => {
    addAnimations(this.delta, this.entities, this.spritesheet);

    // TODO: replace with something more robust
    if (this.delta.removed.includes(this.clientId)) {
      this.delta.removed = this.delta.removed.filter(id => id !== this.clientId);
    }

    removeEntities(this.delta, this.entities);
    updateEntities(this.delta, this.entities);

    this.delta.updated = {};
    this.delta.removed = [];
  };

  private syncGameState = async () => {
    await fetchGameSnapshotData()
      .then((snapshot) => syncEntities(snapshot, this.entities));
  };
};

export default Engine;
