import type p5 from "p5";

import type { EntityData, PlayerEntityData } from "../types/entity";
import type { Vector } from "../types/geometry";
import type { Spritesheet } from "../utils/graphics";
import type { Entity } from "./Entity";

const PLAYER_WIDTH = 96;
const MAX_PLYAER_TRAIL_POINTS = 32;

const PLAYER_SPRITE_NAMES = ["alpha", "bravo", "charlie", "delta"];

function chooseSprite(username: string): p5.Image {
  // TODO: consider replacing this with a normal field?
  const asciiSum = [...username].reduce((a, b) => a + b.charCodeAt(0), 0);
  return Player.spritesheet[PLAYER_SPRITE_NAMES[asciiSum % PLAYER_SPRITE_NAMES.length]][0];
}

class Player implements Entity {
  static spritesheet: Spritesheet;

  position: Vector;
  velocity: Vector;
  rotation: number;

  username: string;
  sprite: p5.Image;
  previousPositions: Vector[];

  constructor(data: PlayerEntityData) {
    this.position = data.position;
    this.velocity = data.velocity;
    this.rotation = data.rotation;

    this.username = data.username;
    this.sprite = chooseSprite(data.username);
    this.previousPositions = [];
  }

  update = (data: EntityData) => {
    if (!data.position || !data.rotation) {
      return;
    }

    const previousPosition: Vector = {
      x: this.position.x - Math.cos(this.rotation) * PLAYER_WIDTH / 2,
      y: this.position.y - Math.sin(this.rotation) * PLAYER_WIDTH / 2,
    };
    this.previousPositions.push(previousPosition);
    if (this.previousPositions.length > MAX_PLYAER_TRAIL_POINTS) {
      this.previousPositions.shift();
    }
    this.position = data.position;
    this.rotation = data.rotation;
  };

  removalAnimationName = () => {
    return "smallExplosion";
  };

  draw = (instance: p5, debug?: boolean) => {
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
    instance.rotate(this.rotation);
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
      instance.line(0, 0, Math.cos(this.rotation) * 120, Math.sin(this.rotation) * 120);
      instance.text(`position: (${this.position.x.toFixed(2)}, ${this.position.y.toFixed(2)}), rotation: ${this.rotation.toFixed(2)}`, 0, -85);
      instance.pop();
    }

    instance.pop();
  };

  drawTrail = (instance: p5) => {
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
