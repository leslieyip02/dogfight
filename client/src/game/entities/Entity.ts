import p5 from "p5";

import type { EntityData } from "../../pb/entities";
import type { Vector } from "../../pb/vector";
import type { AnimationStep } from "../graphics/animation";

export interface Entity {
    position: Vector;

    update: (data: EntityData) => void;
    remove: () => AnimationStep | null;

    draw: (instance: p5, debug?: boolean) => void;
    drawIcon?: (instance: p5) => void;
};

export type EntityMap = Record<string, Entity>;
