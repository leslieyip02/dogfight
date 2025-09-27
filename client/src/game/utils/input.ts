import type p5 from "p5";

import { sendMessage } from "../../api/game";
import type { InputEventData, RespawnEventData } from "../types/event";

const MINIMUM_ZOOM = 0.8;
const MAXIMUM_ZOOM = 1.0;
const ZOOM_THRESHOLD = 1.0;
const ZOOM_DELTA = 0.005;

const MOUSE_INPUT_RADIUS = Math.min(window.innerWidth, window.innerHeight) / 2 * 0.8;

function normalizeMouseValues(mouseX: number, mouseY: number): [number, number] {
  const dx = mouseX - window.innerWidth / 2;
  const dy = mouseY - window.innerHeight / 2;
  const theta = Math.atan2(dy, dx);
  const clamped = Math.min(Math.hypot(dx, dy), MOUSE_INPUT_RADIUS) / MOUSE_INPUT_RADIUS;
  return [Math.cos(theta) * clamped, Math.sin(theta) * clamped];
}

export function drawInputHelper(instance: p5) {
  instance.push();
  instance.noFill();
  instance.stroke("#0000ff");
  instance.strokeWeight(2);
  instance.circle(window.innerWidth / 2, window.innerHeight / 2, MOUSE_INPUT_RADIUS * 2);
  instance.pop();
}

class Input {
  clientId: string;
  socket: WebSocket;

  mouseX: number;
  mouseY: number;
  mousePressed: boolean;
  zoom: number;

  constructor(clientId: string, socket: WebSocket) {
    this.clientId = clientId;
    this.socket = socket;

    this.mouseX = 0;
    this.mouseY = 0;
    this.mousePressed = false;
    this.zoom = MAXIMUM_ZOOM;
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

  calculateZoom = () => {
    const throttle = Math.hypot(this.mouseX, this.mouseY);
    if (throttle >= ZOOM_THRESHOLD) {
      this.zoom = Math.max(this.zoom - ZOOM_DELTA, MINIMUM_ZOOM);
    } else {
      this.zoom = Math.min(this.zoom + ZOOM_DELTA, MAXIMUM_ZOOM);
    }
    return this.zoom;
  };
}

export default Input;
