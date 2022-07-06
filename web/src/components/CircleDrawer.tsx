import { useRef, useState } from 'react';

type Circle = {
  x: number;
  y: number;
  r: number;
};

const defaultRadius = 20;

const CircleDrawer = () => {
  type Action = (data: Circle[]) => void; // mutation on data list
  const [actions, setActions] = useState<Action[]>([]);
  const [actionNum, setActionNum] = useState(0);

  const data: Circle[] = []; // init empty data, execute action
  for (let i = 0; i < actionNum; i++) {
    actions[i](data);
  }

  const [selected, setSelected] = useState(-1);
  const sliderRef = useRef<HTMLDivElement>(null);

  const addAction = (a: Action) => {
    if (actionNum >= actions.length) {
      actions.push(a);
    } else {
      actions[actionNum] = a;
    }
    setActionNum(actionNum + 1);
    setActions([...actions]);
  };

  const clear = () => {
    deactivate();
    addAction((ds) => ds.splice(0, ds.length));
  };

  const undo = () => {
    deactivate();
    if (actionNum > 0) {
      setActionNum(actionNum - 1);
    }
  };

  const redo = () => {
    deactivate();
    if (actionNum < actions.length) {
      setActionNum(actionNum + 1);
    }
  };

  const addCircle = (e: React.MouseEvent<SVGElement, MouseEvent>) => {
    const svg = e.target as SVGElement;
    const rect = svg.getBoundingClientRect();

    const c: Circle = {
      r: defaultRadius,
      x: Math.round(e.clientX - rect.x),
      y: Math.round(e.clientY - rect.y),
    };

    const existed =
      data.filter((o) => {
        return o.x == c.x && o.y == c.y && o.r == c.r;
      }).length != 0;

    if (!existed) {
      addAction((ds) => ds.push(c));
    }
  };

  const draw = (c: Circle, i: number) => {
    const fill = i == selected ? 'fill-red-400' : 'fill-transparent';
    return (
      <circle
        onClick={(e) => activate(i, e)}
        className={`hover:fill-gray-400 ${fill} hover:opacity-90`}
        key={JSON.stringify(c)}
        cx={c.x}
        cy={c.y}
        r={c.r}
        fill={'none'}
        stroke={'gray'}
      ></circle>
    );
  };

  const deactivate = () => {
    if (!sliderRef || !sliderRef.current) {
      return;
    }

    setSelected(-1);
    const slider = sliderRef.current;
    slider.style.visibility = 'hidden';
  };

  const activate = (i: number, e: React.MouseEvent<SVGElement, MouseEvent>) => {
    e.stopPropagation();
    if (!sliderRef || !sliderRef.current) {
      return;
    }

    const c = data[i];
    const slider = sliderRef.current;
    const sliderRect = slider.getBoundingClientRect();
    slider.style.visibility = 'visible';
    slider.style.top = c.y + c.r + sliderRect.height + 'px';
    slider.style.left = c.x - sliderRect.width / 2 + 'px';

    const input = slider.getElementsByTagName('input')[0];
    input.value = c.r + '';

    input.oninput = () => {
      const circle = e.target as SVGCircleElement;
      circle.setAttribute('r', input.value);
    };

    input.onchange = () => {
      const r = parseInt(input.value);
      addAction((ds) => {
        ds[i] = { x: ds[i].x, y: ds[i].y, r: r };
      });
    };

    setSelected(i);
  };

  return (
    <div className="flex flex-col gap-2 relative">
      <div className="grid grid-cols-3 gap-2">
        <button disabled={actionNum == 0} onClick={() => undo()}>
          Undo
        </button>
        <button disabled={actionNum == actions.length} onClick={() => redo()}>
          Redo
        </button>
        <button disabled={data.length == 0} onClick={() => clear()}>
          Clear {data.length ? data.length : ''}
        </button>
      </div>

      <svg
        className="w-full h-64 border border-gray-500 rounded-md"
        onClick={(e) => addCircle(e)}
      >
        {data.map((c, i) => draw(c, i))}
      </svg>

      <div
        ref={sliderRef}
        className="flex flex-row gap-2 border rounded-md border-gray-500 p-2 w-fit absolute invisible"
      >
        <div className={'self-center'}>Radius</div>
        <input type="range" min={10} max={100} />
        <button onClick={() => deactivate()}>OK</button>
      </div>
    </div>
  );
};

export default CircleDrawer;
