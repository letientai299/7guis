import { render, screen } from '@testing-library/react';
import { describe, expect, it } from 'vitest';
import App from './App';

describe('App', () => {
  it('should render without error', async () => {
    render(<App />);
    expect(screen.getByText(/7GUIs/i)).toBeDefined();
  });
});
