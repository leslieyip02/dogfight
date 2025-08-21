import type p5 from "p5";
import Player from "./entities/Player";
import type { GameJoinEventData, GameQuitEventData, GameUpdateEventData } from "./GameEvent";

class GameEngine {
  instance: p5;
  players: { [clientId: string]: Player };

  constructor(instance: p5) {
    this.instance = instance;
    this.instance.setup = this.setup;
    this.instance.draw = this.draw;
    this.players = {};
  }

  setup = () => {
    this.instance.createCanvas(window.innerWidth, window.innerHeight);
  };

  draw = () => {
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
