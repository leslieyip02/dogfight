import type p5 from "p5";
import type { Image } from "p5";

import type { EntityMap } from "../entities/Entity";
import Player from "../entities/Player";
import { drawInputHelper } from "./input";

const DEBUG = import.meta.env.VITE_DEBUG;

const GRID_SIZE = 96;
const BACKGROUND_COLOR = "#111111";

const MINIMAP_RADIUS = 100;
const MINIMAP_OFFSET = 128;
const MINIMAP_SCALE = 1 / 800;

const SPRITESHEET_CONFIGS = {
  alpha: {
    path: "alpha.png",
    width: 96,
    height: 96,
  },
  bravo: {
    path: "bravo.png",
    width: 96,
    height: 96,
  },
  charlie: {
    path: "charlie.png",
    width: 96,
    height: 96,
  },
  explosion: {
    path: "explosion.png",
    width: 96,
    height: 96,
  },
};

export type CanvasConfig = {
  x: number;
  y: number;
  zoom: number;
};

export type Spritesheet = Record<string, Image[]>;

export async function loadSpritesheet(instance: p5): Promise<Spritesheet> {
  const loadingPromises = Object.entries(SPRITESHEET_CONFIGS)
    .map(async ([name, config]) => {
      return new Promise<[string, Image[]]>((resolve) => {
        const { path, width, height } = config;
        instance.loadImage(path, (image) => {
          const frames = [];
          for (let y = 0; y < image.height; y += height) {
            for (let x = 0; x < image.width; x += width) {
              const frame = image.get(x, y, width, height);
              frames.push(frame);
            }
          }
          resolve([name, frames]);
        });
      });
    });
  return Promise.all(loadingPromises)
    .then(entries => Object.fromEntries(entries));
}

export function drawBackground(config: CanvasConfig, instance: p5) {
  instance.background(BACKGROUND_COLOR);

  if (DEBUG) {
    drawInputHelper(instance);
  }

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

  // TODO: move this?
  Object.values(entities)
    .filter(entity => entity instanceof Player)
    .forEach(player => player.drawTrail(instance));
  Object.values(entities)
    .forEach(entity => entity.draw(instance, DEBUG));
  instance.pop();
}

export function drawMinimap(origin: CanvasConfig, clientPlayer: Player | null, entities: EntityMap, instance: p5) {
  instance.push();
  instance.translate(window.innerWidth - MINIMAP_OFFSET, window.innerHeight - MINIMAP_OFFSET);
  instance.stroke("#ffffff");
  instance.fill(BACKGROUND_COLOR);
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
    instance.rotate(clientPlayer.position.theta);
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

export function drawRespawnPrompt(instance: p5) {
  instance.push();
  instance.textFont("Courier New");
  instance.textAlign(instance.CENTER);
  instance.stroke("#ffffff");
  instance.fill("#ffffff");
  instance.translate(window.innerWidth / 2, window.innerHeight / 2);

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

function centerCanvas(config: CanvasConfig, instance: p5) {
  instance.scale(config.zoom);
  instance.translate(
    -config.x + (window.innerWidth / 2) / config.zoom,
    -config.y + (window.innerHeight / 2) / config.zoom,
  );
}
