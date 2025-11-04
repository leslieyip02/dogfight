const ROOT_HOST = import.meta.env.VITE_ROOT_HOST;

export type JoinRequest = {
    username: string;
    roomId?: string;
}

export type JoinResponse = {
    clientId: string;
    host: string;
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

  return await fetch(`http://${ROOT_HOST}/api/join`, payload)
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
