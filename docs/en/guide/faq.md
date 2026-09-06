# FAQ

## General

### What is AirlineSim Autobuy?

AirlineSim Autobuy is a tool that monitors the AirlineSim used aircraft market and automatically purchases aircraft based on rules you define. It scans game servers at configurable intervals, evaluates listings against your rules, and optionally purchases matching aircraft.

### Is this tool allowed by the game's terms of service?

AirlineSim Autobuy is designed for personal use and educational purposes. It interacts with the game through the same HTTP interface as a regular browser. However, automated tools may violate the game's terms of service. Use at your own risk. The developers recommend using reasonable polling intervals to avoid placing excessive load on game servers.

### Does the tool work on all AirlineSim servers?

Yes, the tool supports any AirlineSim game server. Configure the server URL in the `servers` section of the configuration file. It works with both free and premium servers.

### Do I need programming experience to use this?

No. Basic configuration is done through a YAML file, and the web interface allows you to manage rules and settings without editing files directly. However, familiarity with editing text files and using a terminal is helpful.

### What is the web UI built with?

The web UI is a Vue 3 single-page application built with Vite. It is embedded in the Go binary, so there are no additional dependencies to run the web interface.

## Configuration

### How do I set up multiple accounts?

Use the `auths` list in the configuration file:

```yaml
auths:
  - username: "account1@example.com"
    password: "password1"
    session_file: "session1.json"
  - username: "account2@example.com"
    password: "password2"
    session_file: "session2.json"
```

Then reference the auth index in your rules using `auth_id`. See the [Configuration Guide](/guide/configuration#auths) for details.

### Can I monitor multiple servers at once?

Yes. Add multiple entries to the `servers` list:

```yaml
servers:
  - host: free1
    base_url: https://free1.airlinesim.aero
  - host: free2
    base_url: https://free2.airlinesim.aero
```

The engine starts a separate monitoring loop for each server.

### How do I change the polling interval?

Edit the `monitor.interval` value in the configuration file. The value is in seconds. For example, `interval: 60` scans the market once per minute.

### What does the jitter setting do?

Jitter adds a random delay to each polling interval. This makes the request pattern less predictable, which helps avoid rate limiting. With `interval: 30` and `jitter: 10`, the actual delay is between 30 and 40 seconds.

### How do I disable the web UI?

Set `webui.enabled: false` in the configuration file:

```yaml
webui:
  enabled: false
```

### How do I make the web UI accessible only from localhost?

Set the bind address to `127.0.0.1`:

```yaml
webui:
  host: 127.0.0.1
```

### Can I use environment variables for configuration?

Not currently. All configuration is read from the YAML file. This feature may be added in a future release.

## Technical Issues

### Authentication fails with "login failed"

Check the following:

1. **Username and password** -- Verify they are correct in the configuration file
2. **Server URL** -- Ensure the `base_url` is correct (e.g., `https://free1.airlinesim.aero`)
3. **Network connectivity** -- Ensure the machine running the tool can reach `sar.simulogics.games` and the game server
4. **Account status** -- Ensure your AirlineSim account is active and not suspended
5. **Session file** -- Delete the session file (default: `session.json`) and try again

### The engine starts but finds no aircraft

1. **Check the market** -- Log into the game manually and verify there are aircraft listed on the market
2. **Verify rules** -- Ensure your rules are enabled and have correct match conditions
3. **Check logs** -- Look for warnings or errors in the console output
4. **Session health** -- The session may have expired. Check for "session expired, re-authenticating" messages

### The purchase fails with "bid failed"

1. **Insufficient balance** -- Check your account balance in the game
2. **Auction expired** -- The auction may have ended before your bid was placed
3. **Already outbid** -- Another player may have placed a higher bid
4. **Session expired** -- The session may have expired between scanning and purchasing
5. **Check the aircraft URL** -- The aircraft may have been removed from the market

### The web UI is not loading

1. **Check if the engine is running** -- The web UI is embedded in the binary; if the application is running, the web UI should be available
2. **Verify the port** -- Ensure the configured port is not in use by another application
3. **Check the bind address** -- If bound to `127.0.0.1`, you can only access it from the same machine
4. **Firewall** -- Ensure the port is not blocked by a firewall
5. **Check the logs** -- Look for web UI startup errors

### The configuration file is not being reloaded

1. **File watcher** -- The configuration file is polled every 5 seconds for changes
2. **File format** -- Ensure the YAML file is valid (check for indentation errors)
3. **Section limitations** -- Only rule changes are hot-reloaded. Changes to servers, auths, monitor, and webui sections require a restart

### I get "too many requests" errors from the game server

1. **Increase the interval** -- Set `monitor.interval` to a higher value (e.g., 60 seconds)
2. **Increase jitter** -- Add more jitter to make request timing less predictable
3. **Reduce servers** -- Monitor fewer servers simultaneously

### The binary is flagged by antivirus software

The Go binary is compiled from source code. Some antivirus software may flag Go binaries due to their packaging. You can:

1. Build from source yourself using `make build`
2. Add an exclusion for the binary in your antivirus software
3. Verify the binary SHA256 checksum against a trusted build

## Troubleshooting

### Debug Mode

Run the application with debug-level logging to get more detailed output:

```bash
# Temporarily edit main.go to change log level:
slog.LevelDebug

# Or build a debug version:
go build -gcflags="-N -l" -o autobuy-debug ./cmd/autobuy
```

### Session Issues

If you encounter persistent session problems:

1. Stop the application
2. Delete the session file (`session.json` by default)
3. Restart the application
4. A fresh login will be performed

### Checking the Market URL

The tool discovers the market URL automatically. To verify the discovery is working:

1. Enable debug logging
2. Look for the line: `INFO discovered market URL url=...`
3. If discovery fails, the tool falls back to a direct URL

### Validating the Configuration

Use a YAML validator to check your configuration file for syntax errors:

```bash
# Using Python
python -c "import yaml; yaml.safe_load(open('configs/config.yaml'))"

# Using yq (if installed)
yq eval configs/config.yaml > /dev/null
```

### Common Error Messages

| Error | Likely Cause | Solution |
|---|---|---|
| `login failed` | Invalid credentials or network issue | Verify username/password and server URL |
| `failed to fetch market page` | Network issue or expired session | Check connectivity; session will auto-renew |
| `failed to parse HTML` | Market page format changed | Report as a bug on GitHub |
| `purchase returned status 404` | Aircraft no longer available | The listing was removed or sold |
| `bid failed: unexpected response` | Auction format may have changed | Check the aircraft page manually |
| `address ... is not allowed` | URL validation blocked the request | Check server URL configuration |

### Getting Help

If you encounter issues not covered here:

1. **Check the logs** -- Look for error messages in the console output
2. **Search GitHub Issues** -- Check if the problem has been reported at [github.com/Maicarons/airlinesim-autobuy/issues](https://github.com/Maicarons/airlinesim-autobuy/issues)
3. **Open a new issue** -- Provide the following information:
   - Application version (commit hash or binary build date)
   - Configuration file (with sensitive fields redacted)
   - Log output from the session
   - Steps to reproduce the issue