# Create Container

![Create Container](_media/single-create-container.png)

> Description

- `Name`: Container name
- `Image`:  The name of the image used to create the container, such as docker.io/helloworld:latest
- `Command`: The command to use when running the container, and if it is empty, use the predefined commands in Dockerfile
- `Network Mode`: Network mode, `Host`(Direct use of the host network), `Bridge`(Through the Docker network card to the host network), `Customer`(To use a custom network, you need to specify `Network Name`)
  - `Network Name`: When `Network Mode` is` Customer` can be set
  - `HostName`: When `Network Mode` is not` Host` 'can be set
- `Restart Policy`: When the container exits the restart policy, `always` means In any case, the container exits (except for the user to manually stop) will restart, `on-failure` means only when the container does not normally exit will restart (exit code is not 0), the default is `none` does not restart
  - `Max Retry Count`: When `Restart Policy` is `on-failure`, you need to specify the maximum number of retries
- `DNS`: The container needs to use the DNS server
- `CPU Shares`: Specifies the CPU usage weight for the current container relative to other containers
- `Memory Limit`: Maximum available memory
- `Ports Binging`: When `Network Mode` is not `Host`, you can set port mapping
- `Volumes Binding`: Container volume mapping
- `Environment Variables`: Set container environment variables, such as "VAR=value"
- `Links`:  A link between a container and a container, such as container_name:alias
