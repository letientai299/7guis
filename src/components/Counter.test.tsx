import { fireEvent, render, screen } from '@testing-library/react';
import { describe, expect, it } from 'vitest';
import Counter from './Counter';

describe('Counter', () => {
  it('click on button should increase the counter', async () => {
    render(<Counter />);
    const btn = screen.getByRole('button');
    expect(btn).toBeDefined();

    const textBox = screen.getByRole('textbox');
    expect(textBox).toBeDefined();
    expect(textBox).toBeInstanceOf(HTMLInputElement);
    const input = textBox as HTMLInputElement;
    expect(input.value).toContain(0);

    const n = 10;
    for (let i = 0; i < n; i++) {
      fireEvent.click(btn);
    }

    expect(input.value).toContain(n);
  });
});
