
import { PLAYER_MAX_SPEED } from "../entities/Player";
import type { GraphicsGUIContext } from "./context";

const MINIMAP_RADIUS = 100;
const MINIMAP_OFFSET = 128;
const MINIMAP_SCALE = 1 / 800;

export function drawMinimap(context: GraphicsGUIContext) {
  const { instance, entities, canvasConfig, getClientPlayer } = context;
  const clientPlayer = getClientPlayer();

  instance.push();
  instance.translate(window.innerWidth - MINIMAP_OFFSET, window.innerHeight - MINIMAP_OFFSET);
  instance.stroke("#ffffff");
  instance.fill("#111111");
  instance.circle(0, 0, MINIMAP_RADIUS * 2);

  Object.values(entities)
    .forEach(entity => {
      if (!entity.drawIcon || entity === clientPlayer) {
        return;
      }

      const dx = entity.position.x - canvasConfig.x;
      const dy = entity.position.y - canvasConfig.y;
      const theta = Math.atan2(dy, dx);
      const clamped = Math.min(Math.hypot(dx, dy) * MINIMAP_SCALE, 1.0) * MINIMAP_RADIUS;

      instance.push();
      instance.translate(Math.cos(theta) * clamped, Math.sin(theta) * clamped);
      entity.drawIcon(instance);
      instance.pop();
    });

  if (clientPlayer) {
    instance.push();
    instance.rotate(clientPlayer.rotation);
    instance.noStroke();
    instance.fill("#ffffff");
    instance.triangle(8, 0, -8, 8, -8, -8);
    instance.pop();
  }

  instance.pop();
}

export function drawHUD(context: GraphicsGUIContext) {
  drawSpeedometer(context);
  drawScore(context);
}

export function drawSpeedometer(context: GraphicsGUIContext) {
  const { instance, getClientPlayer, getInput } = context;
  const clientPlayer = getClientPlayer();
  const input = getInput();

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

export function drawScore(context: GraphicsGUIContext) {
  const { instance, getClientPlayer } = context;
  const clientPlayer = getClientPlayer();

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

export function drawRespawnPrompt(context: GraphicsGUIContext) {
  const { instance, getClientPlayer } = context;
  const clientPlayer = getClientPlayer();
  if (clientPlayer) {
    return;
  }

  instance.fill("#11111155");
  instance.rect(0, 0, window.innerWidth, window.innerHeight);

  instance.push();
  instance.textFont("Courier New");
  instance.textAlign(instance.CENTER);
  instance.translate(window.innerWidth / 2, window.innerHeight / 2);

  instance.stroke("#ffffff");
  instance.fill("#ffffff");
  instance.push();
  instance.textSize(32);
  instance.text("splashed!", 0, -32);
  instance.pop();

  instance.push();
  instance.textSize(16);
  instance.text("click to respawn", 0, 8);
  instance.pop();

  instance.pop();
}
