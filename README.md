The msplat cli is a command line interface for managing a microservices platform using the msplat concept. This project is very new and hasn't even reached it's alpha version.

# Build
You will need go version 1.11.4 or higher, as this project uses go modules.

To build you can execute the `build.sh` file in the root of the project:
```bash
./build.sh
```

This will produce a binary at `dist/mtk` which can be installed globally:
```bash
sudo ln -s `pwd`/dist/mtk /usr/local/bin/mtk
```

# Usage
```bash
$ mtk --help

NAME:
   mtk - A new cli application
USAGE:
   mtk [global options] command [command options] [arguments...]

COMMANDS:
   project         tasks for managing projects
   service         tasks for managing services
   stack           tasks for managing stacks
   help, h         Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --config value  configuration settings for msplat (default: "./config.yml") [$MSPLAT_CONFIG]
   --help, -h      show help
   --version, -v   print the version

VERSION:
   0.0.1
```

# License
This project is licensed under the terms of the GPLv3 license.