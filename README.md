# README

<p align="center">
    <img src="assets/logo/cloud-computing.svg" width="250">
</p>

<p>
    <a href="https://travis-ci.com/rumblefrog/source-chat-relay">
        <img src="https://img.shields.io/travis/com/rumblefrog/source-chat-relay.svg?style=for-the-badge">
    </a>
    <a href="https://discord.gg/TZ4BsrQ">
        <img src="https://img.shields.io/discord/443915420324331521.svg?style=for-the-badge">
    </a>
    <a href="https://github.com/rumblefrog/source-chat-relay/issues">
        <img src="https://img.shields.io/github/issues/rumblefrog/source-chat-relay.svg?style=for-the-badge">
    </a>
    <a href="https://github.com/rumblefrog/source-chat-relay/blob/master/LICENSE">
        <img src="https://img.shields.io/github/license/rumblefrog/source-chat-relay.svg?style=for-the-badge">
    </a>
    <a href="https://www.patreon.com/bePatron?u=962681">
        <img src="assets/become_a_patron_button.png" height="28">
    </a>
</p>

<img src="assets/preview.gif">

Communicate between Discord & In-Game, monitor server without being in-game, control the flow of messages and user base engagement!

## Features

* Receive and send messages bidirectionally
* Channel configuration for powerful setups
* Setup is incredibly easy with Discord bot commands and simple config files
* Upon disconnect, game servers will attempt to reconnect at a fixed interval
* Filter out certain unwanted messages using regex expressions
* Set in-game prefixes to send a message with ability to configure flag permission for the prefix

## Prerequisites

* Server to host the relay binary on \(with MySQL if not external\)
* Source game server with the [socket](https://forums.alliedmods.net/showthread.php?t=67640) extension
* A Discord bot token \([https://discordapp.com/developers/applications/](https://discordapp.com/developers/applications/)\)

## Recommended hosts

Many people ask me what host this can work on, so here are some referrals

* [Vultr](https://www.vultr.com/?ref=7553630)
* [Digitalocean](https://m.do.co/c/87ffbbdddbe9)

