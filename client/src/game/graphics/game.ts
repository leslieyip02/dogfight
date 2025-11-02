import p5 from "p5";

import type { EntityMap } from "../entities/Entity";
import Player, { PLAYER_MAX_SPEED } from "../entities/Player";

const DEBUG = import.meta.env.VITE_DEBUG;

const GRID_SIZE = 96;
const MINIMUM_ZOOM = 0.6;
const MAXIMUM_ZOOM = 1.0;

export type CanvasConfig = {
  x: number;
  y: number;
  zoom: number;
};

export function updateCanvasConfig(config: CanvasConfig, clientPlayer: Player) {
  const speed = Math.hypot(clientPlayer.velocity.x, clientPlayer.velocity.y);
  config.x = clientPlayer.position.x;
  config.y = clientPlayer.position.y;
  config.zoom = MAXIMUM_ZOOM - (speed / PLAYER_MAX_SPEED) * (MAXIMUM_ZOOM - MINIMUM_ZOOM);
}

export function drawBackground(config: CanvasConfig, instance: p5) {
  instance.background("#111111");

  instance.push();
  centerCanvas(config, instance);

  const worldLeft = config.x - (window.innerWidth / 2) / config.zoom;
  const worldRight = config.x + (window.innerWidth / 2) / config.zoom;
  const worldTop = config.y - (window.innerHeight / 2) / config.zoom;
  const worldBottom = config.y + (window.innerHeight / 2) / config.zoom;

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

export function drawEntities(config: CanvasConfig, entities: EntityMap, instance: p5) {
  instance.push();
  centerCanvas(config, instance);

  // called separately so that all entities render above their trails
  Object.values(entities)
    .filter(entity => entity instanceof Player)
    .forEach(player => player.drawTrail(instance, DEBUG));

  Object.values(entities)
    .filter(entity => {
      return Math.abs(config.x - entity.position.x) <= window.innerWidth
        && Math.abs(config.y - entity.position.y) <= window.innerHeight;
    })
    .forEach(entity => entity.draw(instance, DEBUG));
  instance.pop();
}

function centerCanvas(config: CanvasConfig, instance: p5) {
  instance.scale(config.zoom);
  instance.translate(
    -config.x + (window.innerWidth / 2) / config.zoom,
    -config.y + (window.innerHeight / 2) / config.zoom,
  );
}
