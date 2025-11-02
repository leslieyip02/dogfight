import type p5 from "p5";

import type { EntityData } from "../../pb/entities";
import type { Vector } from "../../pb/vector";
import { SOUNDS } from "../audio/sounds";
import Spritesheet from "../graphics/sprites";
import { type AbilityFlag,isAbilityActive, SHIELD_ABILITY_FLAG } from "../logic/abilities";
import type { BaseEntity } from "./Entity";

export const PLAYER_MAX_SPEED = 20.0;

const PLAYER_WIDTH = 96;
const PLAYER_MAX_TRAIL_POINTS = 24;

class Player implements BaseEntity {
  position: Vector;
  velocity: Vector;
  rotation: number;
  flags: AbilityFlag;

  username: string;
  score: number;
  sprite: p5.Image;
  previousPositions: Vector[];

  constructor(data: EntityData) {
    if (!data.position || !data.velocity) {
      throw new Error(`expected entity data but got ${data}`);
    }
    this.position = data.position;
    this.velocity = data.velocity;
    this.rotation = data.rotation;

    const playerData = data.playerData;
    if (!playerData) {
      throw new Error(`expected player data but got ${data}`);
    }
    this.username = playerData.username;
    this.score = playerData.score;
    this.flags = playerData.flags;

    const sprite = Spritesheet.getPlayerSprite(this.username);
    if (!sprite) {
      throw new Error("unable to load player sprite");
    }
    this.sprite = sprite;

    this.previousPositions = [];
  }

  update = (data: EntityData) => {
    if (!data.position || !data.velocity || !data.rotation) {
      return;
    }
    this.position = data.position;
    this.velocity = data.velocity;
    this.rotation = data.rotation;

    const playerData = data.playerData;
    if (playerData) {
      if (this.score < playerData.score) {
        this.score = playerData.score;
        SOUNDS["score"].play();
      }

      if (this.flags != playerData.flags) {
        this.flags = playerData.flags;
        SOUNDS["pickup"].play();
      }
    }
  };

  removalAnimationName = () => {
    return "explosionBig";
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
