import p5 from "p5";

import type { Entity as EntityData } from "../../pb/entities";
import type { Vector } from "../../pb/vector";
import { type AbilityFlag, MULTISHOT_ABILITY_FLAG, SHIELD_ABILITY_FLAG, WIDE_BEAM_ABILITY_FLAG } from "../utils/abilities";
import type { BaseEntity } from "./Entity";

const POWERUP_WIDTH = 20;

class Powerup implements BaseEntity {
  position: Vector;
  rotation: number;
  ability: AbilityFlag;

  constructor(data: EntityData) {
    if (!data.position || !data.rotation) {
      throw new Error(`expected entity data but got ${data}`);
    }
    this.position = data.position;
    this.rotation = data.rotation;

    const powerupData = data.powerupData;
    if (!powerupData) {
      throw new Error(`expected powerup data but got ${data}`);
    }
    this.ability = powerupData.ability;
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
