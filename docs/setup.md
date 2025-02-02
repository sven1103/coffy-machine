# Setup

## Configuration file

Coffy Machine expects a configuration file on startup. You can take the provided example one ``example_config.yaml`` for starters and then start to adjust it to your needs.

### Server

Currently, you can only change the port that Coffy Machine listens to incoming HTTP requests (default: `8080`):

```yaml
server:
  port: 8080 # Change to your requirements
```

### Database

To ensure Coffy Machine's simplicity to setup and run, currently only SQLite is supported and creates a file
on the host's filesystem.

The default will create a database file in the current working directory.

Of course, you can set any path to the host's filesystem. Ensure Coffy Machine can 
access the target directory.

```yaml
database:
  path: ./coffy_machine.db # The default, change it if needed
```
