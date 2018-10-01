# Source Chat Relay

A self-hosted Discord <=> Source game relay

## Features
 - Bidirectional - You can receive/send on both Discord and game server side!
 - Channel - Get creative! With the ability to specify receive/send channels, you can even send messages to other game servers, not only to Discord!
 - Ease of use - Setup is incrediblily easy with Discord bot commands and simple config files
 - Reliable - Upon disconnect, game servers will attempt to reconnect at a fixed interval

## Requirement
 - A MySQL server
 - A game server with the socket extension
 - A server to host the relay server on
 - A Discord bot token

## Getting started
 1. Head over to [Releases](https://github.com/rumblefrog/source-chat-relay/releases) and download the latest package for your operating system
 2. Relay Server

    1. Upload the binary to the server
    2. Configure `config.toml.example` and rename it to `config.toml`
    3. Start it by running `./server`

3. Game server

    1. Upload the plugin to `addons/sourcemod/plugins`
    2. Load the plugin
    3. Configure the config at `cfg/sourcemod/Source-Chat-Relay.cfg` with the host/port of the relay server
    4. Reload the plugin

4. If all goes properly, the client should connect to the relay server and authenticate itself

## Bot Commands

    Before any clients can send messages, you must set the receive/send channels on them

 - r/receivechannel #Channel/EntityID Channel
 - r/sendchannel #Channel/EntityID Channel
 - r/entities ?channel/server

## Credits
 - Ron for being a human linter

## License

This project is licensed under GPL 3.0
 