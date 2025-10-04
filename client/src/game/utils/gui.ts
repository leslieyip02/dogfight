import type p5 from "p5";

import type { EntityMap } from "../entities/Entity";
import Player, { PLAYER_MAX_SPEED } from "../entities/Player";
import type Input from "./input";
import type { CanvasConfig } from "./sprites";

const MINIMAP_RADIUS = 100;
const MINIMAP_OFFSET = 128;
const MINIMAP_SCALE = 1 / 800;

export function drawMinimap(origin: CanvasConfig, clientPlayer: Player | null, entities: EntityMap, instance: p5) {
  instance.push();
  instance.translate(window.innerWidth - MINIMAP_OFFSET, window.innerHeight - MINIMAP_OFFSET);
  instance.stroke("#ffffff");
  instance.fill("#111111");
  instance.circle(0, 0, MINIMAP_RADIUS * 2);

  Object.values(entities)
    .forEach(entity => {
      const drawIcon = entity.drawIcon;
      if (!drawIcon || entity === clientPlayer) {
        return;
      }

      const dx = entity.position.x - origin.x;
      const dy = entity.position.y - origin.y;
      const theta = Math.atan2(dy, dx);
      const clamped = Math.min(Math.hypot(dx, dy) * MINIMAP_SCALE, 1.0) * MINIMAP_RADIUS;

      instance.push();
      instance.translate(Math.cos(theta) * clamped, Math.sin(theta) * clamped);
      drawIcon(instance);
      instance.pop();
    });

  if (clientPlayer) {
    instance.push();
    instance.rotate(clientPlayer.rotation);
    instance.noStroke();
    instance.fill("#ffffff");
    instance.triangle(
      8, 0,
      -8, 8,
      -8, -8,
    );
    instance.pop();
  }

  instance.pop();
}

export function drawHUD(clientPlayer: Player | null, input: Input, instance: p5) {
  drawSpeedometer(clientPlayer, input, instance);
  drawScore(clientPlayer, instance);
}

export function drawSpeedometer(clientPlayer: Player | null, input: Input, instance: p5) {
  instance.push();
  instance.translate(window.innerWidth - MINIMAP_OFFSET, window.innerHeight - MINIMAP_OFFSET);

  const start = 3 / 4 * instance.PI;
  const end = 7 / 4 * instance.PI;

  instance.push();
  instance.stroke("#ffffff22");
  instance.strokeWeight(4);
  instance.strokeCap(instance.SQUARE);
  instance.noFill();
  instance.arc(0, 0, 220, 220, start, end);
  instance.pop();

  const throttle = Math.hypot(input.mouseX, input.mouseY);
  instance.push();
  instance.stroke("#ffffff");
  instance.strokeWeight(4);
  instance.strokeCap(instance.ROUND);
  instance.noFill();
  instance.arc(0, 0, 220, 220, start, start + throttle * (end - start));
  instance.pop();

  const speed = !clientPlayer ? 0 : Math.hypot(clientPlayer.velocity.x, clientPlayer.velocity.y);
  const interval = 0.16;
  const gap = (Math.PI - interval * 16) / 15;
  instance.push();
  instance.stroke("#ffffff");
  instance.strokeWeight(12);
  instance.strokeCap(instance.SQUARE);
  instance.noFill();
  for (let i = 0; i < speed / PLAYER_MAX_SPEED * 16; i++) {
    instance.stroke(i >= 12 ? "#ec1f26" : "#29cc49");
    instance.arc(0, 0, 248, 248, start + ((interval + gap) * i), start + (interval * (i + 1)) + gap * i);
  }
  instance.pop();

  instance.pop();
}

export function drawScore(clientPlayer: Player | null, instance: p5) {
  const score = `${clientPlayer?.score ?? "?"}`;
  instance.push();
  instance.translate(window.innerWidth - 280, window.innerHeight - 40);
  instance.noFill();
  instance.stroke("#ffffff");
  instance.strokeWeight(1);
  instance.rectMode(instance.CENTER);
  instance.textAlign(instance.CENTER);
  instance.textFont("Courier New");
  instance.text(`score: ${score}`, 0, 0);
  instance.pop();
}
