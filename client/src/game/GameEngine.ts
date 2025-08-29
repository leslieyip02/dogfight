import p5 from "p5";
import Player from "./entities/Player";
import type { GameInputEventData, GameJoinEventData, GameQuitEventData, GameUpdatePositionEventData, GameUpdatePowerupEventData } from "./GameEvent";
import Projectile from "./entities/Projectile";
import Explosion from "./entities/Explosion";
import Powerup from "./entities/Powerup";

const API_URL = import.meta.env.VITE_API_URL;

const FPS = 60;
const BACKGROUND_COLOR = "#111111";

class GameEngine {
  instance: p5;

  clientId: string;
  roomId: string;

  players: { [id: string]: Player };
  projectiles: { [id: string]: Projectile };
  powerups: { [id: string]: Powerup };
  explosions: { [id: string]: Explosion };

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

    this.pressed = false;
    this.sendInput = sendInput;
  }

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
      await this.fetchPlayers();
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
  };

  mousePressed = () => {
    this.pressed = true;
  };

  receive = (event: MessageEvent) => {
    const data = JSON.parse(event.data);
    switch (data["type"]) {
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
    this.players[data.id] = new Player(data.username, data.position);
  };

  private handleQuit = (data: GameQuitEventData) => {
    delete this.players[data.id];
  };

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

    destroyedIds.forEach(id => {
      this.players[id]?.destroy();

      if (!this.explosions[id]) {
        this.explosions[id] = new Explosion({ ...this.players[id].position });
      }

      if (id === this.clientId) {
        return;
      }
      delete this.players[id];
    });

    if (needFetch) {
      await this.fetchPlayers();
    }
  };

  private fetchPlayers = async () => {
    // fallback
    await fetch(`${API_URL}/room/players?roomId=${this.roomId}`)
      .then(response => response.json())
      .then(data => {
        data.map((player: GameJoinEventData) => {
          this.players[player.id] = new Player(player.username, player.position);
        });
      });
  };

  private updateProjectiles = (data: GameUpdatePositionEventData) => {
    const destroyedIds = new Set(Object.keys(this.projectiles));
    Object.entries(data.projectiles)
      .forEach(entry => {
        const [id, position] = entry;
        const projectile = this.projectiles[id];
        if (!projectile) {
          this.projectiles[id] = new Projectile(position);
          return;
        }
        projectile.update(position);
        destroyedIds.delete(id);
      });

    destroyedIds.forEach(id => delete this.projectiles[id]);
  };

  private updatePowerups = (data: GameUpdatePowerupEventData) => {
    if (!data.position) {
      delete this.powerups[data.id];
      return;
    }

    this.powerups[data.id] = new Powerup(data.type, data.position);
  };

  private normalize = (value: number, full: number): number => {
    const delta = value - full / 2;
    return Math.sign(delta) * Math.min(Math.abs(delta / (full / 2 * 0.8)), 1.0);
  };
};

export default GameEngine;
