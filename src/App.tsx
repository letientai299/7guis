import { useState } from 'react';

function App() {
  const [count, setCount] = useState(0);

  return (
    <main className="grid w-screen h-screen place-content-center text-center">
      <p>Hello Vite + React!</p>
      <p>
        Env: <code>{process.env.NODE_ENV}</code>
      </p>
      <p>
        <button
          className={
            'border border-solid border-1 border-amber-400 p-1 rounded-md'
          }
          type="button"
          onClick={() => setCount((count) => count + 1)}
        >
          count: {count}
        </button>
      </p>
      <p>
        Edit <code>App.tsx</code> and save to test HMR updates.
      </p>
    </main>
  );
}

export default App;
