import p5 from "p5";

import Player from "./entities/Player";
import type {
  GameEvent,
  GameEventType,
  GameInputEventData,
  GameJoinEventData,
  GameQuitEventData,
  GameUpdatePositionEventData,
  GameUpdatePowerupEventData
} from "./GameEvent";
import Explosion from "./entities/Explosion";
import Minimap from "./Minimap";
import Powerup from "./entities/Powerup";
import Projectile from "./entities/Projectile";

type FetchedGameState = {
  players: GameJoinEventData[],
  powerups: GameUpdatePowerupEventData[],
};

const API_URL = import.meta.env.VITE_API_URL;
const DEBUG = import.meta.env.VITE_DEBUG;

const FPS = 60;
const GRID_SIZE = 64;

export const BACKGROUND_COLOR = "#111111";

class GameEngine {
  instance: p5;

  zoom: number;

  clientId: string;
  roomId: string;
  
  players: { [id: string]: Player };
  projectiles: { [id: string]: Projectile };
  powerups: { [id: string]: Powerup };
  explosions: { [id: string]: Explosion };

  updateEventBuffer: GameEvent[];

  minimap: Minimap;

  pressed: boolean;
  sendInput: (data: GameInputEventData) => void;

  constructor(instance: p5, clientId: string, roomId: string, sendInput: (data: GameInputEventData) => void) {
    this.instance = instance;
    this.instance.setup = this.setup;
    this.instance.draw = this.draw;
    this.instance.mousePressed = this.mousePressed;

    this.zoom = 1.0;

    this.clientId = clientId;
    this.roomId = roomId;

    this.players = {};
    this.projectiles = {};
    this.powerups = {};
    this.explosions = {};

    this.updateEventBuffer = [];

    this.minimap = new Minimap();

    this.pressed = false;
    this.sendInput = sendInput;
  }

  init = async () => {
    await this.fetchState();
  };

  // p5.js
  
  setup = () => {
    this.instance.createCanvas(window.innerWidth, window.innerHeight);
    this.instance.frameRate(FPS);
  };

  draw = async () => {
    const mouseX = this.normalize(this.instance.mouseX, window.innerWidth);
    const mouseY = this.normalize(this.instance.mouseY, window.innerHeight);
    this.sendInput({
      clientId: this.clientId,
      mouseX,
      mouseY,
      mousePressed: this.pressed,
    });
    this.pressed = false;

    this.instance.background(BACKGROUND_COLOR);

    const clientPlayer = this.players[this.clientId];
    if (!clientPlayer) {
      await this.fetchState();
      return;
    }

    const throttle = Math.min(Math.sqrt(mouseX * mouseX + mouseY * mouseY), 1.0);
    if (throttle > 0.8) {
      this.zoom = Math.max(this.zoom - 0.005, 0.8);
    } else {
      this.zoom = Math.min(this.zoom + 0.005, 1.0);
    }
    
    this.drawGrid();
    this.drawEntities();
    this.minimap.draw(this.instance, clientPlayer, this.players, this.powerups);
  };

  private drawGrid = () => {
    this.instance.push();

    const clientPlayer = this.players[this.clientId];
    const dx = clientPlayer.position.x % GRID_SIZE;
    const dy = clientPlayer.position.y % GRID_SIZE;

    const rows = Math.ceil(window.innerHeight / GRID_SIZE) + 1;
    const cols = Math.ceil(window.innerWidth / GRID_SIZE) + 1;

    this.instance.stroke("#ffffff33");
    this.instance.strokeWeight(2);
    for (let r = 0; r < rows; r++) {
      this.instance.line(0, r * GRID_SIZE - dy, window.innerWidth, r * GRID_SIZE - dy);
    }
    for (let c = 0; c < cols; c++) {
      this.instance.line(c * GRID_SIZE - dx, 0, c * GRID_SIZE - dx,  window.innerHeight);
    }
    this.instance.pop();
  };

  private drawEntities = () => {
    const clientPlayer = this.players[this.clientId];
    if (!clientPlayer) {
      return;
    }

    this.updateEntities();

    this.instance.push();
    this.instance.scale(this.zoom);
    this.instance.translate(
      -clientPlayer.position.x + window.innerWidth / 2 / this.zoom,
      -clientPlayer.position.y + window.innerHeight / 2 / this.zoom,
    );

    Object.values(this.players)
      .forEach(player => player.drawTrail(this.instance));
    Object.values(this.players)
      .forEach(player => player.draw(this.instance, DEBUG));
    Object.values(this.projectiles)
      .forEach(projectile => projectile.draw(this.instance, DEBUG));
    Object.values(this.powerups)
      .forEach(powerup => powerup.draw(this.instance, DEBUG));
    Object.values(this.explosions)
      .forEach(explosion => explosion.draw(this.instance));
    this.instance.pop();
  };

  mousePressed = () => {
    this.pressed = true;
  };

  // server messaging

  receive = (event: GameEvent) => {
    switch (event["type"] as GameEventType) {
    case "join":
      this.handleJoin(event.data as GameJoinEventData);
      break;
    case "quit":     
      this.handleQuit(event.data as GameQuitEventData);
      break;
    case "position":
    case "powerup":
      this.updateEventBuffer.push(event);
      break;
    default:
      return;
    }
  };

  private handleJoin = (data: GameJoinEventData) => {
    this.players[data.id] = new Player(data.username, data.position, this.onRemovePlayer(data.id));
  };

  private handleQuit = (data: GameQuitEventData) => {
    delete this.players[data.id];
  };

  private fetchState = async () => {
    await fetch(`${API_URL}/room/state?roomId=${this.roomId}`)
      .then(response => response.json())
      .then((data: FetchedGameState) => {
        data.players.map(player => {
          this.players[player.id] = new Player(
            player.username,
            player.position,
            this.onRemovePlayer(player.id)
          );
        });
        data.powerups.map(powerup => this.updatePowerups(powerup));
      });
  };

  // updates

  private updateEntities = () => {
    let updatePositionEventData: GameUpdatePositionEventData | null = null;
    this.updateEventBuffer.forEach(event => {
      switch (event.type) {
      case "position": {
        const data = event.data as GameUpdatePositionEventData;
        if (!updatePositionEventData || data.timestamp > updatePositionEventData.timestamp) {
          updatePositionEventData = data;
        }
        break;
      };
      case "powerup": {
        const data = event.data as GameUpdatePowerupEventData;
        this.updatePowerups(data);
        break;
      };
      default:
        break;
      }
    });

    this.updatePlayers(updatePositionEventData);
    this.updateProjectiles(updatePositionEventData);
    this.updateEventBuffer = [];

    Object.values(this.powerups)
      .forEach(powerup => powerup.update());
    Object.values(this.explosions)
      .forEach(explosion => explosion.update());
  };

  private updatePlayers = async (data: GameUpdatePositionEventData | null) => {
    if (!data) {
      return;
    }

    let needFetch = false;
    const destroyedIds = new Set(Object.keys(this.players));
    Object.entries(data.players)
      .forEach(entry => {
        const [id, position] = entry;
        const player = this.players[id];
        if (!player) {
          needFetch = true;
          return;
        }
        player.update(position);
        destroyedIds.delete(id);
      });
      
    if (needFetch) {
      await this.fetchState();
    }

    destroyedIds.forEach(id => this.players[id].remove());
  };

  private updateProjectiles = (data: GameUpdatePositionEventData | null) => {
    if (!data) {
      return;
    }

    const destroyedIds = new Set(Object.keys(this.projectiles));
    Object.entries(data.projectiles)
      .forEach(entry => {
        const [id, position] = entry;
        const projectile = this.projectiles[id];
        if (!projectile) {
          this.projectiles[id] = new Projectile(position, () => {
            delete this.projectiles[id];
          });
          return;
        }
        projectile.update(position);
        destroyedIds.delete(id);
      });

    destroyedIds.forEach(id => this.projectiles[id].remove());
  };

  private updatePowerups = (data: GameUpdatePowerupEventData) => {
    if (data.position) {
      if (this.powerups[data.id]) {
        return;
      }

      this.powerups[data.id] = new Powerup(data.type, data.position, () => {
        delete this.powerups[data.id];
      });
    } else {
      this.powerups[data.id]?.onRemove();
    }
  };

  private onRemovePlayer = (id: string) => {
    return () => {
      if (!this.explosions[id]) {
        this.explosions[id] = new Explosion({ ...this.players[id].position }, () => {
          delete this.explosions[id];
        });
      }

      // do not remove the client's player
      // to continue showing the background
      if (id === this.clientId) {
        return;
      }
      delete this.players[id];
    };
  };

  // helpers

  private normalize = (value: number, full: number): number => {
    const delta = value - full / 2;
    return Math.sign(delta) * Math.min(Math.abs(delta / (full / 2 * 0.8)), 1.0);
  };
};

export default GameEngine;
