import type p5 from "p5";

import type { EntityData } from "../../pb/entities";
import type { Vector } from "../../pb/vector";
import type { BaseEntity } from "./Entity";

class Animation implements BaseEntity {
  position: Vector;

  frames: p5.Image[];
  index: number;
  onFinish: () => void;

  constructor(position: Vector, frames: p5.Image[], onFinish: () => void) {
    this.position = position;
    this.frames= frames;
    this.index = 0;
    this.onFinish = onFinish;
  }

  // eslint-disable-next-line @typescript-eslint/no-unused-vars, unused-imports/no-unused-vars
  update = (_data: EntityData) => {};

  removalAnimationName = () => {
    return null;
  };

  draw = (instance: p5) => {
    if (this.index >= this.frames.length) {
      this.onFinish();
      return;
    }

    instance.push();
    const frame = this.frames[this.index];
    instance.translate(this.position.x - frame.width / 2, this.position.y - frame.height / 2);
    instance.image(this.frames[this.index], 0, 0);
    instance.pop();

    this.index++;
  };
}

export default Animation;
