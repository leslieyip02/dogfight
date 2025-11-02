import type p5 from "p5";

import configs from "./sprites.json";

type SpriteConfig = {
  path: string;
  width: number;
  height: number;
  isPlayerSprite?: boolean;
};

class Spritesheet {
  private static sprites: Record<string, p5.Image[]> = {};
  static isLoaded: boolean = false;

  private static playerSpriteNames: string[] = Object.entries(configs)
    .filter(([, config]: [string, SpriteConfig]) => config.isPlayerSprite)
    .map(([name]) => name);

  static loadAll = async (instance: p5) => {
    const loadingPromises = Object.entries(configs)
      .map(async ([name, config]) => Spritesheet.loadOne(name, config, instance));

    Spritesheet.sprites = await Promise.all(loadingPromises)
      .then(entries => Object.fromEntries(entries));
    Spritesheet.isLoaded = true;
  };

  private static loadOne = async (
    name: string,
    config: SpriteConfig,
    instance: p5,
  ): Promise<[string, p5.Image[]]> => {
    const { path, width, height } = config;
    return new Promise<[string, p5.Image[]]>((resolve) =>
      instance.loadImage(path, (image) => {
        const frames = [];
        for (let y = 0; y < image.height; y += height) {
          for (let x = 0; x < image.width; x += width) {
            const frame = image.get(x, y, width, height);
            frames.push(frame);
          }
        }
        resolve([name, frames]);
      }),
    );
  };

  static get = (name: string): p5.Image | null => {
    const frames = Spritesheet.sprites[name];
    if (!frames) {
      return null;
    }
    return frames[0];
  };

  static getAnimationFrames = (name: string): p5.Image[] | null => {
    return Spritesheet.sprites[name] ?? null;
  };

  static getPlayerSpritePath = (username: string): string => {
    const spriteName = Spritesheet.getPlayerSpriteName(username);
    return (configs as Record<string, SpriteConfig>)[spriteName].path;
  };

  static getPlayerSprite = (username: string): p5.Image | null => {
    const spriteName = Spritesheet.getPlayerSpriteName(username);
    return Spritesheet.get(spriteName);
  };

  private static getPlayerSpriteName = (username: string): string => {
    const asciiSum = [...username].reduce((a, b) => a + b.charCodeAt(0), 0);
    return Spritesheet.playerSpriteNames[asciiSum % Spritesheet.playerSpriteNames.length];
  };
};

export default Spritesheet;
