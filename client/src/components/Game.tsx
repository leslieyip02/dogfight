import { useEffect, useLayoutEffect, useRef, useState } from "react";

import "./Game.css";
import GameEngine from "../game/GameEngine";
import p5 from "p5";

const WS_URL = import.meta.env.VITE_WS_URL;

type Props = {
  clientId: string,
  roomId: string,
}

const Game: React.FC<Props> = ({ clientId, roomId }) => {
  const gameEngineRef = useRef<GameEngine | null>(null);
  const containerRef = useRef<HTMLDivElement>(null);
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
    gameEngineRef.current?.receive(event);
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

  useLayoutEffect(() => {
    const sketch = (instance: p5) => {
      gameEngineRef.current = new GameEngine(instance);
    };
    
    const instance = new p5(sketch, containerRef.current!);
    return () => instance.remove();
  }, []);

  return <div className="game__container" ref={containerRef} />;
};

export default Game;