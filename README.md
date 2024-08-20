# Serverless cloud functions library

The repository contains cloud functions and related packages for different tasks.

## Telegram notification function

The function triggers Telegram Bot API to deliver messages/notifications to specified channels.

- Tested on GCP
- Trigger: Pub/Sub
- Dependencies:
    - pkg/telegram
- Environment variables:
    - `BOT_TOKEN` - a secutiry bot token given by [BotFather](https://t.me/BotFather)
    - `DEFAULT_CHAT_ID` - specifies the default chat/channel ID to deliver the messages
    - `{CHANNEL}_CHAT_ID` - (optional) can be used to deliver messages to a particular chat or channel
