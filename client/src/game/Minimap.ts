import type p5 from "p5";
import { BACKGROUND_COLOR } from "./GameEngine";
import type Player from "./entities/Player";
import type Powerup from "./entities/Powerup";

const MINIMAP_RADIUS = 100;
const OFFSET = 128;
const MINIMAP_SCALE = 1 / 800;

class Minimap {
  draw = (
    instance: p5,
    clientPlayer: Player,
    players: { [id: string]: Player },
    powerups: { [id: string]: Powerup },
  ) => {
    instance.push();
    instance.translate(window.innerWidth - OFFSET, window.innerHeight - OFFSET);

    instance.stroke("#ffffff");
    instance.fill(BACKGROUND_COLOR);
    instance.circle(0, 0, MINIMAP_RADIUS * 2);

    instance.push();
    instance.rotate(clientPlayer.position.theta);
    instance.noStroke();
    instance.fill("#ffffff");

    if (clientPlayer.removed) {
      // TODO: different icon when dead
    }
    instance.triangle(
      8, 0,
      -8, 8,
      -8, -8,
    );
    instance.pop();

    instance.fill("#ff0000");
    Object.values(players)
      .forEach(player => {
        if (player === clientPlayer) {
          return;
        }

        const dx = player.position.x - clientPlayer.position.x;
        const dy = player.position.y - clientPlayer.position.y;
        const theta = Math.atan2(dy, dx);
        const distance = Math.min(Math.sqrt(dx * dx + dy * dy) * MINIMAP_SCALE, 1.0) * MINIMAP_RADIUS;
        instance.circle(Math.cos(theta) * distance, Math.sin(theta) * distance, 8);
      });

    instance.fill("#00ff00");
    Object.values(powerups)
      .forEach(powerup => {
        const dx = powerup.position.x - clientPlayer.position.x;
        const dy = powerup.position.y - clientPlayer.position.y;
        const theta = Math.atan2(dy, dx);
        const distance = Math.min(Math.sqrt(dx * dx + dy * dy) * MINIMAP_SCALE, 1.0) * MINIMAP_RADIUS;
        instance.circle(Math.cos(theta) * distance, Math.sin(theta) * distance, 8);
      });

    instance.pop();
  };
};

export default Minimap;
