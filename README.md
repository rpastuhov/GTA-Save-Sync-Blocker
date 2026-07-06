
# GTA Save Sync Blocker

A Windows utility that toggles a firewall rule to block outbound traffic to a specific IP, controlled with global hotkeys. Built around the GTA Online save-blocking trick some players use during heists.

## How it works

Checks for admin rights on startup (needed to edit firewall rules) and exits if it doesn't have them.

Prints the instructions below, then checks whether Windows Firewall is on. If it's off, the program turns it on - otherwise the block rule does nothing - and restores it to whatever state it found on exit.

Installs a low-level keyboard hook and listens for two hotkeys:
- **Ctrl+F9** - adds an outbound block rule named `GTA_SYNC` targeting `192.81.241.171`.
- **Ctrl+F12** - removes that rule.

On exit (Ctrl+C or a termination signal): pulls the block rule if it's still active, restores the firewall to its original state.

## Usage instructions

```
========================================
              HOW TO USE
========================================

1. Start a heist or mission.

2. Press Ctrl+F9 at any point between loading into the job and the moment
   you receive the payout from Madrazo's people (or the mission finale).

3. After loading into the session, switch to Story Mode.

4. Wait for the cutscene to finish, then press Ctrl+F12.

5. Join a session via invite.

6. Deposit the money into your bank account.

7. Force a save (the easiest way is to change your outfit in the interaction
   menu - you can pick the same outfit and press Enter).

IMPORTANT:
- Cayo Perico limit: 4 heists per hour from the moment you receive your
  first payout (exceeding this limit is not recommended).
- Your cut is calculated as: potential take * 0.88 * your percentage.
- Contract limit: 2 per 30 minutes from the moment you receive your first
  payout (a third one won't pay out, and may result in a ban).

HOTKEYS:
- Ctrl+F9  - Enable save-blocking mode
- Ctrl+F12 - Disable save-blocking mode
```

## Building

```bash
go mod init gta_sync
go mod tidy
go build -o gta_sync.exe -ldflags "-s -w" .
```

Run `gta_sync.exe` as Administrator.

## Notes

That's outside Rockstar's ToS, and probably against the game's rules outright. Use it and you're taking on the risk yourself. For educational and research purposes only.

Apache License 2.0
