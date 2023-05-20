![release](https://github.com/tomkcey/valheim_plus_config_override/actions/workflows/release.yml/badge.svg)
![tag](https://github.com/tomkcey/valheim_plus_config_override/actions/workflows/tag.yml/badge.svg)

### Why?

For those using the Valheim Plus mod on the game Valheim, sometimes you have a couple configurations you'd like to change but you might have to dig through the file to find each key/value pair. How tiring is that?

### How?

Just execute the program, enter the absolute path to source file, then to target file (with which you'll override the source file), and it will create a `*.cfg` file in your temp directory, which may be different depending on your operating system. The file's name is a time (Unix) number.

### How to compile?

If you have Go installed, position yourself at root and then enter `go build main.go` in your terminal. This will generate a `main.exe` file which you can run from the terminal.

### How to use with docker?

Position yourself at root, then run `docker build . -t <name>`, then `docker run -it <name>`.
