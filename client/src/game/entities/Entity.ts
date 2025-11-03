import p5 from "p5";

import type { EntityData } from "../../pb/entities";
import type { Vector } from "../../pb/vector";
import type { AnimationStep } from "../graphics/animation";

/**
 * Represents a game Entity.
 * Encapsulates updates and rendering for one entity.
 */
export interface Entity {
    position: Vector;

    /**
     * Updates the entity's state to match data.
     * @param data incoming data
     */
    update: (data: EntityData) => void;

    /**
     * Handles effects (e.g. audio) associated with removing the entity.
     * Should only be called if the entity is on screen.
     * @returns removal animation if any
     */
    onRemove: () => AnimationStep | null;

    /**
     * Renders the entity.
     * @param instance p5.js instance
     * @param debug whether to draw debug info
     */
    draw: (instance: p5, debug?: boolean) => void;

    /**
     * Renders the entity's icon (e.g. for minimap).
     * @param instance p5.js instance
     */
    drawIcon?: (instance: p5) => void;
};

/**
 * Represents a record of IDs to entities.
 */
export type EntityMap = Record<string, Entity>;
