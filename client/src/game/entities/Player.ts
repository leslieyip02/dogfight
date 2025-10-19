import type p5 from "p5";

import type { EntityData, PlayerEntityData } from "../types/entity";
import type { Vector } from "../types/geometry";
import { type AbilityFlag,isAbilityActive, SHIELD_ABILITY_FLAG } from "../utils/abilities";
import type { Spritesheet } from "../utils/sprites";
import type { Entity } from "./Entity";

export const PLAYER_MAX_SPEED = 20.0;

const PLAYER_WIDTH = 96;
const PLAYER_MAX_TRAIL_POINTS = 24;

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
  flags: AbilityFlag;

  username: string;
  score: number;
  sprite: p5.Image;
  previousPositions: Vector[];

  constructor(data: PlayerEntityData) {
    this.position = data.position;
    this.velocity = data.velocity;
    this.rotation = data.rotation;
    this.flags = data.flags;

    this.username = data.username;
    this.score = data.score;
    this.sprite = chooseSprite(data.username);
    this.previousPositions = [];
  }

  update = (data: EntityData) => {
    if (!data.position || !data.velocity || !data.rotation) {
      return;
    }
    this.position = data.position;
    this.velocity = data.velocity;
    this.rotation = data.rotation;

    const playerEntityData = data as PlayerEntityData;
    if (playerEntityData) {
      this.score = playerEntityData.score;
      this.flags = playerEntityData.flags;
    }
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

    if (isAbilityActive(this.flags, SHIELD_ABILITY_FLAG)) {
      instance.push();
      instance.stroke("#ffff0088");
      instance.strokeWeight(4);
      instance.noFill();
      instance.circle(this.sprite.width / 2, this.sprite.height / 2, this.sprite.width * 1.2);
      instance.pop();
    }

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
    instance.rectMode(instance.CENTER);
    instance.textAlign(instance.CENTER);
    instance.textFont("Courier New");
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

  drawTrail = (instance: p5, debug?: boolean) => {
    const previousPosition: Vector = {
      x: this.position.x - Math.cos(this.rotation) * PLAYER_WIDTH / 2,
      y: this.position.y - Math.sin(this.rotation) * PLAYER_WIDTH / 2,
    };
    this.previousPositions.push(previousPosition);
    if (this.previousPositions.length > PLAYER_MAX_TRAIL_POINTS) {
      this.previousPositions.shift();
    }

    instance.push();
    instance.strokeWeight(4);
    const color = instance.color("#ffa320");
    for (let i = 0; i < this.previousPositions.length - 1; i++) {
      color.setAlpha(Math.min(i/(PLAYER_MAX_TRAIL_POINTS / 4), 1) * 255);
      instance.stroke(color);
      instance.line(
        this.previousPositions[i].x, this.previousPositions[i].y,
        this.previousPositions[i + 1].x, this.previousPositions[i + 1].y,
      );
    }

    if (debug && this.previousPositions.length > 0) {
      instance.push();
      instance.stroke("#ff0000");
      instance.strokeWeight(1);
      instance.circle(this.previousPositions[0].x, this.previousPositions[0].y, 10);
      instance.pop();
    }

    instance.pop();
  };
};

export default Player;
