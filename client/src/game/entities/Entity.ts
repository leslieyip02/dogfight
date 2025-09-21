import p5 from "p5";
import type { EntityPosition } from "../GameEvent";

export interface Entity {
    position: EntityPosition;

    update: (position?: EntityPosition) => void;
    // TODO: add debug argument
    draw: (instance: p5) => void;
    remove: () => void;
};
