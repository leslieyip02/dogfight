import { Howl } from "howler";

import config from "./audio.json";

class Audiosheet {
  private static audio: Record<string, Howl> = {};

  static loadAll = async () => {
    Audiosheet.audio = Object.fromEntries(
      Object.entries(config)
        .map(([name,config])=> [name, new Howl({ src:config.path,html5:true })]),
    );
  };

  static get = (name: string): Howl | null => {
    return Audiosheet.audio[name] ?? null;
  };
}

export default Audiosheet;
