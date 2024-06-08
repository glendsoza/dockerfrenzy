# DOCKER FRENZY

## An Easy and Convenient Way to Interact with Docker Daemon and Hosts

### Features

- List all hosts and SSH into them.
- List containers and images.
- View logs, execute commands, and perform actions on containers.
- View image details and create containers from them.

### Installing

#### Using Docker Compose

1. Go to the root of the repository.
2. Run `docker-compose up`.

#### From Source

Building from source requires Go 1.22 and Node v16.20.2.

Refer to the Docker and docker-compose files for more instructions on how to build from source.

Note: Host on which server is running requires `ssh` and `sshpass` application installed

### Config File and Environment Variables

The config file is in the following format:

```yaml
configList:
  - sshConfig:
      passwordAuth:
        username: <>
        password: <>
      ips:
        - <>
  - sshConfig:
      sshAuth:
        username: 
        PrivateKeyFile: 

```

Refer to the Docker Compose file and .env file for more information about environment variables.

`PrivateKeyFile` should be present in the location pointed by `CONFIG_FOLDER` environment variable

Note
This project is still in beta and should not be used in production.

Roadmap

- Add authentication
- Add unit tests 
- Store credentials in config file more securely 
- Better notification and alerting system
