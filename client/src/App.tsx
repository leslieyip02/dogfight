import "./App.css";

import { useState } from "react";
import { generateUsername } from "unique-username-generator";

import { joinRoom } from "./api/room";
import Game from "./components/Game";

const App = () => {
  const [username, setUsername] = useState<string>(generateUsername("-"));
  const [roomId, setRoomId] = useState<string>("");
  const [clientId, setClientId] = useState<string | null>(null);
  const [token, setToken] = useState<string | null>(null);

  const onSubmit = async (event: React.FormEvent<HTMLFormElement>) => {
    event.preventDefault();
    await joinRoom(username, roomId)
      .then(response => {
        setClientId(response.clientId);
        setToken(response.token);
      });
  };

  if (!clientId || !token) {
    return (
      <>
        <form onSubmit={onSubmit}>
          <h1 className="form__header">dogfight</h1>
          <div className="form__field">
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
          <div className="form__field">
            <label htmlFor="roomId">room id:</label><br />
            <input
              type="text"
              id="room id"
              name="roomId"
              aria-required="true"
              value={roomId}
              onChange={(e) => setRoomId(e.target.value)}
            />
          </div>
          <button className="form__submit" type="submit">Join</button>
        </form>
      </>
    );
  }

  return <Game clientId={clientId} token={token} />;
};

export default App;
