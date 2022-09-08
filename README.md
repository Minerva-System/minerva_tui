# Minerva TUI

A CLI interface for the [Minerva System](https://github.com/luksamuk/minerva-system).

Still a work-in-progress.

Check out for binary downloads on the Releases page.

## Running with Docker

To run this interactively through console using Docker, just...

```bash
docker run -it luksamuk/minerva_tui
```

Or, to create a container named `minerva_tui` and run it:

```bash
./run.sh
```

It will work better if your current terminal has color capabilities
(specifically, it expects that your terminal supports an environment
variable `TERM` as `xterm-256color`).

You don't have to build the container, it is already hosted on DockerHub.

## License

This software is distributed under the GPLv3 license. See LICENSE for details.


