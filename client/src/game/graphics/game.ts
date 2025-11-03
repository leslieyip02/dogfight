import Player, { PLAYER_MAX_SPEED } from "../entities/Player";
import { shouldCullEntity } from "../logic/update";
import type { AnimationStep } from "./animation";
import type { CanvasConfig, GraphicsGameContext } from "./context";

const DEBUG = import.meta.env.VITE_DEBUG;

const GRID_SIZE = 96;
const MINIMUM_ZOOM = 0.6;
const MAXIMUM_ZOOM = 1.0;

export function updateCanvasConfig(config: CanvasConfig, clientPlayer: Player) {
  const speed = Math.hypot(clientPlayer.velocity.x, clientPlayer.velocity.y);
  config.x = clientPlayer.position.x;
  config.y = clientPlayer.position.y;
  config.zoom = MAXIMUM_ZOOM - (speed / PLAYER_MAX_SPEED) * (MAXIMUM_ZOOM - MINIMUM_ZOOM);
}

export function centerCanvas(context: GraphicsGameContext) {
  const { instance, canvasConfig } = context;
  instance.scale(canvasConfig.zoom);
  instance.translate(
    -canvasConfig.x + (window.innerWidth / 2) / canvasConfig.zoom,
    -canvasConfig.y + (window.innerHeight / 2) / canvasConfig.zoom,
  );
}

export function drawBackground(context: GraphicsGameContext) {
  const { instance, canvasConfig } = context;
  instance.background("#111111");

  const worldLeft = canvasConfig.x - (window.innerWidth / 2) / canvasConfig.zoom;
  const worldRight = canvasConfig.x + (window.innerWidth / 2) / canvasConfig.zoom;
  const worldTop = canvasConfig.y - (window.innerHeight / 2) / canvasConfig.zoom;
  const worldBottom = canvasConfig.y + (window.innerHeight / 2) / canvasConfig.zoom;
  const startCol = Math.floor(worldLeft / GRID_SIZE) * GRID_SIZE;
  const endCol = Math.ceil(worldRight / GRID_SIZE) * GRID_SIZE;
  const startRow = Math.floor(worldTop / GRID_SIZE) * GRID_SIZE;
  const endRow = Math.ceil(worldBottom / GRID_SIZE) * GRID_SIZE;

  instance.push();
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

export function drawEntities(context: GraphicsGameContext) {
  Object.values(context.entities)
    .filter(entity => !shouldCullEntity(entity.position, context.canvasConfig))
    .forEach(entity => entity.draw(context.instance, DEBUG));
}

export function drawAnimations(
  context: GraphicsGameContext,
  animations: AnimationStep[],
): AnimationStep[] {
  animations = animations
    .map(animation => animation(context.instance))
    .filter(animation => animation !== null);
  return animations;
}
