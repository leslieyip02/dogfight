import { useState } from "react";
import "./App.css";

const API_URL = import.meta.env.VITE_API_URL ?? "/api";

function App() {

  // TODO: add websocket logic

  const [username, setUsername] = useState<string>("testificate");

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
        const playerId = data.playerId;
        console.log(playerId);
      });
  };

  return (
    <>
      <h1>Join a room</h1>
      <form onSubmit={onSubmit}>
        <div>
          <label htmlFor="username">Username</label><br />
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
        <br />
        <button type="submit">Join</button>
      </form>
    </>
  );
}

export default App;
