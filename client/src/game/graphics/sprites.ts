import type p5 from "p5";

import configs from "./sprites.json";

type SpriteConfig = {
  path: string;
  width: number;
  height: number;
  isPlayerSprite?: boolean;
};

/**
 * A static class which provides sprites.
 */
class Spritesheet {
  private static sprites: Record<string, p5.Image[]> = {};
  static isLoaded: boolean = false;

  private static playerSpriteNames: string[] = Object.entries(configs)
    .filter(([, config]: [string, SpriteConfig]) => config.isPlayerSprite)
    .map(([name]) => name);

  /**
   * Loads all sprites from sprites.json.
   * Should only be called once when setting up the p5.js canvas.
   */
  static loadAll = async (instance: p5) => {
    const loadingPromises = Object.entries(configs)
      .map(async ([name, config]) => Spritesheet.loadOne(name, config, instance));

    Spritesheet.sprites = await Promise.all(loadingPromises)
      .then(entries => Object.fromEntries(entries));
    Spritesheet.isLoaded = true;
  };

  /**
   * Loads a single sprite.
   * Will tile the loaded image based on the provided config.
   * @returns a loading promise that resolves once the sprite is loaded
   */
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

  /**
   * Get a single frame of a sprite.
   * @param name sprite name
   * @returns sprite if loaded, null otherwise
   */
  static get = (name: string): p5.Image | null => {
    const frames = Spritesheet.sprites[name];
    if (!frames) {
      return null;
    }
    return frames[0];
  };

  /**
   * Get all frames of a sprite.
   * @param name sprite name
   * @returns frames if loaded, null otherwise
   */
  static getAnimationFrames = (name: string): p5.Image[] | null => {
    return Spritesheet.sprites[name] ?? null;
  };

  /**
   * Gets the path to the player's sprite based on their username.
   * @returns arbitrary sprite path relative to the public directory
   */
  static getPlayerSpritePath = (username: string): string => {
    const spriteName = Spritesheet.getPlayerSpriteName(username);
    return (configs as Record<string, SpriteConfig>)[spriteName].path;
  };

  /**
   * Gets the player's sprite based on their username.
   * @returns arbitrary sprite
   */
  static getPlayerSprite = (username: string): p5.Image | null => {
    const spriteName = Spritesheet.getPlayerSpriteName(username);
    return Spritesheet.get(spriteName);
  };

  /**
   * Returns a arbitrary sprite name by encoding username into an index.
   * @returns arbitrary sprite name
   */
  private static getPlayerSpriteName = (username: string): string => {
    const asciiSum = [...username].reduce((a, b) => a + b.charCodeAt(0), 0);
    return Spritesheet.playerSpriteNames[asciiSum % Spritesheet.playerSpriteNames.length];
  };
};

export default Spritesheet;
