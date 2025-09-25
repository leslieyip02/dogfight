import { useState } from "react";
import Game from "./components/Game";

import "./App.css";

const API_URL = import.meta.env.VITE_API_URL;

const App = () => {
  const [username, setUsername] = useState<string>("testificate");
  const [clientId, setClientId] = useState<string | null>(null);
  const [roomId, setRoomId] = useState<string | null>(null);

  const onSubmit = async (event: React.FormEvent<HTMLFormElement>) => {
    event.preventDefault();

    const body = {
      "username": username,
    };

    const payload = {
      method: "POST",
      body: JSON.stringify(body),
    };

    await fetch(`${API_URL}/room/join`, payload)
      .then(response => response.json())
      .then(data => {
        setClientId(data.clientId);
        setRoomId(data.roomId);
      });
  };

  if (!clientId || !roomId) {
    return (
      <>
        <form onSubmit={onSubmit}>
          <h1 className="form__header">dogfight</h1>
          <div className="form__username">
            <label htmlFor="username">username:</label><br />
            <input
              type="text"
              id="username"
              name="username"
              required
              aria-required="true"
              value={username}
              onChange={(e) => setUsername(e.target.value)}
            />
          </div>
          <button className="form__submit" type="submit">Join</button>
        </form>
      </>
    );
  }

  return <Game clientId={clientId} roomId={roomId} />;
};

export default App;
