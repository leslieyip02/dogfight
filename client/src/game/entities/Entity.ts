import p5 from "p5";

import type { EntityData } from "../types/entity";
import type { Vector } from "../types/geometry";

export interface Entity {
    position: Vector;

    update: (data: EntityData) => void;
    removalAnimationName: () => string | null;

    draw: (instance: p5, debug?: boolean) => void;
    drawIcon?: (instance: p5) => void;
};

export type EntityMap = {
    [id: string]: Entity,
}
