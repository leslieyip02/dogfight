import p5 from "p5";

import type { EntityData } from "../../pb/entities";
import type { Vector } from "../../pb/vector";
import Spritesheet from "../graphics/sprites";
import { type AbilityFlag, MULTISHOT_ABILITY_FLAG, SHIELD_ABILITY_FLAG, toAbilityName, WIDE_BEAM_ABILITY_FLAG } from "../logic/abilities";
import type { Entity } from "./Entity";

const POWERUP_WIDTH = 80;

class Powerup implements Entity {
  position: Vector;
  rotation: number;
  ability: AbilityFlag;

  sprite: p5.Image | null;

  constructor(data: EntityData) {
    if (!data.position) {
      throw new Error(`expected entity data but got ${data}`);
    }
    this.position = data.position;
    this.rotation = data.rotation;

    const powerupData = data.powerupData;
    if (!powerupData) {
      throw new Error(`expected powerup data but got ${data}`);
    }
    this.ability = powerupData.ability;

    this.sprite = Spritesheet.get(this.spriteName());
  }

  update = (data: EntityData) => {
    if (!data.position) {
      return;
    }
    this.position = data.position;
    this.rotation = data.rotation;
  };

  removalAnimationName = () => {
    return null;
  };

  draw = (instance: p5, debug?: boolean) => {
    this.drawModel(instance);
    if (debug) {
      this.drawDebug(instance);
    }
  };

  drawIcon = (instance: p5) => {
    instance.push();
    instance.fill(this.iconFill());
    instance.circle(0, 0, 8);
    instance.pop();
  };

  private drawModel = (instance: p5) => {
    if (!this.sprite) {
      this.sprite = Spritesheet.get(this.spriteName());
    }
    if (!this.sprite) {
      return;
    }

    instance.push();
    instance.translate(this.position.x, this.position.y);
    instance.translate(-this.sprite.width/2, -this.sprite.height/2);
    instance.image(this.sprite, 0, 0);
    instance.pop();
  };

  private drawDebug = (instance: p5) => {
    instance.push();
    instance.translate(this.position.x, this.position.y);
    instance.noFill();
    instance.stroke("#ff0000");
    instance.rect(-POWERUP_WIDTH / 2, -POWERUP_WIDTH / 2, POWERUP_WIDTH, POWERUP_WIDTH);
    instance.pop();
  };

  private spriteName = () => {
    const spriteName = toAbilityName(this.ability);
    if (!spriteName) {
      throw new Error(`unexpected ability ${this.ability}`);
    }
    return spriteName;
  };

  private iconFill = () => {
    switch (this.ability) {
    case MULTISHOT_ABILITY_FLAG:
      return "#fac811";
    case WIDE_BEAM_ABILITY_FLAG:
      return "#f073ff";
    case SHIELD_ABILITY_FLAG:
      return "#36d9b0";
    default:
      throw new TypeError("invalid flags");
    }
  };
}

export default Powerup;
