import { JoinRequest, JoinResponse } from "../pb/join";

const ROOT_HOST = import.meta.env.VITE_ROOT_HOST;

export async function joinRoom(username: string, roomId: string): Promise<JoinResponse> {
  const body: JoinRequest = {
    username,
    ...(roomId === "" ? {} : { roomId: roomId }),
  };
  const payload = {
    method: "POST",
    body: JoinRequest.encode(body).finish(),
  };

  return await fetch(`http://${ROOT_HOST}/api/join`, payload)
    .then(async response => {
      if (!response.ok) {
        const message = await response.text();
        throw new Error(message);
      }
      return response.arrayBuffer();
    })
    .then(buffer => {
      const message = new Uint8Array(buffer);
      const joinResponse = JoinResponse.decode(message);
      localStorage.setItem("jwt", joinResponse.token);
      return joinResponse;
    });
}
