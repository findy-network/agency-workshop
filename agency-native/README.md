# Running Agency Natively from Sources

## Start script

Open a terminal window to this folder and set up the agency to your localhost
using an installation script:

```shell
./setup.sh
```

If you don't want to install something like tmux and tmuxinator, just comment
the function calls away.

Otherwise, the script starts by cloning several 'findy-network' repos that are needed
agency services.

## Shutting down

`tmux kill-session`

