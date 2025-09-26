import p5 from "p5";

import Player from "./entities/Player";
import { fetchSnapshot as fetchGameSnapshotData } from "../api/game";
import type {
  DeltaEventData,
  Event,
  EventType,
  InputEventData,
  JoinEventData,
  QuitEventData
} from "./types/event";
import type { Entity } from "./entities/Entity";
import { mergeDeltas, removeEntities, syncEntities, updateEntities } from "./utils/update";
import { drawBackground, drawEntities, drawMinimap } from "./utils/graphics";

const FPS = 60;

class Engine {
  instance: p5;
  zoom: number;

  clientId: string;
  entities: { [id: string]: Entity };
  delta: DeltaEventData;

  pressed: boolean;
  sendInput: (data: InputEventData) => void;

  constructor(
    instance: p5,
    clientId: string,
    sendInput: (data: InputEventData) => void
  ) {
    this.instance = instance;
    this.instance.setup = this.setup;
    this.instance.draw = this.draw;
    this.instance.mousePressed = this.mousePressed;

    this.zoom = 1.0;

    this.clientId = clientId;
    this.entities = {};

    this.delta = {
      "timestamp": 0,
      "updated": {},
      "removed": [],
    };

    this.pressed = false;
    this.sendInput = sendInput;
  }

  init = async () => {
    await this.syncGameState();
  };

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

    drawBackground(clientPlayer, this.zoom, this.instance);
    drawEntities(clientPlayer, this.entities, this.zoom, this.instance);
    drawMinimap(clientPlayer, this.entities, this.instance);
  };

  mousePressed = () => {
    this.pressed = true;
  };

  handleInput = () => {
    const mouseX = this.normalize(this.instance.mouseX, window.innerWidth);
    const mouseY = this.normalize(this.instance.mouseY, window.innerHeight);
    this.sendInput({
      id: this.clientId,
      mouseX,
      mouseY,
      mousePressed: this.pressed,
    });
    this.pressed = false;

    const throttle = Math.min(Math.sqrt(mouseX * mouseX + mouseY * mouseY), 1.0);
    if (throttle > 0.8) {
      this.zoom = Math.max(this.zoom - 0.005, 0.8);
    } else {
      this.zoom = Math.min(this.zoom + 0.005, 1.0);
    }
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

  // updates

  private syncGameState = async () => {
    await fetchGameSnapshotData()
      .then((snapshot) => syncEntities(snapshot, this.entities));
  };

  private handleUpdates = () => {
    // TODO: restore some sort of removal animation
    updateEntities(this.delta, this.entities);
    removeEntities(this.delta, this.entities);

    this.delta.updated = {};
    this.delta.removed = [];
  };

  // helpers

  private normalize = (value: number, full: number): number => {
    const delta = value - full / 2;
    return Math.sign(delta) * Math.min(Math.abs(delta / (full / 2 * 0.8)), 1.0);
  };
};

export default Engine;
