import type p5 from "p5";

import { sendInputMessage, sendRespawnMessage } from "../../api/game";
import type { Event_InputEventData, Event_RespawnEventData } from "../../pb/event";
import { SOUNDS } from "./sounds";

const MOUSE_INPUT_RADIUS = Math.min(window.innerWidth, window.innerHeight) / 2 * 0.8;

function normalizeMouseValues(mouseX: number, mouseY: number): [number, number] {
  const dx = mouseX - window.innerWidth / 2;
  const dy = mouseY - window.innerHeight / 2;
  const theta = Math.atan2(dy, dx);
  const clamped = Math.min(Math.hypot(dx, dy), MOUSE_INPUT_RADIUS) / MOUSE_INPUT_RADIUS;
  return [Math.cos(theta) * clamped, Math.sin(theta) * clamped];
}

class Input {
  clientId: string;
  socket: WebSocket;

  mouseX: number;
  mouseY: number;
  mousePressed: boolean;

  constructor(clientId: string, socket: WebSocket) {
    this.clientId = clientId;
    this.socket = socket;

    this.mouseX = 0;
    this.mouseY = 0;
    this.mousePressed = false;
  }

  handleMousePress = () => {
    this.mousePressed = true;
  };

  handleInput = (instance: p5) => {
    [this.mouseX, this.mouseY] = normalizeMouseValues(instance.mouseX, instance.mouseY);
    const data: Event_InputEventData = {
      id: this.clientId,
      mouseX: this.mouseX,
      mouseY: this.mouseY,
      mousePressed: this.mousePressed,
    };
    sendInputMessage(this.socket, data);

    if (this.mousePressed) {
      SOUNDS["shoot"].play();
    }
    this.mousePressed = false;
  };

  handleRespawn = () => {
    if (!this.mousePressed) {
      return;
    }

    const data: Event_RespawnEventData = {
      id: this.clientId,
    };
    sendRespawnMessage(this.socket, data);

    this.mousePressed = false;
  };
}

export default Input;
