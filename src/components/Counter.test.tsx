import { fireEvent, render, screen } from '@testing-library/react';
import { describe, expect, it } from 'vitest';
import Counter from './Counter';

describe('Counter', () => {
  it('click on button should increase the counter', async () => {
    render(<Counter />);
    const btn = screen.getByRole('button');
    expect(btn).toBeDefined();

    const div = screen.getByText(/Click/i);
    expect(div).toBeDefined();
    expect(div.textContent).toContain(0);

    const n = 10;
    for (let i = 0; i < n; i++) {
      fireEvent.click(btn);
    }

    expect(div.textContent).toContain(n);
  });
});
