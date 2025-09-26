import type p5 from "p5";
import type { Event, InputEventData } from "../types/event";

const MOUSE_MAXIMUM_EXTENT = 0.8;
const ZOOM_THRESHOLD = 0.8;
const MINIMUM_ZOOM = 0.8;
const MAXIMUM_ZOOM = 1.0;
const ZOOM_DELTA = 0.005;

function normalizeMouseValue(value: number, full: number): number {
  const delta = value - full / 2;
  const radius = full / 2 * MOUSE_MAXIMUM_EXTENT;
  return Math.sign(delta) * Math.min(Math.abs(delta / radius), 1.0);
};

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
    if (this.socket.readyState !== WebSocket.OPEN) {
      return;
    }

    const full = Math.max(window.innerWidth, window.innerHeight);
    this.mouseX = normalizeMouseValue(instance.mouseX, full);
    this.mouseY = normalizeMouseValue(instance.mouseX, full);

    const data: InputEventData = {
      id: this.clientId,
      mouseX: this.mouseX,
      mouseY: this.mouseY,
      mousePressed: this.mousePressed,
    };
    const event: Event = {
      type: "input",
      data,
    };
    this.socket.send(JSON.stringify(event));

    this.mousePressed = false;
  };

  calculateZoom = () => {
    const throttle = Math.min(Math.hypot(this.mouseX, this.mouseY), 1.0);
    if (throttle > ZOOM_THRESHOLD) {
      this.zoom = Math.max(this.zoom - ZOOM_DELTA, MINIMUM_ZOOM);
    } else {
      this.zoom = Math.min(this.zoom + ZOOM_DELTA, MAXIMUM_ZOOM);
    }
    return this.zoom;
  };
}

export default Input;
