import { useEffect, useState } from "react";

const WS_URL = import.meta.env.VITE_WS_URL;

type Props = {
  clientId: string,
  roomId: string,
}

const Game: React.FC<Props> = ({ clientId, roomId }) => {
  const [socket, setSocket] = useState<WebSocket | null>(null);

  const onopen = () => {
    console.log("open");
  };

  const onclose = () => {
    console.log("close");
    setSocket(null);
  };

  const onmessage = (event: MessageEvent) => {
    console.log("message:", event);
  };

  useEffect(() => { 
    if (socket !== null) {
      return;
    }

    const ws = new WebSocket(`${WS_URL}?clientId=${clientId}&roomId=${roomId}`);
    ws.onopen = onopen;
    ws.onclose = onclose;
    ws.onmessage = onmessage;
    setSocket(ws);
  }, [clientId, roomId, socket]);

  return (
    <>
      <h1>
        Your client ID is <strong>{clientId}</strong>
      </h1>
    </>
  );
};

export default Game;