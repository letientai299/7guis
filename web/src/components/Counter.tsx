import React, { useState } from 'react';

const Counter = () => {
  const [count, setCount] = useState(0);
  return (
    <div className={'grid grid-cols-2 grid-rows-1 gap-2'}>
      <div className="col-span-1 self-center">{`Clicked ${count} time(s)`}</div>
      <button className="col-span-1" onClick={() => setCount(count + 1)}>
        Count
      </button>
    </div>
  );
};

export default Counter;
