import { fireEvent, render, screen } from '@testing-library/react';
import { describe, expect, it } from 'vitest';
import TemperatureConverter from './TemperatureConverter';

describe('TemperatureConverter', () => {
  it('convert C to F correctly', async () => {
    render(<TemperatureConverter />);
    const c = screen.getByAltText('Celsius');
    expect(c).toBeDefined();
    const f = screen.getByAltText('Fahrenheit');
    expect(f).toBeDefined();

    fireEvent.change(c, { target: { value: 100 } });
    expect((f as HTMLInputElement).value).toEqual('212');
  });

  it('convert F to C correctly', async () => {
    render(<TemperatureConverter />);
    const c = screen.getByAltText('Celsius');
    expect(c).toBeDefined();
    const f = screen.getByAltText('Fahrenheit');
    expect(f).toBeDefined();

    fireEvent.change(f, { target: { value: 122 } });
    expect((c as HTMLInputElement).value).toEqual('50');
  });

  it(`change color when input invalid`, async () => {
    render(<TemperatureConverter />);
    const f = screen.getByAltText('Fahrenheit');
    // note: if the value is '122x' instead, the test won't work,
    // as the received value is somehow just '122'
    fireEvent.change(f, { target: { value: 'x122' } });
    expect(f.className).toContain('bg');
  });
});
