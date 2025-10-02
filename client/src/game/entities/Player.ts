import type p5 from "p5";

import type { EntityPosition } from "../types/entity";
import type { Spritesheet } from "../utils/graphics";
import type { Entity } from "./Entity";

const PLAYER_WIDTH = 96;
const MAX_PLYAER_TRAIL_POINTS = 32;

const PLAYER_SPRITE_NAMES = ["alpha", "bravo", "charlie"];

function chooseSprite(username: string): p5.Image {
  // TODO: consider replacing this with a normal field?
  const asciiSum = [...username].reduce((a, b) => a + b.charCodeAt(0), 0);
  return Player.spritesheet[PLAYER_SPRITE_NAMES[asciiSum % PLAYER_SPRITE_NAMES.length]][0];
}

class Player implements Entity {
  static spritesheet: Spritesheet;

  position: EntityPosition;
  username: string;
  sprite: p5.Image;

  roll: number;
  removed: boolean;

  previousPositions: EntityPosition[];

  constructor(position: EntityPosition, username: string) {
    this.position = position;
    this.username = username;
    this.sprite = chooseSprite(username);

    this.roll = 0;
    this.removed = false;

    this.previousPositions = [];
  }

  update = (position?: EntityPosition) => {
    if (!position || this.removed) {
      return;
    }

    const previousPosition: EntityPosition = {
      x: this.position.x - Math.cos(this.position.theta) * PLAYER_WIDTH / 2,
      y: this.position.y - Math.sin(this.position.theta) * PLAYER_WIDTH / 2,
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

    this.drawModel(instance, debug);
    this.drawUsername(instance, debug);
  };

  drawIcon = (instance: p5) => {
    instance.fill("#ff0000");
    instance.circle(0, 0, 8);
  };

  drawModel = (instance: p5, debug?: boolean) => {
    instance.push();
    instance.translate(this.position.x, this.position.y);
    instance.rotate(this.position.theta);
    instance.translate(-this.sprite.width / 2, -this.sprite.height / 2);
    instance.image(this.sprite, 0, 0);

    if (debug) {
      instance.push();
      instance.stroke("#ff0000");
      instance.noFill();
      instance.rect(0, 0, this.sprite.width, this.sprite.height);
      instance.pop();
    }

    instance.pop();
  };

  drawUsername = (instance: p5, debug?: boolean) => {
    instance.push();
    instance.translate(this.position.x, this.position.y);
    instance.noFill();
    instance.stroke("#ffffff");
    instance.strokeWeight(1);
    instance.textAlign(instance.CENTER);
    instance.rectMode(instance.CENTER);
    instance.text(this.username, 0, -65);

    if (debug) {
      instance.push();
      instance.stroke("#ff0000");
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
