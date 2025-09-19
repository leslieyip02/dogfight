import type p5 from "p5";
import { BACKGROUND_COLOR } from "./GameEngine";
import type Player from "./entities/Player";

const MINIMAP_RADIUS = 96;
const OFFSET = 128;

class Minimap {
  draw = (
    instance: p5,
    clientPlayer: Player,
    players: { [id: string]: Player },
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
        instance.circle(Math.cos(theta) * MINIMAP_RADIUS, Math.sin(theta) * MINIMAP_RADIUS, 8);
      });

    instance.pop();
  };
};

export default Minimap;
