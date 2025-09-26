import p5 from "p5";

import Player from "./entities/Player";
import { fetchSnapshot as fetchGameSnapshotData } from "../api/game";
import type {
  DeltaEventData,
  Event,
  EventType,
  JoinEventData,
  QuitEventData
} from "./types/event";
import type { EntityMap } from "./entities/Entity";
import { mergeDeltas, removeEntities, syncEntities, updateEntities } from "./utils/update";
import { drawBackground, drawEntities, drawMinimap } from "./utils/graphics";
import Input from "./utils/input";

const FPS = 60;

class Engine {
  instance: p5;

  clientId: string;
  entities: EntityMap;

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

    this.input = new Input(clientId, socket);
    this.delta = {
      "timestamp": 0,
      "updated": {},
      "removed": [],
    };
  }

  setup = () => {
    this.instance.createCanvas(window.innerWidth, window.innerHeight);
    this.instance.frameRate(FPS);
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

  // eslint-disable-next-line @typescript-eslint/no-unused-vars
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
    // TODO: restore some sort of removal animation
    updateEntities(this.delta, this.entities);
    removeEntities(this.delta, this.entities);

    this.delta.updated = {};
    this.delta.removed = [];
  };

  private syncGameState = async () => {
    await fetchGameSnapshotData()
      .then((snapshot) => syncEntities(snapshot, this.entities));
  };
};

export default Engine;
