import {
  Event,
  type Event_InputEventData,
  Event_RespawnEventData,
  Event_SnapshotEventData,
  EventType,
} from "../pb/event";

const API_URL = import.meta.env.VITE_API_URL;

export async function fetchSnapshot(): Promise<Event_SnapshotEventData | null> {
  const token = localStorage.getItem("jwt");
  if (!token) {
    return Promise.reject(null);
  }
  return await fetch(`${API_URL}/room/snapshot?token=${token}`)
    .then(response => response.arrayBuffer())
    .then(buffer => {
      const message = new Uint8Array(buffer);
      return Event.decode(message).snapshotEventData ?? null;
    });
}

export function sendInputMessage(
  socket: WebSocket,
  data: Event_InputEventData,
) {
  if (socket.readyState !== socket.OPEN) {
    return;
  }

  const message = Event.encode({
    type: EventType.EVENT_TYPE_INPUT,
    inputEventData: data,
  }).finish();
  socket.send(message);
}

export function sendRespawnMessage(
  socket: WebSocket,
  data: Event_RespawnEventData,
) {
  if (socket.readyState !== socket.OPEN) {
    return;
  }

  const message = Event.encode({
    type: EventType.EVENT_TYPE_RESPAWN,
    respawnEventData: data,
  }).finish();
  socket.send(message);
}
