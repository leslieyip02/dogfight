import p5 from "p5";

import type { EntityData } from "../../pb/entities";
import type { Vector } from "../../pb/vector";
import { type AbilityFlag,isAbilityActive, WIDE_BEAM_ABILITY_FLAG } from "../logic/abilities";
import type { Entity } from "./Entity";

const PROJECTILE_WIDTH = 20;

class Projectile implements Entity {
  position: Vector;
  rotation: number;
  flags: AbilityFlag;
  lifetime: number;

  constructor(data: EntityData) {
    if (!data.position || !data.velocity) {
      throw new Error(`expected entity data but got ${data}`);
    }
    this.position = data.position;
    this.rotation = data.rotation;

    const projectileData = data.projectileData;
    if (!projectileData) {
      throw new Error(`expected projectile data but got ${data}`);
    }
    this.flags = projectileData.flags;
    this.lifetime = projectileData.lifetime;
  }

  update = (data: EntityData) => {
    if (!data.position || !data.rotation) {
      return;
    }
    this.position = data.position;
    this.rotation = data.rotation;

    const projectileData = data.projectileData;
    if (projectileData) {
      this.lifetime = projectileData.lifetime;
    }
  };

  removalAnimationName = (): string => {
    // only explode on hit
    return this.lifetime > 1 ? "explosionSmall" : "";
  };

  draw = (instance: p5, debug?: boolean) => {
    instance.push();
    instance.translate(this.position.x, this.position.y);

    instance.push();
    instance.rotate(this.rotation);

    if (isAbilityActive(this.flags, WIDE_BEAM_ABILITY_FLAG)) {
      instance.noFill();
      instance.stroke("#ffffff");
      instance.strokeWeight(8);
      instance.arc(-40, 0, 60, 80, 3 / 2 * Math.PI + 0.5, 1 / 2 * Math.PI - 0.5);
    } else {
      instance.noStroke();
      instance.fill("#ffffff");
      instance.circle(-20, 0, 10);
      instance.rect(-20, -5, 20, 10);
      instance.circle(0, 0, 10);
    }

    instance.pop();

    if (debug) {
      instance.push();
      instance.stroke("#ff0000");
      instance.noFill();

      instance.push();
      instance.rotate(this.rotation);
      instance.rect(-PROJECTILE_WIDTH / 2, -PROJECTILE_WIDTH / 2, PROJECTILE_WIDTH, PROJECTILE_WIDTH);
      instance.pop();

      instance.line(0, 0, Math.cos(this.rotation) * 120, Math.sin(this.rotation) * 120);
      instance.text(`position: (${this.position.x.toFixed(2)}, ${this.position.y.toFixed(2)}), rotation: ${this.rotation.toFixed(2)}`, 0, -85);
      instance.pop();
    }

    instance.pop();
  };
}

export default Projectile;
