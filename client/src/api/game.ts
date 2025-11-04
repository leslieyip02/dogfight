import {
  Event,
  Event_SnapshotEventData,
} from "../pb/event";

export async function fetchSnapshot(host: string): Promise<Event_SnapshotEventData | null> {
  const token = localStorage.getItem("jwt");
  if (!token) {
    return Promise.reject(null);
  }
  return await fetch(`http://${host}/api/room/snapshot?token=${token}`)
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
