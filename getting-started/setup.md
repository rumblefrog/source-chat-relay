---
description: Setup instructions
---

# Setup

## Prerequisites

Make sure you have the following before proceeding with the installation

* Machine to host the relay binary on \(with MySQL if not external\) \[Port forward/Firewall if necessary\]
* Source game server with the [socket](https://forums.alliedmods.net/showthread.php?t=67640) extension
* A Discord bot token \([https://discordapp.com/developers/applications/](https://discordapp.com/developers/applications/)\)

## Download

Download the latest release from [releases](https://github.com/rumblefrog/source-chat-relay/releases) for either linux/windows

> Note: If you need additional OS/Architectural support, you may build the server binary yourself

## Installation

* Upload the binary file \(server\[.exe\]\) to the server
* Upload the plugin to the Sourcemod plugins folder and load the plugin by either restart, change map, and/or via `sm plugins load Source-Chat-Relay`

## Configuration

* Configure the `config.toml.example` on the relay server and rename it to `config.toml`
* Configure `cfg/sourcemod/Source-Chat-Relay.cfg` on the game server

## Connecting

This is the part that everyone gets confused on, so please follow it carefully.

### How the system works \(Prelude\)

Before actually connecting everything, you should familarize yourself with how the system actually works.

The concept is similar to TV **channels**, you can broadcast on a channel and you can view/receive on a channel, with the only difference is that you can may view/receive on multiple channels at once.

Entities are bidirectional \(both direction send/receive\), meaning they are both capable of broadcasting and view/receiving.

### Connecting

If all goes well and the game server successfully connects to the relay server, you should see an entity from each game server when you run `r/entities` as a command to the Discord bot \(If you don't, see [Troubleshooting](../support/troubleshooting.md)\)

You may now began linking things, copy the entity ID of the game server and use `r/receivechannel` and `r/sendchannel` to configure the receive and send properties

Examples:

* `r/receivechannel !@#$%^klmwxyz01bc2345nopqr 1`
* `r/sendchannel !@#$%^klmwxyz01bc2345nopqr 1`

Similarily, you can do the same for Discord text channel

Examples:

* `r/receivechannel #general 1`
* `r/sendchannel #general 1`

Just like, two entities are configured to both send and receive on channel `1`

More in depth explanation of the commands is available at [Bot Commands](bot-commands.md)

