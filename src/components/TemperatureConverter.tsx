import { useEffect, useState } from 'react';

const TempInput = (p: {
  name: string;
  v: number;
  onChange: (n: number) => void;
}) => {
  const [value, setValue] = useState(`${p.v}`);
  useEffect(() => {
    setValue(`${p.v}`);
  }, [p.v]);

  const handleChange = (s: string) => {
    setValue(s);
    const n = Number.parseFloat(s ? s : '0');

    if (n == null || Number.isNaN(n)) {
      return;
    }

    p.onChange(n);
  };

  return (
    <label className="grid grid-cols-2 basis-1/2 grow gap-1">
      <input
        className={isNaN(Number(value)) ? 'bg-red-400' : ''}
        alt={p.name}
        type="text"
        value={value}
        onChange={(e) => handleChange(e.target.value)}
      />
      <div className="self-center">{p.name}</div>
    </label>
  );
};

const TemperatureConverter = () => {
  const [c, setC] = useState(0);
  const [f, setF] = useState(32);

  const cChange = (c: number) => {
    const v = c * (9 / 5) + 32;
    setF(Math.round(v * 100) / 100);
  };

  const fChange = (f: number) => {
    const v = ((f - 32) * 5) / 9;
    setC(Math.round(v * 100) / 100);
  };

  return (
    <div className="flex gap-2">
      <TempInput name="Celsius" onChange={cChange} v={c} />
      <div className="self-center grow-0">=</div>
      <TempInput name="Fahrenheit" onChange={fChange} v={f} />
    </div>
  );
};

export default TemperatureConverter;
