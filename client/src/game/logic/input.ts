import type p5 from "p5";

import { type Event, EventType } from "../../pb/event";

const MOUSE_INPUT_RADIUS = Math.min(
  window.innerWidth,
  window.innerHeight,
) / 2 * 0.8;

/**
 * Represents user input.
 */
export type Input = {
  mouseX: number;
  mouseY: number;
  mousePressed: boolean;
};

export function initInput(): Input {
  return {
    mouseX: 0,
    mouseY: 0,
    mousePressed: false,
  };
}

/**
 * Returns a new input given a mouse press.
 * @param current current input
 * @returns new input
 */
export function handleMousePress(
  current: Input,
): Input {
  return {
    ...current,
    mousePressed: true,
  };
}

/**
 * Returns a new input given the current p5.js instance.
 * @param current current input
 * @param instance p5.js instance
 * @returns new input
 */
export function handleMouseMove(
  current: Input,
  instance: p5,
): Input {
  const dx = instance.mouseX - window.innerWidth / 2;
  const dy = instance.mouseY - window.innerHeight / 2;
  const theta = Math.atan2(dy, dx);
  const clamped = Math.min(
    Math.hypot(dx, dy),
    MOUSE_INPUT_RADIUS,
  ) / MOUSE_INPUT_RADIUS;

  return {
    ...current,
    mouseX: Math.cos(theta) * clamped,
    mouseY: Math.sin(theta) * clamped,
  };
}

/**
 * Converts the input into an Event for serialization.
 * Also resets the input.
 * @param current current input
 * @returns new input and the converted event
 */
export function convertInputToEvent(
  current: Input,
  clientId: string,
  isRespawn: boolean,
): [Input, Event | null] {
  if (isRespawn) {
    return [
      initInput(),
      current.mousePressed ? {
        type: EventType.EVENT_TYPE_RESPAWN,
        respawnEventData: {
          id: clientId,
        },
      } : null,
    ];
  };

  return [
    initInput(),
    {
      type: EventType.EVENT_TYPE_INPUT,
      inputEventData: {
        id: clientId,
        ...current,
      },
    },
  ];
}
