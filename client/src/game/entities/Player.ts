import p5 from "p5";
import type { Entity } from "./Entity";
import type { EntityPosition } from "../GameEvent";

const DEBUG = import.meta.env.VITE_DEBUG;

const GRID_SIZE = 64;

class Player implements Entity {
  position: EntityPosition;
  username: string;
  roll: number;

  constructor(username: string, position: EntityPosition) {
    this.username = username;
    this.position = position;
    this.roll = 0;
  }

  update = (position: EntityPosition) => {
    this.roll = Math.sign(position.theta - this.position.theta);
    this.position = position;
  };

  draw = (instance: p5) => {
    instance.push();

    instance.translate(this.position.x, this.position.y);
    
    instance.push();
    instance.rotate(this.position.theta);
    instance.fill("#ffffff");
    instance.triangle(
      40, 0,
      -40, 40,
      -40, -40,
    );
    instance.pop();

    instance.noFill();
    instance.stroke("#ffffff");
    instance.strokeWeight(1);
    instance.textAlign(instance.CENTER);
    instance.rectMode(instance.CENTER);
    instance.text(this.username, 0, -65);

    if (DEBUG) {
      instance.stroke("#ff0000");
      instance.circle(0, 0, 80);
      instance.line(0, 0, Math.cos(this.position.theta) * 120, Math.sin(this.position.theta) * 120);
      instance.text(`position: (${this.position.x.toFixed(2)}, ${this.position.y.toFixed(2)}), theta: ${this.position.theta.toFixed(2)}`, 0, -85);
    }

    instance.pop();
  };

  drawGrid = (instance: p5) => {
    instance.push();

    const dx = this.position.x % GRID_SIZE;
    const dy = this.position.y % GRID_SIZE;

    const rows = Math.ceil(window.innerHeight / GRID_SIZE) + 1;
    const cols = Math.ceil(window.innerWidth / GRID_SIZE) + 1;

    instance.stroke("#ffffff33");
    instance.strokeWeight(2);
    for (let r = 0; r < rows; r++) {
      instance.line(0, r * GRID_SIZE - dy, window.innerWidth, r * GRID_SIZE - dy);
    }
    for (let c = 0; c < cols; c++) {
      instance.line(c * GRID_SIZE - dx, 0, c * GRID_SIZE - dx,  window.innerHeight);
    }

    instance.pop();
  };
};

export default Player;
