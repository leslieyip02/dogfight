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

  remove = () => {};

  // eslint-disable-next-line @typescript-eslint/no-unused-vars, unused-imports/no-unused-vars
  draw = (instance: p5, _debug?: boolean) => {
    instance.push();
    instance.translate(this.position.x, this.position.y);

    instance.rotate(this.position.theta);
    instance.noFill();
    instance.stroke("#ffffff");
    instance.strokeWeight(2);
    for (let i = 0; i < this.points.length; i++) {
      const j = (i + 1) % this.points.length;
      instance.line(this.points[i].x, this.points[i].y, this.points[j].x, this.points[j].y);
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
