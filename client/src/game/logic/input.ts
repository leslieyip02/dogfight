import type p5 from "p5";

import { type Event, EventType } from "../../pb/event";

const MOUSE_INPUT_RADIUS = Math.min(window.innerWidth, window.innerHeight) / 2 * 0.8;

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

export function handleMousePress(current: Input, mousePressed: boolean): Input {
  return {
    ...current,
    mousePressed: current.mousePressed || mousePressed,
  };
}

export function handleMouseMove(current: Input, instance: p5): Input {
  const dx = instance.mouseX - window.innerWidth / 2;
  const dy = instance.mouseY - window.innerHeight / 2;
  const theta = Math.atan2(dy, dx);
  const clamped = Math.min(Math.hypot(dx, dy), MOUSE_INPUT_RADIUS) / MOUSE_INPUT_RADIUS;

  return {
    ...current,
    mouseX: Math.cos(theta) * clamped,
    mouseY: Math.sin(theta) * clamped,
  };
}

export function convertInputToEvent(current: Input, clientId: string, isRespawn: boolean): [Input, Event | null] {
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
