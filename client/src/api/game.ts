import {
  Event,
  Event_SnapshotEventData,
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

export function sendEvent(
  socket: WebSocket,
  event: Event,
) {
  if (socket.readyState !== socket.OPEN) {
    return;
  }

  const message = Event.encode(event).finish();
  socket.send(message);
}
