import p5 from "p5";
import type { Entity } from "./Entity";
import type { EntityPosition } from "../GameEvent";

const RADIUS = 40;
const MAX_TRAIL_POINTS = 32;

class Player implements Entity {
  position: EntityPosition;
  username: string;
  roll: number;

  previousPositions: EntityPosition[];

  removed: boolean;
  onRemove: () => void;

  constructor(username: string, position: EntityPosition, onRemove: () => void) {
    this.username = username;
    this.position = position;
    this.roll = 0;

    this.previousPositions = [];

    this.removed = false;
    this.onRemove = onRemove;
  }

  update = (position?: EntityPosition) => {
    if (!position || this.removed) {
      return;
    }

    this.previousPositions.push({ ...this.position });
    if (this.previousPositions.length > MAX_TRAIL_POINTS) {
      this.previousPositions.shift();
    }

    this.roll = Math.sign(position.theta - this.position.theta);
    this.position = position;
  };

  draw = (instance: p5, debug?: boolean) => {
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

    if (debug) {
      instance.push();
      instance.stroke("#ff0000");
      instance.circle(0, 0, 2 * RADIUS);
      instance.line(0, 0, Math.cos(this.position.theta) * 120, Math.sin(this.position.theta) * 120);
      instance.text(`position: (${this.position.x.toFixed(2)}, ${this.position.y.toFixed(2)}), theta: ${this.position.theta.toFixed(2)}`, 0, -85);
      instance.pop();
    }

    instance.pop();
  };

  drawTrail = (instance: p5) => {
    if (this.removed) {
      return;
    }

    instance.push();
    instance.stroke("#888888");
    instance.strokeWeight(4);
    instance.noFill();
    for (let i = 0; i < this.previousPositions.length - 1; i++) {
      instance.line(
        this.previousPositions[i].x, this.previousPositions[i].y,
        this.previousPositions[i + 1].x, this.previousPositions[i + 1].y,
      );
    }
    instance.pop();
  };

  remove = () => {
    this.removed = true;
    this.onRemove();
  };
};

export default Player;
