import { Howl } from "howler";

import config from "./audio.json";

/**
 * A static class which provides audio.
 */
class Audiosheet {
  private static audio: Record<string, Howl> = {};

  /**
   * Loads all sprites from audio.json.
   * Should only be called once when setting up the p5.js canvas.
   */
  static loadAll = async () => {
    Audiosheet.audio = Object.fromEntries(
      Object.entries(config)
        .map(([name, config])=> [name, new Howl({ src: config.path, html5: true })]),
    );
  };

  /**
   * Gets an audio clip.
   * @param name audio clip name
   * @returns audio if it exists, null otherwise
   */
  static get = (name: string): Howl | null => {
    return Audiosheet.audio[name] ?? null;
  };
}

export default Audiosheet;
