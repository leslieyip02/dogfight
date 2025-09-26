import p5 from "p5";

import Player from "./entities/Player";
import Minimap from "./Minimap";
import Powerup from "./entities/Powerup";
import Projectile from "./entities/Projectile";
import { fetchGameSnapshot as fetchGameSnapshotData } from "../api/game";
import type {
  DeltaEventData,
  Event,
  EventType,
  InputEventData,
  JoinEventData,
  QuitEventData
} from "./types/event";
import type { Entity } from "./entities/Entity";
import type { EntityData, PlayerEntityData, PowerupEntityData } from "./types/entity";

const DEBUG = import.meta.env.VITE_DEBUG;

const FPS = 60;
const GRID_SIZE = 96;

export const BACKGROUND_COLOR = "#111111";

class Engine {
  instance: p5;
  zoom: number;

  clientId: string;
  entities: { [id: string]: Entity };
  delta: DeltaEventData;

  minimap: Minimap;

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

    this.minimap = new Minimap();

    this.pressed = false;
    this.sendInput = sendInput;
  }

  init = async () => {
    await this.syncGameState();
  };

  // p5.js

  setup = () => {
    this.instance.createCanvas(window.innerWidth, window.innerHeight);
    this.instance.frameRate(FPS);
  };

  draw = () => {
    this.handleInput();
    this.updateEntities();

    this.drawGrid();
    this.drawEntities();
    this.drawMinimap();
  };

  private drawGrid = () => {
    const clientPlayer = this.entities[this.clientId] as Player;
    if (!clientPlayer) {
      return;
    }

    this.instance.background(BACKGROUND_COLOR);
    this.instance.push();
    this.instance.stroke("#ffffff33");
    this.instance.strokeWeight(2);

    this.instance.scale(this.zoom);
    this.instance.translate(
      -clientPlayer.position.x + (window.innerWidth / 2) / this.zoom,
      -clientPlayer.position.y + (window.innerHeight / 2) / this.zoom,
    );

    const worldLeft = clientPlayer.position.x - (window.innerWidth / 2) / this.zoom;
    const worldRight = clientPlayer.position.x + (window.innerWidth / 2) / this.zoom;
    const worldTop = clientPlayer.position.y - (window.innerHeight / 2) / this.zoom;
    const worldBottom = clientPlayer.position.y + (window.innerHeight / 2) / this.zoom;

    const startCol = Math.floor(worldLeft / GRID_SIZE) * GRID_SIZE;
    const endCol = Math.ceil(worldRight / GRID_SIZE) * GRID_SIZE;
    for (let x = startCol; x <= endCol; x += GRID_SIZE) {
      this.instance.line(x, worldTop, x, worldBottom);
    }

    const startRow = Math.floor(worldTop / GRID_SIZE) * GRID_SIZE;
    const endRow = Math.ceil(worldBottom / GRID_SIZE) * GRID_SIZE;
    for (let y = startRow; y <= endRow; y += GRID_SIZE) {
      this.instance.line(worldLeft, y, worldRight, y);
    }

    this.instance.pop();
  };

  private drawEntities = () => {
    const clientPlayer = this.entities[this.clientId];
    if (!clientPlayer) {
      return;
    }

    this.instance.push();
    this.instance.scale(this.zoom);
    this.instance.translate(
      -clientPlayer.position.x + (window.innerWidth / 2) / this.zoom,
      -clientPlayer.position.y + (window.innerHeight / 2) / this.zoom,
    );

    Object.values(this.entities)
      .filter(entity => entity instanceof Player)
      .forEach(player => player.drawTrail(this.instance));
    Object.values(this.entities)
      .forEach(entity => entity.draw(this.instance, DEBUG));
    this.instance.pop();
  };

  private drawMinimap = () => {
    const clientPlayer = this.entities[this.clientId] as Player;
    if (!clientPlayer) {
      return;
    }
    this.minimap.draw(this.instance, clientPlayer, this.entities);
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

  // server messaging

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
    const overwrite = data.timestamp > this.delta.timestamp;
    Object.entries(data.updated)
      .forEach(entry => {
        const [id, data] = entry;
        if (!overwrite && this.delta.updated[id]) {
          return;
        }
        this.delta.updated[id] = data;
      });
    this.delta.removed = [...this.delta.removed, ...data.removed];
    this.delta.timestamp = data.timestamp;
  };

  private handleEntityData(id: string, data: EntityData, ) {
    if (this.entities[id]) {
      this.entities[id].update(data.position);
      return;
    }

    switch (data.type) {
    case "player":
      this.entities[id] = new Player(data.position, (data as PlayerEntityData).username);
      break;

    case "projectile":
      this.entities[id] = new Projectile(data.position);
      break;

    case "powerup":
      this.entities[id] = new Powerup(data.position, (data as PowerupEntityData).ability);
      break;

    default:
      break;
    }
  }

  private syncGameState = async () => {
    await fetchGameSnapshotData()
      .then((snapshot) => {
        if (!snapshot) {
          return;
        }

        Object.entries(snapshot.entities)
          .forEach(entry => {
            const [id, data] = entry;
            this.handleEntityData(id, data);
          });
      });
  };

  // updates

  private updateEntities = () => {
    // TODO: restore some sort of removal animation
    this.delta.removed
      .forEach(id => {
        const entity = this.entities[id];
        if (entity instanceof Player) {
          (entity as Player).removed = true;
        }

        // keep client's player
        if (id === this.clientId) {
          return;
        }
        delete this.entities[id];
      });

    Object.entries(this.delta.updated)
      .forEach(entry => {
        const [id, data] = entry;
        this.handleEntityData(id, data);
      });

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
