import p5 from "p5";
import type { Entity } from "./Entity";
import type { EntityPosition } from "../GameEvent";

const DEBUG = import.meta.env.VITE_DEBUG;

const RADIUS = 40;

class Player implements Entity {
  position: EntityPosition;
  username: string;
  roll: number;
  removed: boolean;
  onRemove: () => void;

  constructor(username: string, position: EntityPosition, onRemove: () => void) {
    this.username = username;
    this.position = position;
    this.roll = 0;
    this.removed = false;
    this.onRemove = onRemove;
  }

  update = (position?: EntityPosition) => {
    if (!position || this.removed) {
      return;
    }
    this.roll = Math.sign(position.theta - this.position.theta);
    this.position = position;
  };

  draw = (instance: p5) => {
    if (this.removed) {
      return;
    }

    instance.push();

    instance.translate(this.position.x, this.position.y);

    instance.push();
    instance.rotate(this.position.theta);
    instance.fill("#ffffff");
    // TODO: consider changing to a sprite
    instance.triangle(
      RADIUS, 0,
      -RADIUS, RADIUS,
      -RADIUS, -RADIUS,
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
      instance.circle(0, 0, 2 * RADIUS);
      instance.line(0, 0, Math.cos(this.position.theta) * 120, Math.sin(this.position.theta) * 120);
      instance.text(`position: (${this.position.x.toFixed(2)}, ${this.position.y.toFixed(2)}), theta: ${this.position.theta.toFixed(2)}`, 0, -85);
    }

    instance.pop();
  };

  remove = () => {
    this.removed = true;
    this.onRemove();
  };
};

export default Player;
