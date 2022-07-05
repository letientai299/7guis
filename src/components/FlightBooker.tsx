import { useState } from 'react';

enum BookType {
  OneWay,
  RoundTrip,
}

const dateFormat = 'dd.mm.yyyy';

const isLeapYear = (y: number) => {
  return (y % 4 == 0 && y % 100 !== 0) || y % 400 == 0;
};

const parseDate = (s: string): { ok: boolean; date?: Date } => {
  if (s == '') {
    return { ok: true };
  }

  if (s.match(/[^\d\.]/)) {
    return { ok: false };
  }

  const ss = s.split('.');
  if (
    ss.length != 3 ||
    ss[0].length != 2 ||
    ss[1].length != 2 ||
    ss[2].length != 4
  ) {
    return { ok: false };
  }

  const [d, m, y] = ss.map((v) => parseInt(v));
  if (d > 31 || d <= 0) {
    return { ok: false };
  }

  let dMax = 28;
  if (m == 2 && isLeapYear(y)) {
    dMax = 29;
  }
  if ([1, 3, 5, 7, 8, 10, 12].indexOf(m) >= 0) {
    dMax = 31;
  } else if (m != 2) {
    dMax = 30;
  }

  if (d > dMax) {
    return { ok: false };
  }

  const r = new Date(y, m - 1, d);
  return { ok: true, date: r };
};

const FlightBooker = () => {
  const [bookType, setBookType] = useState(BookType.OneWay);
  const [departure, setDeparture] = useState('');
  const [arrival, setArrival] = useState('');

  const bookTypeChange = (i: number) => {
    setBookType(i as BookType);
  };

  const labelClass = 'self-center';
  const canBook = () => {
    if (bookType == BookType.OneWay) {
      return departure != '' && parseDate(departure).ok;
    }

    if (departure == '' || arrival == '') {
      return false;
    }

    const d = parseDate(departure);
    const a = parseDate(arrival);
    if (!d.ok || !a.ok || !d.date || !a.date) {
      return false;
    }

    return a.date.getTime() > d.date.getTime();
  };

  const inform = () => {
    alert('Booked!');
  };

  return (
    <div className={'grid grid-cols-2 gap-2'}>
      <div className={labelClass}>Book Type</div>
      <select onChange={(e) => bookTypeChange(parseInt(e.currentTarget.value))}>
        <option value={BookType.OneWay}>One way</option>
        <option value={BookType.RoundTrip}>Round trip</option>
      </select>

      <div className={labelClass}>Departure</div>
      <input
        type="text"
        onChange={(e) => setDeparture(e.target.value)}
        className={parseDate(departure).ok ? '' : 'bg-red-400'}
        placeholder={dateFormat}
        value={departure}
      />

      <div className={labelClass}>Arrival</div>

      <input
        type="text"
        value={arrival}
        className={parseDate(arrival).ok ? '' : 'bg-red-400'}
        onChange={(e) => setArrival(e.target.value)}
        placeholder={dateFormat}
        disabled={bookType == BookType.OneWay}
      />

      <button
        className={'col-span-2'}
        disabled={!canBook()}
        onClick={() => inform()}
      >
        Book
      </button>
    </div>
  );
};

export default FlightBooker;
