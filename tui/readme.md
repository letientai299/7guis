# TUI implementation in Go of 7GUIs task

## Developing

For quick edit and review loop, use `nodemon` to watch for code change, rebuild
and restart the application via:

```sh
$ nodemon -w . -e .go -x './restart.sh'
```

## Lesson learned

- `tview.Form` overrides each input fields style. Hence, after Temperature
  Converter, for other form component, I rather use the field directly to
  support color changes during validation.

- `tview` doesn't have a reactive system, hence, it's very painful to implement
  these exercises in it. Or, perhaps I'm just too stupid to figure out the
  correct way to implement them.
