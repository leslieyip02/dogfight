import { Howl } from "howler";

export const SOUNDS: Record<string, Howl> = {
  explosionBig: new Howl({
    src: ["audio/explosion/big.wav"],
    html5: true,
  }),
  explosionSmall: new Howl({
    src: ["audio/explosion/small.wav"],
    html5: true,
  }),
  pickup: new Howl({
    src: ["audio/player/pickup.wav"],
    html5: true,
  }),
  shoot: new Howl({
    src: ["audio/player/shoot.wav"],
    html5: true,
  }),
  score: new Howl({
    src: ["audio/player/score.wav"],
    html5: true,
  }),
};
