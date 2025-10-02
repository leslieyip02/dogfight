import p5 from "p5";

import type { EntityPosition } from "../types/entity";
import type { Entity } from "./Entity";

export type PowerupAbility = "multishot";

const POWERUP_WIDTH = 20;

class Powerup implements Entity {
  position: EntityPosition;

  ability: PowerupAbility;

  constructor(position: EntityPosition, ability: PowerupAbility) {
    this.position = position;
    this.ability = ability;
  }

  update = (position?: EntityPosition) => {
    if (!position) {
      return;
    }
    this.position = position;
  };

  removalAnimationName = () => {
    return null;
  };

  draw = (instance: p5, debug?: boolean) => {
    instance.push();
    instance.translate(this.position.x, this.position.y);
    instance.noStroke();
    instance.fill("#00ff00");
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
    instance.fill("#00ff00");
    instance.circle(0, 0, 8);
    instance.pop();
  };
}

export default Powerup;
