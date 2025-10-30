const API_URL = import.meta.env.VITE_API_URL;

export type JoinRequest = {
    username: string;
    roomId?: string;
}

export type JoinResponse = {
    clientId: string;
    token: string;
}

export async function joinRoom(username: string, roomId: string): Promise<JoinResponse> {
  const body: JoinRequest = {
    username,
    ...(roomId === "" ? {} : { roomId: roomId }),
  };
  const payload = {
    method: "POST",
    body: JSON.stringify(body),
  };

  return await fetch(`${API_URL}/room/join`, payload)
    .then(async response => {
      if (!response.ok) {
        const message = await response.text();
        throw new Error(message || "unknown error");
      }
      return response.json();
    })
    .then(response => {
      localStorage.setItem("jwt", response.token);
      return response;
    });
}
