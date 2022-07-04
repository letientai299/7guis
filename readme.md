# Vite React Typescripts and tools templates

This repo contains the design for my ideal document system.

## Usage

```sh
# extract the repo
$ npx degit https://github.com/letientai299/vite-react-ts-template
# this install dependencies, pnpm and setup git hooks
$ pnpm run setup
# or
$ npm run setup
```

After that, use `pnmp` to

# Added libraries

- [TailwindCSS](https://tailwindcss.com/)

## Configured tools

- Auto format code before commit via
  - [Prettier](https://prettier.io/)
  - [`husky`](https://github.com/typicode/husky)
  - [`lint-staged`](https://github.com/okonet/lint-staged)
- Minify CSS for production via [cssnano](https://cssnano.co/)
- Testing and coverage with `vitest`, `jsdom` and many other libs.
