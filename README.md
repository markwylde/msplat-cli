![msplat logo](https://raw.githubusercontent.com/msplat/branding/master/msplat-logo-nameonly-sm.png)

[![Build Status](https://travis-ci.org/msplat/msplat-cli.svg?branch=master)](https://travis-ci.org/msplat/msplat-cli)

The msplat cli is a command line interface for managing a microservices platform using the msplat concept. This project is very new and hasn't even reached it's alpha version.

# Tutorials
The wiki has some tutorials and getting started guides for creating your first msplat environment.

- [Create a new Development Environment](https://github.com/msplat/msplat-cli/wiki/Create-a-new-development-environment)
- [Create a new Production Environment](https://github.com/msplat/msplat-cli/wiki/Create-a-production-environment)
- [...more tutorials](https://github.com/msplat/msplat-cli/wiki)

# Build
You will need go version 1.11.4 or higher, as this project uses go modules.

To build you can execute the `build.sh` file in the root of the project:
```bash
## To build ignoring autocomplete
./build.sh

## To build and activate autocomplete in the active shell
. ./build.sh
```

This will produce a binary at `dist/mtk` which can be installed globally:
```bash
sudo ln -s `pwd`/dist/mtk /usr/local/bin/mtk
```

## Autocomplete
The CLI comes with autocomplete functionality. To activate it you will need to run the following for every shell. Alternatively you can add it to your `.bashrc` file.

```
cd dist && PROG=mtk source ../autocomplete/bash_autocomplete
```

# Usage
Before using the cli you will need to define the config file for the environment.

You can do this by specifying the `MSPLAT_CONFIG` environment variable or providing a `--config` flag.

```bash
$ export MSPLAT_CONFIG=~/Documents/example/config.yml
$ mtk --help

NAME:
   msplat toolkit - A cli for managing msplat environments
USAGE:
   mtk [global options] command [command options] [arguments...]

COMMANDS:
   projects, pr         Tasks for managing projects
   services, sv         Tasks for managing services
   stacks, st           Tasks for managing stacks
   help, h              Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --config value  configuration settings for msplat (default: "?cwd") [$MSPLAT_CONFIG]
   --help, -h      show help
   --version, -v   print the version

VERSION:
   0.0.1
```

# License
This project is licensed under the terms of the GPLv3 license.
