import "./App.css";

import { useState } from "react";

import { joinRoom } from "./api/room";
import Game from "./components/Game";

const App = () => {
  const [username, setUsername] = useState<string>("testificate");
  const [clientId, setClientId] = useState<string | null>(null);
  const [token, setToken] = useState<string | null>(null);

  const onSubmit = async (event: React.FormEvent<HTMLFormElement>) => {
    event.preventDefault();
    await joinRoom(username)
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

  return <Game clientId={clientId} token={token} />;
};

export default App;
