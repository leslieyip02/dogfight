import p5 from "p5";
import Player from "./entities/Player";
import type { GameInputEventData, GameJoinEventData, GameQuitEventData, GameUpdateEventData } from "./GameEvent";

const FPS = 60;
const BACKGROUND_COLOR = "#111111";
const GRID_SIZE = 64;

class GameEngine {
  instance: p5;
  clientId: string;
  players: { [clientId: string]: Player };
  sendInput: (data: GameInputEventData) => void;

  constructor(instance: p5, clientId: string, sendInput: (data: GameInputEventData) => void) {
    this.instance = instance;
    this.instance.setup = this.setup;
    this.instance.draw = this.draw;
    this.clientId = clientId;
    this.players = {};
    this.sendInput = sendInput;
  }

  setup = () => {
    this.instance.createCanvas(window.innerWidth, window.innerHeight);
    this.instance.frameRate(FPS);
  };

  draw = () => {
    this.sendInput({
      clientId: this.clientId,
      mouseX: this.normalize(this.instance.mouseX, window.innerWidth),
      mouseY: this.normalize(this.instance.mouseY, window.innerHeight),
    });

    this.instance.background(BACKGROUND_COLOR);
    this.drawGrid();

    const player = this.players[this.clientId];
    player.draw(this.instance);

    // TODO: filter nearby players
    // Object.values(this.players)
    //   .forEach(player => {
    //     this.instance.circle(player.position.x, player.position.y, 80);
    //     this.instance.text(player.username, player.position.x, player.position.y);
    //   });
  };

  private drawGrid = () => {
    const player = this.players[this.clientId];
    const dx = player.position.x % GRID_SIZE;
    const dy = player.position.y % GRID_SIZE;

    const rows = Math.ceil(window.innerHeight / GRID_SIZE) + 1;
    const cols = Math.ceil(window.innerWidth / GRID_SIZE) + 1;

    this.instance.stroke("#ffffff33");
    this.instance.strokeWeight(2);
    for (let r = 0; r < rows; r++) {
      this.instance.line(0, r * GRID_SIZE + dy, window.innerWidth, r * GRID_SIZE + dy);
    }
    for (let c = 0; c < cols; c++) {
      this.instance.line(c * GRID_SIZE + dx, 0, c * GRID_SIZE + dx,  window.innerHeight);
    }
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
      this.handleUpdate(data.data as GameUpdateEventData);
      break;
    default:
      return;
    }
  };

  private handleJoin = (data: GameJoinEventData) => {
    this.players[data.clientId] = new Player(data.username, data.x, data.y, data.theta);
  };

  private handleQuit = (data: GameQuitEventData) => {
    delete this.players[data.clientId];
  };

  private handleUpdate = (data: GameUpdateEventData) => {
    Object.entries(data)
      .forEach(entry => {
        const [clientId, {x, y, theta}] = entry;
        this.players[clientId]?.update(x, y, theta);
      });
  };

  private normalize = (value: number, full: number): number => {
    return Math.min((full / 2 - value) / (full / 2 * 0.8), 1.0);
  };
};

export default GameEngine;
