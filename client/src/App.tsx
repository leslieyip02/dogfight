import "./App.css";

function App() {

  // TODO: add websocket logic

  return (
    <>
      <h1>Join a room</h1>
      <form>
        <div>
          <label htmlFor="username">Username</label><br />
          <input type="text" id="username" name="username" required aria-required="true" />
        </div>
        <br />
        <div>
          <label htmlFor="roomId">Room ID</label><br />
          <input type="text" id="roomId" name="roomId" />
        </div>
        <br />
        <button type="submit">Submit</button>
      </form>
    </>
  );
}

export default App;
