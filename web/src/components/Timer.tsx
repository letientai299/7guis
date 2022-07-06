import { useEffect, useRef, useState } from 'react';

const min = 1000;
const max = 20 * 1000;

const Timer = () => {
  const [ms, setMs] = useState(5000);

  // for the Reset to button to work, we use a dummy value that will always
  // change when Reset is clicked, to trigger effect execution
  const [trigger, setTrigger] = useState(0);
  const elapseRef = useRef<HTMLSpanElement>(null);
  const gaugeRef = useRef<HTMLDivElement>(null);

  useEffect(() => {
    let id: number;
    let remain = 0;
    if (elapseRef && elapseRef.current) {
      const s = elapseRef.current.innerText?.replace('s', '');
      if (s) {
        remain = parseFloat(s) * 1000;
      }
    }

    let last = new Date();

    const draw = () => {
      id = requestAnimationFrame(draw);
      const now = new Date();
      remain += now.getTime() - last.getTime();
      last = now;

      if (remain > ms) {
        remain = ms;
        cancelAnimationFrame(id);
      }

      if (elapseRef && elapseRef.current) {
        elapseRef.current.innerText = `${(remain / 1000).toFixed(2)}s`;
      }

      if (gaugeRef && gaugeRef.current) {
        const done = (remain * 100) / ms;
        const w = `${done}%`;
        gaugeRef.current.style.width = w;
      }
    };

    draw();

    return () => {
      cancelAnimationFrame(id);
    };
  }, [ms, trigger]);

  const reset = () => {
    if (elapseRef && elapseRef.current) {
      elapseRef.current.innerText = `0.00s`;
    }
    setTrigger((t) => t + 1);
  };

  return (
    <div className={'grid grid-cols-6 gap-2'}>
      <div className={'self-center'}>Elapsed</div>
      <div className="col-span-4 h-2/3 self-center border border-gray-500 drop-shadow-d rounded-md">
        <div
          ref={gaugeRef}
          className="bg-cyan-500 h-full rounded-md border-2 border-white"
        />
      </div>
      <span ref={elapseRef}>0.00s</span>

      <div className={'self-center'}>Duration</div>
      <input
        type="range"
        min={min}
        max={max}
        value={ms}
        className="col-span-4 bg-cyan-400"
        onChange={(e) => setMs(parseInt(e.target.value))}
      />
      <span>{(ms / 1000).toFixed(2)}s</span>

      <button className={'col-span-full'} onClick={reset}>
        Reset
      </button>
    </div>
  );
};
export default Timer;
