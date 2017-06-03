# skasaha

Granblue Fantasy bot for Discord.

This bot exists as a free and open source alternative to
the proprietary [risend/vampy][] bot.

[risend/vampy]: <https://risend.github.io/vampy/>

## Features

* Emoji from [risend/vampy][].
* Event list.
* Character lookup.

## Install

```bash
go get github.com/KuroiKitsu/skasaha/...
```

## Configure

Create a file `skasaha.json` with this content:

```json
{
  "token": "YOUR_TOKEN",
  "prefix": "!"
}
```

## Start

```bash
skasaha
```

## Commands

| Command | Arguments | Description |
|---|---|---|
| `help` || Display help. |
| `events` || List of events. |
| `emo`, `emoji` | `name` | Displays emoji `name`. |
| `s`, `search` | Any text. | Searches for something. |

As a special case, emoji have a short form. For example, if you want to
display the emoji `stare`, and your prefix is `!`, then you can do `!!stare`.
If your prefix is `$` then it becomes `$$stare`, and so on.

## License

Apache 2.0
