import React from 'react';
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
  }[] = [
    { name: 'Counter', com: Counter },
    { name: 'Temperature Converter', com: TemperatureConverter },
    { name: 'Flight Booker', com: FlightBooker },
    { name: 'Timer', com: Timer },
    { name: 'CRUD', com: CRUD },
    { name: 'Circle Drawer', com: CircleDrawer },
    // { name: 'Cells', com: <Cells /> },
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
            <strong className="col-span-1">{g.name}</strong>
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
