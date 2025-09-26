import p5 from "p5";
import type { Entity } from "./Entity";
import type { EntityPosition } from "../types/entity";

class Projectile implements Entity {
  position: EntityPosition;

  constructor(position: EntityPosition) {
    this.position = position;
  }

  update = (position?: EntityPosition) => {
    if (!position) {
      return;
    }
    this.position = position;
  };

  remove = () => {};

  draw = (instance: p5, debug?: boolean) => {
    instance.push();
    instance.translate(this.position.x, this.position.y);

    instance.push();
    instance.rotate(this.position.theta);
    instance.noStroke();
    instance.fill("#ffffff");
    instance.circle(-20, 0, 10);
    instance.rect(-20, -5, 20, 10);
    if (debug) {
      instance.stroke("#ff0000");
    }
    instance.circle(0, 0, 10);
    instance.pop();

    if (debug) {
      instance.stroke("#ff0000");
      instance.noFill();
      instance.line(0, 0, Math.cos(this.position.theta) * 120, Math.sin(this.position.theta) * 120);
      instance.text(`position: (${this.position.x.toFixed(2)}, ${this.position.y.toFixed(2)}), theta: ${this.position.theta.toFixed(2)}`, 0, -85);
    }

    instance.pop();
  };
}

export default Projectile;
