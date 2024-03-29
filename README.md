## This branch is working branch for [v3](https://github.com/rumblefrog/source-chat-relay/discussions/44). For the current stable release, the readme is available in [master](https://github.com/rumblefrog/source-chat-relay/blob/1333456609b283a893f6305617818e5a30998181/README.md)

<p align="center">
    <img src="docs/static/logo/cloud-computing.svg" width="250">
</p>

<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->
## Table of Contents

- [This branch is working branch for v3. For the current stable release, the readme is available in master](#this-branch-is-working-branch-for-v3-for-the-current-stable-release-the-readme-is-available-in-master)
- [Table of Contents](#table-of-contents)
- [Introduction](#introduction)
- [Features](#features)
- [Prerequisites](#prerequisites)
- [Getting started](#getting-started)
- [Alliedmods Thread](#alliedmods-thread)
- [Natives](#natives)
- [License](#license)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

## Introduction

<p>
    <img src="https://img.shields.io/github/workflow/status/rumblefrog/source-chat-relay/CI?style=for-the-badge">
    <a href="https://discord.gg/HUc67zN">
        <img src="https://img.shields.io/discord/335290997317697536.svg?style=for-the-badge">
    </a>
    <a href="https://github.com/rumblefrog/source-chat-relay/issues">
        <img src="https://img.shields.io/github/issues/rumblefrog/source-chat-relay.svg?style=for-the-badge">
    </a>
    <a href="https://github.com/rumblefrog/source-chat-relay/blob/master/LICENSE">
        <img src="https://img.shields.io/github/license/rumblefrog/source-chat-relay.svg?style=for-the-badge">
    </a>
    <a href="https://www.patreon.com/bePatron?u=962681">
        <img src="docs/static/become_a_patron_button.png" height="28">
    </a>
</p>

<p align="center">
    <img align="center" src="docs/src/assets/preview_2.gif">
</p>

Communicate between Discord & In-Game, monitor server without being in-game, control the flow of messages and user base engagement!

## Features
 - Receive and send messages bidrectionally
 - Channel and type configuration for powerful setups
 - Setup is incrediblily easy with Discord bot commands and simple config files
 - Upon disconnect, game servers will attempt to reconnect at a fixed interval
 - Filter out certain unwanted messages using regex expressions
 - Set ingame prefixes to send a message with ability to configure flag permission for the prefix
 - Natives to expand upon the functionality of the plugin (Custom events, team chat only relays, etc)

## Prerequisites
 - Server to host the relay binary on (with MySQL if not external)
 - Source game server with the [socket](https://forums.alliedmods.net/showthread.php?t=67640) extension
 - A Discord bot token (https://discordapp.com/developers/applications/)

## Getting started
 - [Setup Guide](https://rumblefrog.me/source-chat-relay/setup)

 > For additional support, feel free to leave a reply on the Alliedmods Thread

## Alliedmods Thread
 - [Thread](https://forums.alliedmods.net/showthread.php?t=311079)

## Natives

Message dispatchers and forwards are available within `client/include/`

## License

This project is licensed under GPL 3.0

Icons made by <a href="https://www.flaticon.com/authors/itim2101" title="itim2101">itim2101</a> from <a href="https://www.flaticon.com/" title="Flaticon">www.flaticon.com</a> is licensed by <a href="http://creativecommons.org/licenses/by/3.0/" title="Creative Commons BY 3.0" target="_blank">CC 3.0 BY</a>
