<p align="center">
<b>barn</b> is a ASCOM Alpaca compatible safety monitor driver, allowing you to expose local files or remote URLs as ASCOM compatible safety monitor devices.
</p>

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

## Todo
- Regular expressions and JSON support. For now **barn** only looks for "true" / "1" in files and urls. 
