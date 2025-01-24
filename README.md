![image](https://github.com/goldomatic/shipo-cli/blob/main/shipo-cli-logo.png)


# shipo-CLI

Welcome to **shipo-CLI**: The ultimate shitposting command-line tool for the modern shitposter! This tool lets you unleash your creativity (or chaos) across Bluesky and Twitter… because what’s life without a good shitpost?

---

## Features
- **Post to Bluesky**: Because one platform isn’t enough for your brilliant ideas.
- **Post to Twitter**: Keep the chaos alive on X (formerly Twitter).
- **Dual Posting**: Use one command to post on both platforms simultaneously.
- **Daily Post Limits**: Keeps you from going overboard. Or not… it’s up to you.
- **Default Hashtags**: Automatically slap `#shipo-CLI` onto your posts so everyone knows how awesome you are.

---

## Installation

1. Clone the repo:
   ```bash
   git clone https://github.com/goldomatic/shipo-cli.git
   cd shipo-cli
   ```

2. Build the executable:
   ```bash
   go build -o shipo-cli
   ```

3. Optionally, move it to your `PATH`:
   ```bash
   mv shipo-cli /usr/local/bin/
   ```

---

## Configuration

Create a configuration file at `~/.config/shipo-cli/config`:

```plaintext
# Bluesky Credentials
handle = your-bluesky-handle.bsky.social
password = your-bluesky-password

# Daily Post Limit
limit = 5

# Twitter Credentials
twitter_consumer_key = YOUR_CONSUMER_KEY
twitter_consumer_secret = YOUR_CONSUMER_SECRET
twitter_access_token = YOUR_ACCESS_TOKEN
twitter_access_secret = YOUR_ACCESS_SECRET
```

---

## Usage

### Basic Commands

Post a message to **Bluesky**:
```bash
shipo-cli -p b -c "This is a Bluesky post!"
```

Post a message to **Twitter**:
```bash
shipo-cli -p t -c "This is a tweet!"
```

Post to **both platforms** at once:
```bash
shipo-cli -p bt -c "This is a cross-platform shitpost!"
```

### Command-Line Flags
| Flag  | Description                                      |
|-------|--------------------------------------------------|
| `-p`  | Platform: `b` for Bluesky, `t` for Twitter, `bt` for both |
| `-c`  | Content: The text you want to post               |

---

## Examples

1. **Posting to Bluesky**:
   ```bash
   shipo-cli -p b -c "Bluesky, here I come!"
   ```

2. **Posting to Twitter**:
   ```bash
   shipo-cli -p t -c "Hello, X! (but I’ll always call you Twitter)."
   ```

3. **Cross-Platform Posting**:
   ```bash
   shipo-cli -p bt -c "Shitposting on both platforms at once. Efficiency!"
   ```

---

## FAQ

### What happens if I hit my daily post limit?
shipo-CLI will gently (or not-so-gently) tell you to take a break and try again tomorrow. Or just set a higher limit. You’re the boss.

### Can I add hashtags automatically?
You bet! `#shipo-CLI` is added to every post because branding is important.

### Does shipo-CLI work on Windows?
Yes, but you’ll need to build the binary for Windows:
```bash
GOOS=windows GOARCH=amd64 go build -o shipo-cli.exe
```

---

## Contributing
Feel free to open an issue or create a pull request. Feature suggestions and bug reports are always welcome. Let’s make this the best (and silliest) CLI tool together!

---

## License
This project is licensed under the MIT License. Do whatever you want with it, but don’t forget to shitpost responsibly.

---

Happy shitposting! 🚀

