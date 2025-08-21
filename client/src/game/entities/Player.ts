import p5 from "p5";
import type { Entity } from "./Entity";

class Player implements Entity {
  position: p5.Vector;
  username: string;

  constructor(username: string, x: number, y: number) {
    this.username = username;
    this.position = new p5.Vector(x, y);
  }

  update = (x: number, y: number) => {
    this.position.x = x;
    this.position.y = y;
  };
};

export default Player;
