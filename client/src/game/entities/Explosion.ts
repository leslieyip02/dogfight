import p5 from "p5";
import type { Entity } from "./Entity";
import type { EntityPosition } from "../GameEvent";

class Explosion implements Entity {
  position: EntityPosition;
  diameter: number;
  opacity: number;

  constructor(position: EntityPosition) {
    this.position = position;
    this.diameter = 10;
    this.opacity = 1;
  }

  update = () => {
    this.diameter *= 1.1;
    this.opacity *= 0.95;
  };

  draw = (instance: p5) => {
    instance.push();

    instance.translate(this.position.x, this.position.y);

    instance.stroke(255, 255 * this.opacity);
    instance.strokeWeight(2);
    instance.noFill();
    instance.circle(0, 0, this.diameter);

    instance.pop();
  };
}

export default Explosion;
