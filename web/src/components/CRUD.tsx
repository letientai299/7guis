import { useEffect, useState } from 'react';

type Person = {
  first: string;
  last: string;
};

const initData: Person[] = [...new Array(50)].map((_, i) => {
  return { first: `${i} John`, last: `Doe ${i}` };
});

const CRUD = () => {
  const [data, setData] = useState<Person[]>(initData);
  const [first, setFirst] = useState('');
  const [last, setLast] = useState('');
  const [selected, setSelected] = useState<Person>({ first: '', last: '' });
  const [filterPrefix, setFilterPrefix] = useState('');
  const [filtered, setFiltered] = useState<Person[]>(data);

  useEffect(() => {
    const r = data.filter(
      (d) =>
        d.first.startsWith(filterPrefix) || d.last.startsWith(filterPrefix),
    );
    setFiltered(r);
  }, [filterPrefix, data]);

  const select = (p: Person) => {
    setSelected(p);
  };

  useEffect(() => {
    const i = filtered.indexOf(selected);
    if (i < 0) {
      return;
    }

    const row = document.getElementById('selected-record');
    const div = document.getElementById('data-table');
    if (!row || !div) {
      return;
    }

    const pos = row.offsetHeight * i;
    if (div.scrollTop < pos && div.scrollTop + div.offsetHeight >= pos) {
      return; // no need to scroll
    }

    div.scrollTop = row.offsetHeight * i;
  }, [selected, filtered]);

  const create = () => {
    const p: Person = { first: first, last: last };
    setData([...data, p]);
    setSelected(p);
  };

  const update = () => {
    const i = filtered.indexOf(selected);
    if (i >= 0) {
      const ui = data.indexOf(selected);
      if (ui >= 0) {
        const p: Person = { first: first, last: last };
        data[ui] = p;
        setData([...data]);
        setSelected(p);
      }
    }
  };

  const remove = () => {
    const i = filtered.indexOf(selected);
    if (i >= 0) {
      const ui = data.indexOf(selected);
      if (ui >= 0) {
        data.splice(ui, 1);
        setSelected({ first: '', last: '' });
        setData([...data]);
      }
    }
  };

  return (
    <div className="flex gap-2 flex-col">
      <div className="flex gap-2">
        <div>Filter prefix</div>
        <input
          type="text"
          placeholder="name..."
          value={filterPrefix}
          onChange={(e) => setFilterPrefix(e.target.value)}
        />
      </div>

      <div className="flex gap-2 content-start">
        <div className="basis-2/3 h-64 overflow-scroll" id="data-table">
          <table className="m-0">
            <tbody>
              {filtered.map((d, i) => {
                const highlight = d == selected ? 'text-white bg-cyan-500' : '';
                return (
                  <tr onClick={() => select(d)} key={i} id="selected-record">
                    <td
                      className={`border-b border-gray-500 ` + highlight}
                    >{`${d.first}, ${d.last}`}</td>
                  </tr>
                );
              })}
            </tbody>
          </table>
        </div>

        <div className="basis-1/3 grid grid-cols-1 grow-0 h-fit">
          <div>Name</div>
          <input
            type="text"
            value={first}
            onChange={(e) => setFirst(e.target.value)}
          />
          <div>Surname</div>
          <input
            type="text"
            value={last}
            onChange={(e) => setLast(e.target.value)}
          />
        </div>
      </div>

      <div className="flex gap-2">
        <button
          className="grow"
          disabled={first == '' && last == ''}
          onClick={() => create()}
        >
          Create
        </button>
        <button
          className="grow"
          disabled={
            (first == '' && last == '') || filtered.indexOf(selected) < 0
          }
          onClick={() => update()}
        >
          Update
        </button>
        <button
          className="grow"
          disabled={filtered.indexOf(selected) < 0}
          onClick={() => remove()}
        >
          Delete
        </button>
      </div>
    </div>
  );
};
export default CRUD;
