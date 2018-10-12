The filter is useful if you wish to drop certain messages whether it be commands or just unsolicited messages

To enable filters:

1. Set `Filter` to `true` in `config.toml`
2. Create a `filter.txt` in the same directory as the relay server binary
3. Populate each line of the file with a regex expression (**Note**: if you have a blank line, it will filter everything out)
