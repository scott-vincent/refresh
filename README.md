# Refresh

## Rebuild and re-run your Go server when back-end files change
## + Rebuild your Vue application when front-end files change

This project was inspired by [https://github.com/pilu/fresh](https://github.com/pilu/fresh). The lack of updates and response from the maintainer, non-idiomatic codebase, numerous bugs, and lack of detailed reporting made the project a dead end for me to use. Enter `refresh`.

This simple command line application will watch your files, trigger a build of your Go binary and restart the application for you.

Forked Version: This version of refresh which was forked from [https://github.com/markbates/refresh](https://github.com/markbates/refresh) adds an extra YAML setting *go_or_vue* so you can rebuild your *vue* front end when your *.vue* files change. If you do this, you should run a second version of refresh with *go_or_vue* set to *go* to rebuild (and rerun) the backend server.

## Installation

```
$ go get github.com/markbates/refresh
```

## Getting Started

First you'll want to create a `refresh.yml` configuration file:

```
$ refresh init
```

If you want the config file in a different directory:

```
$ refresh init -c path/to/config.yml
```

Set it up the way you want, but I believe the defaults really speak for themselves, and will probably work for 90% of the use cases out there.

## Usage

Once you have your configuration all set up, all you need to do is run it:

```
$ refresh run
```

That's it! Now, as you change your code the binary will be re-built and re-started for you.

## HTTP Handler

Refresh is nice enough to ship with an `http.Handler` that you can wrap around your requests. Why would you want to do that?
Well, if there is an error doing a build, the built in `http.Handler` will print the error in your browser in giant text so you'll know that there was a problem, and where to fix it (hopefully).

```go
...
m := http.NewServeMux()
err = http.ListenAndServe(":3000", web.ErrorChecker(m))
...
```

## Configuration Settings

```yml
# The root of your application relative to your configuration file.
app_root: .
# List of folders you don't want to watch. The more folders you ignore, the 
# faster things will be.
ignored_folders:
  - vendor
  - log
  - tmp
# List of file extensions you want to watch for changes.
included_extensions:
  - .go
# The directory you want to build your binary in.
build_path: /tmp
# `fsnotify` can trigger many events at once when you change a file. To minimize
# unnecessary builds, a delay is used to ignore extra events.
build_delay: 200ms
# If you have a specific sub-directory of your project you want to build.
build_target_path : "./cmd/cli"
# Build vue instead of go (binary/command settings will be ignored as there is nothing to run).
go_or_vue: "vue"
# What you would like to name the built binary.
binary_name: refresh-build
# Extra command line flags you want passed to the built binary when running it.
command_flags: ["--env", "development"]
# Extra environment variables you want defined when the built binary is run.
command_env: ["PORT=1234"]
# If you want colors to be used when printing out log messages.
enable_colors: true
```
