import type p5 from "p5";
import Player from "./entities/Player";
import type { GameInputEventData, GameJoinEventData, GameQuitEventData, GameUpdateEventData } from "./GameEvent";

const FPS = 60;

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
      mouseX: this.instance.mouseX,
      mouseY: this.instance.mouseY,
    });

    Object.values(this.players)
      .forEach(player => {
        this.instance.circle(player.position.x, player.position.y, 80);
        this.instance.text(player.username, player.position.x, player.position.y);
      });
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
    this.players[data.clientId] = new Player(data.username, data.x, data.y);
  };

  private handleQuit = (data: GameQuitEventData) => {
    delete this.players[data.clientId];
  };

  private handleUpdate = (data: GameUpdateEventData) => {
    Object.entries(data)
      .forEach(entry => {
        const [clientId, {x, y}] = entry;
        this.players[clientId]?.update(x, y);
      });
  };
};

export default GameEngine;
