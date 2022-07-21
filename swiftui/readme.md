# 7GUIs - SwiftUI version

This is my effort to learn SwiftUI via building the 7GUIs challenge.
The tools I'm using are:

- XCode 13.4.1
- Swift 5.6
- [swift-format](https://github.com/apple/swift-format) CLI, because XCode is so
  bad that it can't format code on save. See [Setup formatter](#setup-formatter)

## Setup formatter

I need to enforce code formating to be more readable while working on this.
Because Xcode is so shitty, I need to use some script to use `swift-format` on
the command line with some [custom config](./swift-format.json). The config is
mostly default, only change `lineLength` to 80.

Here's how:

- Install [swift-format](https://github.com/apple/swift-format) is installed
  (why don't swift just include this tool by default?)

  ```sh
  $ brew install swift-format
  ```

- Install [`modd`](https://github.com/cortesi/modd) to watch for file changes.
  Sadly, nodemon doesn't handle this use case well. `modd` will load its config
  from [`modd.conf`](./modd.conf) and execute `./scripts/fmt.sh` on `*.swift`
  files.

- Because `modd` is too sensitive to file changes, we need `--in-place` flag to
  save formatted code, and `--in-place` modifies the file even without any
  changes needed, `scripts/fmt.sh` has to call `swift-format` in 2 passes (if
  there's only a single argument): first pass check whether the file will
  changes again after formatting, if not, stop, otherwise, do the formatting.
  When `modd` picks up the change later due to second pass, it call
  `scripts/fmt.sh` again. This time, the script won't touch the file.

This setup works acceptably if I'm slowing down the editing and saving loop a
bit. Otherwise, XCode keeps complaining about unable to save file due to
external changes.

I'm aware of XCode extension like
[SwiftFormat](https://github.com/nicklockwood/SwiftFormat), but I don't want to
use them since I'm yet to underdand XCode as the time writing this. I'm trying
to use the bare one effectively before customize it.
