import { Howl } from "howler";

export const SOUNDS: Record<string, Howl> = {
  explosionBig: new Howl({
    src: ["explosion-big.wav"],
    html5: true,
  }),
  explosionSmall: new Howl({
    src: ["explosion-small.wav"],
    html5: true,
  }),
  pickup: new Howl({
    src: ["pickup.wav"],
    html5: true,
  }),
  shoot: new Howl({
    src: ["shoot.wav"],
    html5: true,
  }),
  score: new Howl({
    src: ["score.wav"],
    html5: true,
  }),
};
