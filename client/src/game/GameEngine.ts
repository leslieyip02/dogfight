import p5 from "p5";
import Player from "./entities/Player";
import type { GameInputEventData, GameJoinEventData, GameQuitEventData, GameUpdatePositionEventData } from "./GameEvent";
import Projectile from "./entities/Projectile";

const API_URL = import.meta.env.VITE_API_URL;

const FPS = 60;
const BACKGROUND_COLOR = "#111111";

class GameEngine {
  instance: p5;

  clientId: string;
  roomId: string;

  players: { [clientId: string]: Player };
  projectiles: { [id: string]: Projectile };

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
    case "update":
      this.handleUpdate(data.data as GameUpdatePositionEventData);
      break;
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

  private handleUpdate = async (data: GameUpdatePositionEventData) => {
    let needFetch = false;
    Object.entries(data.players)
      .forEach(entry => {
        const [id, position] = entry;
        const player = this.players[id];
        if (!player) {
          needFetch = true;
          return;
        }
        player.update(position);
      });

    if (needFetch) {
      await this.fetchPlayers();
    }

    const expiredIds = new Set(Object.keys(this.projectiles));
    Object.entries(data.projectiles)
      .forEach(entry => {
        const [id, position] = entry;
        const projectile = this.projectiles[id];
        if (!projectile) {
          this.projectiles[id] = new Projectile(position);
          return;
        }
        projectile.update(position);
        expiredIds.delete(id);
      });

    expiredIds.forEach(id => {
      this.projectiles[id].destroy();
      delete this.projectiles[id];
    });
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

  private normalize = (value: number, full: number): number => {
    const delta = value - full / 2;
    return Math.sign(delta) * Math.min(Math.abs(delta / (full / 2 * 0.8)), 1.0);
  };
};

export default GameEngine;
