# Skasaha (Discord)

Granblue Fantasy bot for Discord.

## Invite

There is no invite link for this bot, and there will never be.

Skasaha is only available in source code form.

If other people run Skasaha on their servers and provide public invites,
that's their problem.

## Features

* Emoji from [risend/vampy](https://risend.github.io/vampy/).
* Event list.
* Free form search.
  * Events
  * Characters

### Policy

Every feature must be directly related to Granblue Fantasy.

That means generic features will not be implemented.

## Install

```bash
go get github.com/KuroiKitsu/skasaha/...
```

## Configure

Create a file `skasaha.json` with this content:

```json
{
  "token": "YOUR_TOKEN",
  "prefix": "!",
  "emoji_dir": "./media/emoji"
}
```

## Start

```bash
skasaha
```

(`skasaha.json` must be in the working directory)

## Commands

| Command | Arguments | Description |
|---|---|---|
| `help` || Display help. |
| `events` || List of events. |
| `emo`, `emoji` | `name` | Display emoji `name`. |
| `s`, `search` | Free-form text. | Search for something. |

As a special case, emoji have a short form. For example, if you want to
display the emoji `stare`, and your prefix is `!`, then you can do `!!stare`.
If your prefix is `$` then it becomes `$$stare`, and so on.

## Contact

* [Discord](https://discord.gg/E7rky88)

## License

Apache 2.0
