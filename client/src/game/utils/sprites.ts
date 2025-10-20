import type p5 from "p5";

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
  delta: {
    path: "delta.png",
    width: 96,
    height: 96,
  },
  explosionBig: {
    path: "explosion-big.png",
    width: 192,
    height: 192,
  },
  explosionSmall: {
    path: "explosion-small.png",
    width: 96,
    height: 96,
  },
};

export type CanvasConfig = {
  x: number;
  y: number;
  zoom: number;
};

export type Spritesheet = Record<string, p5.Image[]>;

export async function loadSpritesheet(instance: p5): Promise<Spritesheet> {
  const loadingPromises = Object.entries(SPRITESHEET_CONFIGS)
    .map(async ([name, config]) => {
      return new Promise<[string, p5.Image[]]>((resolve) => {
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
