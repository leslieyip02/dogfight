import p5 from "p5";

import type { EntityData, PowerupEntityData } from "../types/entity";
import type { Vector } from "../types/geometry";
import { type AbilityFlag, MULTISHOT_ABILITY_FLAG, SHIELD_ABILITY_FLAG, WIDE_BEAM_ABILITY_FLAG } from "../utils/abilities";
import type { Entity } from "./Entity";

const POWERUP_WIDTH = 20;

class Powerup implements Entity {
  position: Vector;
  rotation: number;
  ability: AbilityFlag;

  constructor(data: PowerupEntityData) {
    this.position = data.position;
    this.rotation = data.rotation;
    this.ability = data.ability;
  }

  update = (data: EntityData) => {
    if (!data.position || !data.rotation) {
      return;
    }
    this.position = data.position;
    this.rotation = data.rotation;
  };

  removalAnimationName = () => {
    return null;
  };

  draw = (instance: p5, debug?: boolean) => {
    instance.push();
    instance.translate(this.position.x, this.position.y);
    instance.noStroke();
    instance.fill(this.fill());
    instance.circle(0, 0, 10);

    if (debug) {
      instance.push();
      instance.noFill();
      instance.stroke("#ff0000");
      instance.rect(-POWERUP_WIDTH / 2, -POWERUP_WIDTH / 2, POWERUP_WIDTH, POWERUP_WIDTH);
      instance.pop();
    }

    instance.pop();
  };

  drawIcon = (instance: p5) => {
    instance.push();
    instance.fill(this.fill());
    instance.circle(0, 0, 8);
    instance.pop();
  };

  fill = () => {
    switch (this.ability) {
    case MULTISHOT_ABILITY_FLAG:
      return "#00ffff";
    case WIDE_BEAM_ABILITY_FLAG:
      return "#ff00ff";
    case SHIELD_ABILITY_FLAG:
      return "#ffff00";
    default:
      throw new TypeError("invalid flags");
    }
  };
}

export default Powerup;
