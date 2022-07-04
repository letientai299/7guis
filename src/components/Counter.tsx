import React, { useState } from 'react';

const Counter = (p: React.HTMLAttributes<HTMLElement>) => {
  const [count, setCount] = useState(0);
  const { className, ...rest } = p;
  return (
    <div className={'grid grid-cols-2 grid-rows-1'} {...rest}>
      <input
        className="col-span-1"
        type="text"
        disabled
        value={`Clicked ${count} time(s)`}
      />
      <button className="col-span-1" onClick={() => setCount(count + 1)}>
        Count
      </button>
    </div>
  );
};

export default Counter;
