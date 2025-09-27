import type { Event, EventType, InputEventData, RespawnEventData, SnapshotEventData } from "../game/types/event";

const API_URL = import.meta.env.VITE_API_URL;

export async function fetchSnapshot(): Promise<SnapshotEventData | null> {
  const token = localStorage.getItem("jwt");
  if (!token) {
    return Promise.reject(null);
  }
  return await fetch(`${API_URL}/room/snapshot?token=${token}`)
    .then(response => response.json());
}

export function sendMessage(socket: WebSocket, type: EventType, data: RespawnEventData | InputEventData) {
  if (socket.readyState !== socket.OPEN) {
    return;
  }

  const event: Event = {
    type,
    data: data,
  };
  socket.send(JSON.stringify(event));
}
