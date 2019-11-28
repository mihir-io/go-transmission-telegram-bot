Transmission Telegram Bot
=========================


## What is this?
the Transmission Telegram Bot (TTB for short) is a bot that lets you
monitor and control the Transmission BitTorrent client via Telegram
messages.

## How do I use it?

> This section will be updated in the future.

In short, run `transmission-telegram-bot -h` and run the command in the 
background with
the appropriate arguments. While it's running, you can use the following
commands to control it:

* `/commands` - list out the supported commands
* `/list` - list the torrents in the Transmission queue
* `/add` - add a torrent. Takes a URL string to a torrent file.
* `/remove` - removes a torrent, given the torrent ID as an int.
* `/pause` - stops the torrent, given then torrent ID as an int.
* `/play` - starts the torrent, given the torrent ID as an int.
