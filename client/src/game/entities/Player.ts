import type p5 from "p5";

import type { EntityData } from "../../pb/entities";
import type { Vector } from "../../pb/vector";
import Audiosheet from "../audio/audio";
import { generateExplosionAnimation } from "../graphics/animation";
import Spritesheet from "../graphics/sprites";
import { type AbilityFlag, isAbilityActive, SHIELD_ABILITY_FLAG } from "../logic/abilities";
import type { Entity } from "./Entity";

export const PLAYER_MAX_SPEED = 20.0;
const PLAYER_WIDTH = 80;

class Player implements Entity {
  position: Vector;
  velocity: Vector;
  rotation: number;
  flags: AbilityFlag;

  username: string;
  score: number;
  previousPositions: Vector[];

  sprite: p5.Image | null;

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

    this.previousPositions = [];

    this.sprite = Spritesheet.getPlayerSprite(this.username);
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
        Audiosheet.get("pickup")?.play();
      }

      if (this.flags != playerData.flags) {
        this.flags = playerData.flags;
        Audiosheet.get("pickup")?.play();
      }
    }
  };

  remove = () => {
    return generateExplosionAnimation("explosionBig", this.position);
  };

  draw = (instance: p5, debug?: boolean) => {
    this.drawModel(instance);
    this.drawUsername(instance);
    if (debug) {
      this.drawDebug(instance);
    }
  };

  drawIcon = (instance: p5) => {
    instance.fill("#ff0000");
    instance.circle(0, 0, 8);
  };

  private drawModel = (instance: p5) => {
    if (!this.sprite) {
      this.sprite = Spritesheet.getPlayerSprite(this.username);
    }
    if (!this.sprite) {
      return;
    }

    instance.push();
    instance.translate(this.position.x, this.position.y);
    instance.rotate(this.rotation);
    instance.translate(-this.sprite.width / 2, -this.sprite.height / 2);
    instance.image(this.sprite, 0, 0);

    if (isAbilityActive(this.flags, SHIELD_ABILITY_FLAG)) {
      instance.push();
      instance.stroke("#36d9b0");
      instance.fill("#36d9b033");
      instance.strokeWeight(4);
      instance.circle(this.sprite.width / 2, this.sprite.height / 2, this.sprite.width * 1.2);
      instance.pop();
    }

    instance.pop();
  };

  private drawUsername = (instance: p5) => {
    instance.push();
    instance.translate(this.position.x, this.position.y);
    instance.noFill();
    instance.stroke("#ffffff");
    instance.strokeWeight(1);
    instance.textAlign(instance.CENTER);
    instance.textFont("Courier New");
    instance.text(this.username, 0, -80);
    instance.pop();
  };

  private drawDebug = (instance: p5) => {
    instance.push();
    instance.translate(this.position.x, this.position.y);
    instance.noFill();
    instance.stroke("#ff0000");
    instance.strokeWeight(1);

    instance.push();
    instance.rotate(this.rotation);
    instance.line(0, 0, 120, 0);
    instance.translate(-PLAYER_WIDTH / 2, -PLAYER_WIDTH / 2);
    instance.rect(0, 0, PLAYER_WIDTH);
    instance.pop();

    instance.push();
    instance.fill("#ff0000");
    instance.textAlign(instance.CENTER);
    instance.textSize(16);
    instance.text(`position: (${this.position.x.toFixed(2)}, ${this.position.y.toFixed(2)}), rotation: ${this.rotation.toFixed(2)}`, 0, -100);
    instance.pop();

    instance.pop();
  };
};

export default Player;
