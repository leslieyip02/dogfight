import p5 from "p5";
import type { EntityPosition } from "../types/entity";

export interface Entity {
    position: EntityPosition;

    update: (position?: EntityPosition) => void;
    remove: () => void;

    draw: (instance: p5, debug?: boolean) => void;
    drawIcon?: (instance: p5) => void;
};

export type EntityMap = {
    [id: string]: Entity,
}
