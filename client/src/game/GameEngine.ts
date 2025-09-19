import p5 from "p5";

import Player from "./entities/Player";
import type {
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

const FPS = 60;
export const BACKGROUND_COLOR = "#111111";

class GameEngine {
  instance: p5;

  clientId: string;
  roomId: string;

  players: { [id: string]: Player };
  projectiles: { [id: string]: Projectile };
  powerups: { [id: string]: Powerup };
  explosions: { [id: string]: Explosion };

  minimap: Minimap;

  pressed: boolean;
  sendInput: (data: GameInputEventData) => void;

  constructor(instance: p5, clientId: string, roomId: string, sendInput: (data: GameInputEventData) => void) {
    this.instance = instance;
    this.instance.setup = this.setup;
    this.instance.draw = this.draw;
    this.instance.mousePressed = this.mousePressed;

    this.clientId = clientId;
    this.roomId = roomId;

    this.players = {};
    this.projectiles = {};
    this.powerups = {};
    this.explosions = {};

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
    this.sendInput({
      clientId: this.clientId,
      mouseX: this.normalize(this.instance.mouseX, window.innerWidth),
      mouseY: this.normalize(this.instance.mouseY, window.innerHeight),
      mousePressed: this.pressed,
    });
    this.pressed = false;

    this.instance.background(BACKGROUND_COLOR);

    const clientPlayer = this.players[this.clientId];
    if (!clientPlayer) {
      await this.fetchState();
      return;
    }

    clientPlayer.drawGrid(this.instance);

    this.instance.push();
    this.instance.translate(
      -clientPlayer.position.x + window.innerWidth / 2,
      -clientPlayer.position.y + window.innerHeight / 2
    );

    Object.values(this.players)
      .forEach(player => player.draw(this.instance));
    Object.values(this.projectiles)
      .forEach(projectile => projectile.draw(this.instance));

    Object.values(this.powerups).forEach(powerup => {
      powerup.update();
      powerup.draw(this.instance);
    });
    Object.values(this.explosions).forEach(explosion => {
      explosion.update();
      explosion.draw(this.instance);
    });
    this.instance.pop();

    this.minimap.draw(this.instance, clientPlayer, this.players);
  };

  mousePressed = () => {
    this.pressed = true;
  };

  // server messaging

  receive = (event: MessageEvent) => {
    const data = JSON.parse(event.data);
    switch (data["type"] as GameEventType) {
    case "join":
      this.handleJoin(data.data as GameJoinEventData);
      break;
    case "quit":     
      this.handleQuit(data.data as GameQuitEventData);
      break;
    case "position": {
      const positionData = data.data as GameUpdatePositionEventData;
      this.updatePlayers(positionData);
      this.updateProjectiles(positionData);
      break;
    }
    case "powerup": {
      const powerupData = data.data as GameUpdatePowerupEventData;
      this.updatePowerups(powerupData);
      break;
    }
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

  private updatePlayers = async (data: GameUpdatePositionEventData) => {
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

    destroyedIds.forEach(id => this.players[id].remove());

    if (needFetch) {
      await this.fetchState();
    }
  };

  private updateProjectiles = (data: GameUpdatePositionEventData) => {
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
