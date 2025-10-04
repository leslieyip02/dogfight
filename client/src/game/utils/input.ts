import type p5 from "p5";

import { sendMessage } from "../../api/game";
import type { InputEventData, RespawnEventData } from "../types/event";

const MOUSE_INPUT_RADIUS = Math.min(window.innerWidth, window.innerHeight) / 2 * 0.8;

function normalizeMouseValues(mouseX: number, mouseY: number): [number, number] {
  const dx = mouseX - window.innerWidth / 2;
  const dy = mouseY - window.innerHeight / 2;
  const theta = Math.atan2(dy, dx);
  const clamped = Math.max(Math.min(Math.hypot(dx, dy), MOUSE_INPUT_RADIUS) / MOUSE_INPUT_RADIUS, 0.1);
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
    const data: InputEventData = {
      id: this.clientId,
      mouseX: this.mouseX,
      mouseY: this.mouseY,
      mousePressed: this.mousePressed,
    };
    sendMessage(this.socket, "input", data);
    this.mousePressed = false;
  };

  handleRespawn = () => {
    if (!this.mousePressed) {
      return;
    }

    const data: RespawnEventData = {
      id: this.clientId,
    };
    sendMessage(this.socket, "respawn", data);
    this.mousePressed = false;
  };
}

export default Input;
