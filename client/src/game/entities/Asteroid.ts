import p5 from "p5";

import type { EntityPosition } from "../types/entity";
import type { Entity } from "./Entity";

class Asteroid implements Entity {
  position: EntityPosition;

  points: { x: number, y: number }[];

  constructor(position: EntityPosition, points: { x: number, y: number }[]) {
    this.position = position;
    this.points = points;
  }

  update = (position?: EntityPosition) => {
    if (!position) {
      return;
    }
    this.position = position;
  };

  removalAnimationName = () => {
    return "bigExplosion";
  };

  draw = (instance: p5, debug?: boolean) => {
    instance.push();
    instance.translate(this.position.x, this.position.y);
    instance.rotate(this.position.theta);
    instance.fill("#ffffff88");
    instance.stroke("#ffffff");
    instance.strokeWeight(4);
    instance.beginShape();
    for (let i = 0; i < this.points.length; i++) {
      instance.vertex(this.points[i].x, this.points[i].y);
    }
    instance.endShape(instance.CLOSE);

    instance.noFill();
    instance.stroke("#ffffff");
    instance.strokeWeight(1);
    instance.beginShape(instance.TRIANGLE_STRIP);
    for (let i = 0; i < this.points.length; i++) {
      instance.vertex(this.points[i].x, this.points[i].y);
    }
    instance.endShape();

    if (debug) {
      instance.push();
      instance.stroke("#ff0000");
      instance.fill("#ffffff");
      instance.circle(0, 0, 10);
      instance.pop();
    }

    instance.pop();
  };

  drawIcon = (instance: p5) => {
    instance.push();
    instance.fill("#0000ff");
    instance.circle(0, 0, 8);
    instance.pop();
  };
}

export default Asteroid;
