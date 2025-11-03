import type p5 from "p5";

import { Vector } from "../../pb/vector";
import type Player from "../entities/Player";
import Spritesheet from "./sprites";

const PLAYER_TRAIL_LENGTH = 24;

export type AnimationStep = (instance: p5) => AnimationStep | null;

export function generateAnimation(
  body: AnimationStep,
): AnimationStep | null {
  const step = (): AnimationStep | null => {
    return (instance: p5) => {
      return body(instance);
    };
  };
  return step();
}

export function generateExplosionAnimation(
  name: string,
  position: Vector,
): AnimationStep | null {
  const frames = Spritesheet.getAnimationFrames(name);
  if (!frames) {
    return null;
  }

  let index = 0;
  const x = position.x - frames[index].width / 2;
  const y = position.y - frames[index].height / 2;
  const body = (instance: p5): AnimationStep | null => {
    if (index >= frames.length) {
      return null;
    }

    instance.push();
    instance.translate(x, y);
    instance.image(frames[index], 0, 0);
    instance.pop();
    index++;
    return body;
  };

  return generateAnimation(body);
}

export function generatePlayerTrailAnimation(
  playerRef: WeakRef<Player>,
): AnimationStep | null {
  let index = 0;
  const positions: Vector[] = new Array<Vector>(PLAYER_TRAIL_LENGTH);

  const body = (instance: p5): AnimationStep | null => {
    const player = playerRef.deref();
    if (!player) {
      return null;
    }

    positions[index] = {
      x: player.position.x,
      y: player.position.y,
    };
    index = (index + 1) % PLAYER_TRAIL_LENGTH;

    const color = instance.color("#ffa320");
    instance.push();
    instance.strokeWeight(4);
    for (let i = 0; i < PLAYER_TRAIL_LENGTH - 1; i++) {
      const current = positions[(index + i) % PLAYER_TRAIL_LENGTH];
      const next = positions[(index + i + 1) % PLAYER_TRAIL_LENGTH];
      if (!current || !next) {
        continue;
      }

      color.setAlpha(Math.min(i/(PLAYER_TRAIL_LENGTH / 4), 1) * 255);
      instance.stroke(color);
      instance.line(current.x, current.y, next.x, next.y);
    }
    instance.pop();

    return body;
  };

  return generateAnimation(body);
}
