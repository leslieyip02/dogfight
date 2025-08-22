import p5 from "p5";
import type { Entity } from "./Entity";
import type { EntityPosition } from "../GameEvent";

const DEBUG = import.meta.env.VITE_DEBUG;

class Projectile implements Entity {
  position: EntityPosition;

  constructor(position: EntityPosition) {
    this.position = position;
  }

  update = (position: EntityPosition) => {
    this.position = position;
  };

  draw = (instance: p5) => {
    instance.push();

    instance.translate(this.position.x, this.position.y);

    instance.push();
    instance.rotate(this.position.theta);
    instance.noStroke();
    instance.fill("#eeeeee");
    instance.circle(8, 0, 10);
    instance.circle(-8, 0, 10);
    instance.rect(-8, -5, 16, 10);
    instance.pop();

    if (DEBUG) {
      instance.stroke("#ff0000");
      instance.noFill();

      instance.push();
      instance.rotate(this.position.theta);
      instance.rect(-13, -5, 26, 10);
      instance.pop();

      instance.line(0, 0, Math.cos(this.position.theta) * 120, Math.sin(this.position.theta) * 120);
      instance.text(`position: (${this.position.x.toFixed(2)}, ${this.position.y.toFixed(2)}), theta: ${this.position.theta.toFixed(2)}`, 0, -85);
    }

    instance.pop();
  };

  destroy = () => {
    // do nothing
  };
}

export default Projectile;
