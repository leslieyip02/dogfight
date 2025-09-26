import p5 from "p5";
import type { EntityPosition } from "../types/entity";

export interface Entity {
    position: EntityPosition;

    update: (position?: EntityPosition) => void;
    draw: (instance: p5, debug?: boolean) => void;
};
