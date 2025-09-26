import type p5 from "p5";

import type { EntityMap } from "../entities/Entity";
import Player from "../entities/Player";
import { drawInputHelper } from "./input";

const DEBUG = import.meta.env.VITE_DEBUG;

const GRID_SIZE = 96;
const BACKGROUND_COLOR = "#111111";

const MINIMAP_RADIUS = 100;
const MINIMAP_OFFSET = 128;
const MINIMAP_SCALE = 1 / 800;

function centerCanvas(originX: number, originY: number, zoom: number, instance: p5) {
  instance.scale(zoom);
  instance.translate(
    -originX + (window.innerWidth / 2) / zoom,
    -originY + (window.innerHeight / 2) / zoom,
  );
}

export function drawBackground(clientPlayer: Player, zoom: number, instance: p5) {
  instance.background(BACKGROUND_COLOR);

  if (DEBUG) {
    drawInputHelper(instance);
  }

  instance.push();
  const originX = clientPlayer.position.x;
  const originY = clientPlayer.position.y;
  centerCanvas(originX, originY, zoom, instance);

  const worldLeft = originX - (window.innerWidth / 2) / zoom;
  const worldRight = originX + (window.innerWidth / 2) / zoom;
  const worldTop = originY - (window.innerHeight / 2) / zoom;
  const worldBottom = originY + (window.innerHeight / 2) / zoom;

  const startCol = Math.floor(worldLeft / GRID_SIZE) * GRID_SIZE;
  const endCol = Math.ceil(worldRight / GRID_SIZE) * GRID_SIZE;
  const startRow = Math.floor(worldTop / GRID_SIZE) * GRID_SIZE;
  const endRow = Math.ceil(worldBottom / GRID_SIZE) * GRID_SIZE;

  instance.stroke("#ffffff33");
  instance.strokeWeight(2);
  for (let x = startCol; x <= endCol; x += GRID_SIZE) {
    instance.line(x, worldTop, x, worldBottom);
  }
  for (let y = startRow; y <= endRow; y += GRID_SIZE) {
    instance.line(worldLeft, y, worldRight, y);
  }
  instance.pop();
}

export function drawEntities(clientPlayer: Player, entities: EntityMap, zoom: number, instance: p5) {
  instance.push();
  const originX = clientPlayer.position.x;
  const originY = clientPlayer.position.y;
  centerCanvas(originX, originY, zoom, instance);

  // TODO: move this?
  Object.values(entities)
    .filter(entity => entity instanceof Player)
    .forEach(player => player.drawTrail(instance));
  Object.values(entities)
    .forEach(entity => entity.draw(instance, DEBUG));
  instance.pop();
}

export function drawMinimap(clientPlayer: Player, entities: EntityMap, instance: p5) {
  instance.push();
  instance.translate(window.innerWidth - MINIMAP_OFFSET, window.innerHeight - MINIMAP_OFFSET);

  instance.stroke("#ffffff");
  instance.fill(BACKGROUND_COLOR);
  instance.circle(0, 0, MINIMAP_RADIUS * 2);

  instance.push();
  instance.rotate(clientPlayer.position.theta);
  instance.noStroke();
  instance.fill("#ffffff");
  if (clientPlayer.removed) {
    // TODO: different icon when dead
  }
  instance.triangle(
    8, 0,
    -8, 8,
    -8, -8,
  );
  instance.pop();

  instance.fill("#ff0000");
  Object.values(entities)
    .forEach(entity => {
      const drawIcon = entity.drawIcon;
      if (!drawIcon || entity === clientPlayer) {
        return;
      }

      const dx = entity.position.x - clientPlayer.position.x;
      const dy = entity.position.y - clientPlayer.position.y;
      const theta = Math.atan2(dy, dx);
      const clamped = Math.min(Math.hypot(dx, dy) * MINIMAP_SCALE, 1.0) * MINIMAP_RADIUS;

      instance.push();
      instance.translate(Math.cos(theta) * clamped, Math.sin(theta) * clamped);
      drawIcon(instance);
      instance.pop();
    });

  instance.pop();
}
