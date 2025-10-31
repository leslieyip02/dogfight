import p5 from "p5";

import type { Entity as EntityData } from "../../pb/entities";
import type { Vector } from "../../pb/vector";
import type { BaseEntity } from "./Entity";

class Asteroid implements BaseEntity {
  position: Vector;
  rotation: number;

  points: Vector[];

  constructor(data: EntityData) {
    if (!data.position || !data.rotation) {
      throw new Error(`expected entity data but got ${data}`);
    }
    this.position = data.position;
    this.rotation = data.rotation;

    const asteroidData = data.asteroidData;
    if (!asteroidData) {
      throw new Error(`expected asteroid data but got ${data}`);
    }
    this.points = asteroidData.points;
  }

  update = (data: EntityData) => {
    if (!data.position || !data.rotation) {
      return;
    }
    this.position = data.position;
    this.rotation = data.rotation;
  };

  removalAnimationName = () => {
    return "explosionBig";
  };

  draw = (instance: p5, debug?: boolean) => {
    instance.push();
    instance.translate(this.position.x, this.position.y);
    instance.rotate(this.rotation);
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
