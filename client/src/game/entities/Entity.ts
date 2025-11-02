import p5 from "p5";

import type { EntityData } from "../../pb/entities";
import type { Vector } from "../../pb/vector";

export interface Entity {
    position: Vector;

    update: (data: EntityData) => void;
    removalAnimationName: () => string | null;

    draw: (instance: p5, debug?: boolean) => void;
    drawIcon?: (instance: p5) => void;
};

export type EntityMap = Record<string, Entity>;
