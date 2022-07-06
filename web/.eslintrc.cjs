module.exports = {
  env: {
    browser: true,
    es2021: true,
    node: true,
  },
  ignorePatterns: ['dist'],
  parser: '@typescript-eslint/parser',
  extends: [
    'plugin:@typescript-eslint/recommended',
    'plugin:testing-library/react',
    'prettier',
  ],
  plugins: ['@typescript-eslint', 'testing-library', 'html', 'prettier'],
  rules: {
    'prettier/prettier': 'error',
  },
};
