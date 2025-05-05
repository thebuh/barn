<p align="center">
<b>barn</b> is a ASCOM Alpaca compatible safety monitor driver, allowing you to expose local files or remote URLs as ASCOM compatible safety monitor devices.
</p>

## Running

There are no compiled binaries provided (for now), but you can use docker to run barn.

Build image:
```shell
docker build . -t barn
```
And then run it, adjust ports and config file path if necessary:
```shell
docker run -p8888:8888 -p32227:32227 -it -v <Full path to config file>/barn.yaml:/barn.yaml barn
```

## Configuration

Configuration file should be named **barn.yaml** and goes into same directory as **barn** binary.
Below is an example configuration file to showcase the capabilities of **barn**.

```yaml
api: # (optional)
  port: 8080 # Api port. 8080 by default
discovery: # (optional)
  port: 32227 #Alpaca discovery port. 32227 by default
monitors: # (mandatory)
  http: # Define HTTP safety monitors
    remote: # Monitor ID should be unique
      name: "Some remote url" # Safety monitor name
      description: "Some remote url description" # Description
      url: http://127.0.0.1/test # Url to check
    remote2:
      name: "Second remote url"
      description: "Second remote url description"
      url: http://127.0.0.2/test
  file: # Define local file safety monitor
    local: 
      name: "Local file" 
      description: "Some local file description"
      path: /tmp/test # File path to check
  dummy: # Fake safety monitors (Stays always in defined state) 
    fake:
      is_safe: true # State of the monitor
```
### Rules

By default, barn matches **1** or **true** (case-insensitive) as "safe". You can customize that behaviour with following configuration

```yaml
monitors:
  http:
    remote:
      name: "Some remote url"
      description: "Some remote url description"
      url: http://127.0.0.1/test
      rule:
        invert: false # Invert matching result (if pattern key exists applied after match) 
        pattern: "(open|opening)" # Regular expression to match
```

## Todo
- JSON support.
