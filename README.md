<h1 align="center"> Source Chat Relay </h1> <br>

<p align="center">
    <img src="assets/logo/cloud-computing.svg">
</p>

<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->
## Table of Contents

- [Introduction](#introduction)
- [Features](#features)
- [Prerequisites](#prerequisites)
- [Getting started](#getting-started)
- [Bot Commands](#bot-commands)
- [Credits](#credits)
- [License](#license)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

## Introduction

![Travis (.com)](https://img.shields.io/travis/com/rumblefrog/source-chat-relay.svg?style=for-the-badge)
![Discord](https://img.shields.io/discord/443915420324331521.svg?style=for-the-badge)
![GitHub issues](https://img.shields.io/github/issues/rumblefrog/source-chat-relay.svg?style=for-the-badge)
![GitHub](https://img.shields.io/github/license/rumblefrog/source-chat-relay.svg?style=for-the-badge)

Communicate between Discord & In-Game, monitor server without being in-game, control the flow of messages and user base engagement!

## Features
 - Bidirectional - You can receive/send on both Discord and game server side!
 - Channel - Get creative! With the ability to specify receive/send channels, you can even send messages to other game servers, not only to Discord!
 - Ease of use - Setup is incrediblily easy with Discord bot commands and simple config files
 - Reliable - Upon disconnect, game servers will attempt to reconnect at a fixed interval

## Prerequisites
 - Server to host the relay binary on (with MySQL if not external)
 - Source game server with the socket extension
 - A Discord bot token

## Getting started
 1. Head over to [releases](https://github.com/rumblefrog/source-chat-relay/releases) and download the latest package for your operating system
 2. Relay Server

    1. Upload the binary to the server
    2. Configure `config.toml.example` and rename it to `config.toml`
    3. Start it by running `./server`

3. Game server

    1. Upload the plugin to `addons/sourcemod/plugins`
    2. Load the plugin
    3. Configure the config at `cfg/sourcemod/Source-Chat-Relay.cfg` with the host/port of the relay server
    4. Reload the plugin

4. If all goes properly, you may view the registered entities via `r/entities`
5. Configure the receive/send channels via `r/receivechannels EntityID Channel (comma delimited)` and `r/sendchannels EntityID Channel (comma delimited)` respectively
6. If all done correctly, you will be able to receive and send messages base on the channels you set it to

## Bot Commands

Before any clients can send messages, you must set the receive/send channels on them

 - r/receivechannel #TextChannel/EntityID Channel
    - #TextChannel/EntityID - Either a text channel mention via #channel-name or a game server entity ID obtainable via `r/entities`
    - Channel - List of channels to receive at. If more than one, you may use comma to separate them
 - r/sendchannel #TextChannel/EntityID Channel
    - #TextChannel/EntityID - Either a text channel mention via #channel-name or a game server entity ID obtainable via `r/entities`
    - Channel - List of channels to send to. If more than one, you may use comma to separate them
 - r/entities ?channel/server
    - ?channel/server - Optional argument, allows you filter the return results by either channel or server type

## Credits
 - Ron for being a human linter

## License

This project is licensed under GPL 3.0
 