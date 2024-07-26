# Cowatch

**Run multiple commands/watchers concurrently.**

This primarily targets development environments that rely on compilers but lack
hot module replacement. For example, if you use Go with Templ and Tailwind, the
latter might take too long rebuilding the CSS. So, `cowatch` helps you run
Tailwind in watch mode concurrently with Go and Templ which you could watch
with [Air](https://github.com/air-verse/air).

## Usage

- Install
```bash
go install github.com/waseem-medhat/cowatch@latest
```

- Create a `cowatch.toml` config file in your project root. You can find an
example
[here](https://github.com/waseem-medhat/cowatch/blob/main/cowatch.toml).

- Run `cowatch` to start your commands. You should start seeing the commands'
standard output/error in your console.
