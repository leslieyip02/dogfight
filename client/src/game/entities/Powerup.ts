import p5 from "p5";

import type { EntityPosition } from "../types/entity";
import type { Entity } from "./Entity";

export type PowerupAbility = "multishot";

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

  draw = (instance: p5, debug?: boolean) => {
    instance.push();

    instance.translate(this.position.x, this.position.y);
    instance.rotate(this.position.theta);

    instance.noStroke();
    instance.fill("#00ff00");
    if (debug) {
      instance.stroke("#ff0000");
    }

    instance.circle(0, 0, 10);

    instance.pop();
  };
}

export default Powerup;
