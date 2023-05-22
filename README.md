![tag](https://github.com/tomkcey/valheim_plus_config_override/actions/workflows/tag.yml/badge.svg)

## TODO

- [ ] Sort produced config sections
- [ ] Produce tests
- [ ] Improve error handling
- [ ] Check that input file is a \*.cfg

### Why?

For those using the Valheim Plus mod on the game Valheim, sometimes you have a couple configurations you'd like to change but you might have to dig through the file to find each key/value pair. How tiring is that? Very. Here comes the Valheim Plus Config Overrider.

### How?

Just execute the program, enter the absolute path to source file, then to target file (with which you'll override the source file), and it will create a `*.cfg` file in your temp directory, which may be different depending on your operating system. The file's name is a time (Unix) number.
