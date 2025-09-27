import type p5 from "p5";
import type { Image } from "p5";

import type { EntityPosition } from "../types/entity";
import type { Spritesheet } from "../utils/graphics";
import type { Entity } from "./Entity";

const PLAYER_RADIUS = 40;
const MAX_PLYAER_TRAIL_POINTS = 32;

class Player implements Entity {
  static spritesheet: Spritesheet;

  position: EntityPosition;
  username: string;
  image: Image;

  roll: number;
  removed: boolean;

  previousPositions: EntityPosition[];

  constructor(position: EntityPosition, username: string) {
    this.position = position;
    this.username = username;

    // TODO: replace with something more robust
    const spriteNames = ["alpha", "beta"];
    this.image = Player.spritesheet[spriteNames[Math.floor(Math.random() * spriteNames.length)]][0];

    this.roll = 0;
    this.removed = false;

    this.previousPositions = [];
  }

  update = (position?: EntityPosition) => {
    if (!position || this.removed) {
      return;
    }

    const previousPosition: EntityPosition = {
      x: this.position.x - Math.cos(this.position.theta) * PLAYER_RADIUS,
      y: this.position.y - Math.sin(this.position.theta) * PLAYER_RADIUS,
      theta: this.position.theta,
    };
    this.previousPositions.push(previousPosition);
    if (this.previousPositions.length > MAX_PLYAER_TRAIL_POINTS) {
      this.previousPositions.shift();
    }

    // TODO: apply a transform for this?
    this.roll = Math.sign(position.theta - this.position.theta);
    this.position = position;
  };

  remove = () => {
    this.removed = true;
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
    // instance.triangle(
    //   RADIUS, 0,
    //   -RADIUS, RADIUS,
    //   -RADIUS, -RADIUS,
    // );
    instance.translate(-this.image.width / 2, -this.image.height / 2);
    instance.image(this.image, 0, 0);
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
      instance.circle(0, 0, 2 * PLAYER_RADIUS);
      instance.line(0, 0, Math.cos(this.position.theta) * 120, Math.sin(this.position.theta) * 120);
      instance.text(`position: (${this.position.x.toFixed(2)}, ${this.position.y.toFixed(2)}), theta: ${this.position.theta.toFixed(2)}`, 0, -85);
      instance.pop();
    }

    instance.pop();
  };

  drawIcon = (instance: p5) => {
    instance.fill("#ff0000");
    instance.circle(0, 0, 8);
  };

  drawTrail = (instance: p5) => {
    if (this.removed) {
      return;
    }

    instance.push();
    instance.stroke("#ffa320");
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
};

export default Player;
