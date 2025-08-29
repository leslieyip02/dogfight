import p5 from "p5";
import type { Entity } from "./Entity";
import type { EntityPosition } from "../GameEvent";

export type PowerupType = "multishot";

const DEBUG = import.meta.env.VITE_DEBUG;

class Powerup implements Entity {
  type: PowerupType;
  position: EntityPosition;

  constructor(type: PowerupType, position: EntityPosition) {
    this.type = type;
    this.position = position;
  }

  update = (position?: EntityPosition) => {
    if (!position) {
      return;
    }
    this.position = position;
  };

  draw = (instance: p5) => {
    instance.push();

    instance.translate(this.position.x, this.position.y);

    instance.push();

    instance.rotate(this.position.theta);

    instance.noStroke();
    instance.fill("#00ff00");
    if (DEBUG) {
      instance.stroke("#ff0000");
    }
    instance.circle(0, 0, 10);

    instance.pop();
    instance.pop();
  };
}

export default Powerup;
