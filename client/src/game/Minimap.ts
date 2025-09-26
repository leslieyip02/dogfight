import type p5 from "p5";

import { BACKGROUND_COLOR } from "./Engine";
import type { Entity } from "./entities/Entity";
import Player from "./entities/Player";
import Powerup from "./entities/Powerup";

const MINIMAP_RADIUS = 100;
const OFFSET = 128;
const MINIMAP_SCALE = 1 / 800;

class Minimap {
  draw = (
    instance: p5,
    clientPlayer: Player,
    entities: { [id: string]: Entity },
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
    Object.values(entities)
      .forEach(entity => {
        if (entity === clientPlayer) {
          return;
        }

        if (entity instanceof Player) {
          instance.fill("#ff0000");
        } else if (entity instanceof Powerup) {
          instance.fill("#00ff00");
        } else {
          return;
        }

        const dx = entity.position.x - clientPlayer.position.x;
        const dy = entity.position.y - clientPlayer.position.y;
        const theta = Math.atan2(dy, dx);
        const distance = Math.min(Math.sqrt(dx * dx + dy * dy) * MINIMAP_SCALE, 1.0) * MINIMAP_RADIUS;
        instance.circle(Math.cos(theta) * distance, Math.sin(theta) * distance, 8);
      });
    instance.pop();
  };
};

export default Minimap;
