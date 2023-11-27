# Discord bot written in Golang

This is a discord bot with the following functionality:
- Basic commands such as cat pictures, ping and bookmarking.
- Has Integration with GitHub to track changes in repos.
- Can post cool photos of cats
- Searches wiki for short explanations of things

TODOs:
- Implementation of Slash Commands is on its way.
- More GitHub integrations (track new pushes to repo)
- Need to get OpenAI API key to test functionality of ELI5
- Handle errors when message failed to send.

Needed environment variables: BOT_TOKEN, PREFIX, GITHUB_TOKEN, GITHUB_CHANNEL