import React from 'react';
import Cells from './components/Cells';
import CircleDrawer from './components/CircleDrawer';
import Counter from './components/Counter';
import CRUD from './components/CRUD';
import FlightBooker from './components/FlightBooker';
import TemperatureConverter from './components/TemperatureConverter';
import Timer from './components/Timer';

function App() {
  const guis: {
    name: string;
    com: React.FC<React.HTMLProps<HTMLElement>>;
    remark?: string;
  }[] = [
    { name: 'Counter', com: Counter },
    { name: 'Temperature Converter', com: TemperatureConverter },
    { name: 'Flight Booker', com: FlightBooker },
    { name: 'Timer', com: Timer },
    { name: 'CRUD', com: CRUD },
    {
      name: 'Circle Drawer',
      com: CircleDrawer,
      remark: "Still can't maintain the radius slider while undo/redo",
    },
    {
      name: 'Cells',
      com: Cells,
      remark: `
Not perfect, no shortcut keys, rerender irrelevant cells, doesn't
support functions, ranges, ...
`,
    },
  ];

  return (
    <main className="prose m-auto">
      <h1>
        <a href="https://eugenkiss.github.io/7guis/tasks" target="_blank">
          7GUIs
        </a>
      </h1>
      {guis.map((g) => {
        return (
          <section
            key={g.name}
            className="grid grid-cols-4 h-fit border-b p-3 gap-2"
          >
            <div className="flex flex-col">
              <strong>{g.name}</strong>
              {g.remark ? <div className="italic">{g.remark}</div> : <></>}
            </div>
            <div className="col-span-3 rounded-md border border-gray-500 p-2">
              <g.com />
            </div>
          </section>
        );
      })}
    </main>
  );
}

export default App;
