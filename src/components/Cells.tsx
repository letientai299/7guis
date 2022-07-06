import { useState } from 'react';

class CellData {
  text = ''; // the raw expression
  value = ''; // computed value
  affect: string[] = []; // name of other cells
}

const Cols = 26;
const Rows = 100;

const commonStyle = 'border border-gray-500 p-1';
const thStyle = `font-bold text-right ${commonStyle} bg-zinc-100`;

const cellName = (r: number, c: number) => {
  return `${String.fromCharCode(65 + c)}${r}`;
};

const row = (
  r: number,
  onChange: (name: string, text: string) => void,
  getCellData: (name: string) => CellData,
) => {
  const cells = [
    <th className={thStyle} key={'row-head-' + r}>
      {r}
    </th>,
  ];
  for (let c = 0; c < Cols; c++) {
    const name = cellName(r, c);
    const data = getCellData(name);
    const cell = (
      <td className={commonStyle + ' hover:bg-blue-200'} key={name}>
        <Cell
          text={data.text}
          v={data.value}
          onChange={(s) => onChange(name, s)}
        />
      </td>
    );
    cells.push(cell);
  }
  return <tr key={'row-' + r}>{cells}</tr>;
};

const headerRow = () => {
  const cells = [<th className={thStyle} key={'row-head-0'}></th>];

  let char = 'A'.charCodeAt(0);
  for (let i = 0; i < Cols; i++) {
    const name = String.fromCharCode(char);
    const h = (
      <th className={thStyle} key={`col-header-${name}`}>
        <div className="w-12">{name}</div>
      </th>
    );
    cells.push(h);
    char++;
  }

  return cells;
};

const Cell = (p: {
  text: string;
  v: string;
  onChange: (s: string) => void;
}) => {
  const [editing, setEditing] = useState(false);
  const [content, setContent] = useState(p.text);
  const [modified, setModified] = useState(false);

  const change = (s: string) => {
    setModified(true);
    setContent(s);
  };

  const blur = () => {
    setEditing(false);
    if (modified) {
      p.onChange(content);
    }
    setModified(false);
  };

  return (
    <input
      type="text"
      className="w-24 m-0 border-none shadow-none rounded-none"
      placeholder={editing ? '...' : ''}
      value={editing ? content : p.v}
      onChange={(e) => change(e.target.value)}
      onDoubleClick={() => setEditing(true)}
      readOnly={!editing}
      onBlur={blur}
    ></input>
  );
};

const Cells = () => {
  // raw data as a 2d matrix of text content of all cells
  const [data, setData] = useState(new Map<string, CellData>());

  const findNeededCells = (text: string) => {
    return text.match(/\b[A-Z]\d+\b/g) ?? [];
  };

  const onCellChange = (name: string, text: string) => {
    const cur = data.get(name) ?? new CellData();
    cur.text = text;

    // recompute the cell value if needed
    if (text.startsWith('=')) {
      const needs = findNeededCells(text);
      needs?.forEach((s) => {
        const needCell = data.get(s) ?? new CellData();
        needCell.affect.push(name);
        data.set(s, needCell);
      });

      cur.value = compute(text);
    } else {
      cur.value = text;
    }
    data.set(name, cur);

    // trigger recompute all the affected cells
    cur.affect.forEach((other) => {
      const o = data.get(other) ?? new CellData();
      o.value = compute(o.text);
      data.set(other, o);
    });

    setData(new Map(data));
  };

  const getCellData = (name: string) => {
    return data.get(name) ?? new CellData();
  };

  const compute = (text: string) => {
    console.log('compute for ', text);

    let exp = text.slice(1);
    const needs = findNeededCells(exp);

    needs.forEach((need) => {
      let v = getCellData(need).value;
      if (v == '') {
        v = '0';
      }

      exp = exp.replaceAll(need, v);
    });

    try {
      return eval(exp);
    } catch (error) {
      return `Error! "${text}"`;
    }
  };

  const rows = [...new Array(Rows)].map((_, i) => {
    return row(i + 1, onCellChange, getCellData);
  });

  return (
    <div className="overflow-scroll h-64">
      <table className="border border-collapse min-w-fit m-0">
        <thead className="w-fit">
          <tr>{headerRow()}</tr>
        </thead>
        <tbody>{rows}</tbody>
      </table>
    </div>
  );
};

export default Cells;
