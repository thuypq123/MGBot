# MGBot

# Discord Bot with AWS SQS Integration

This Discord bot integrates with Amazon Simple Queue Service (SQS) for enhanced message handling and processing.

## Features

- Discord message handling
- Voice command support
- AWS SQS integration for message queueing
- Asynchronous message processing

## Prerequisites

- Go 1.x
- AWS account with SQS access
- Discord Bot Token

## Installation

1. Clone the repository:
   ```
   git clone https://github.com/yourusername/discord-bot-sqs.git
   cd discord-bot-sqs
   ```

2. Install dependencies:
   ```
   go mod tidy
   ```

3. Set up environment variables:
   - `DISCORD_BOT_TOKEN`: Your Discord bot token
   - `AWS_ACCESS_KEY_ID`: Your AWS access key
   - `AWS_SECRET_ACCESS_KEY`: Your AWS secret key
   - `AWS_REGION`: Your AWS region
   - `QUEUE_URL`

4. Update the `QueueURL` in `internal/constant/const.go` with your SQS queue URL.

## Configuration

Modify `internal/constant/const.go` to change the bot's command prefix or other constants:
