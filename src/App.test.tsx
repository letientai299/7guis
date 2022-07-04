import { describe, expect, it } from 'vitest';
import App from './App';
import { fireEvent, render, screen } from '@testing-library/react';

describe('App', () => {
  it('should render without error', async () => {
    render(<App />);
    expect(screen.getByText(/Env/i)).toBeDefined();
    expect(screen.getByText(/count/i)).toBeDefined();
  });

  it('should render without error 2', async () => {
    render(<App />);

    const btn = screen.getByRole('button');
    expect(btn).toBeDefined();
    expect(btn.textContent ?? '', 'before clicking').toEqual('count: 0');
    const n = 10;
    for (let i = 0; i < n; i++) {
      fireEvent.click(btn);
    }
    expect(btn.textContent ?? '').toEqual(`count: ${n}`);
  });
});
