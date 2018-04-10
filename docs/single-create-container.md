# Create Container

![Create Container](_media/single-create-container.png)

> Description

- `Name`: Container name
- `Image`:  The name of the image used to create the container, such as docker.io/helloworld:latest
- `Command`: The command running container, and if empty, use the predefined commands in Dockerfile
- `Network Mode`: Network mode, `Host`(Direct use of the host network), `Bridge`(Through the Docker network card to the host network), `Customer`(To use a custom network, you need to specify `Network Name`)
  - `Network Name`: Set when `Network Mode` is` Customer`
  - `HostName`: Set when `Network Mode` is not` Host`
- `Restart Policy`: The restart policy when the container exits. `always` means the container will always restart after exiting (except for the user to manually stop). `on-failure` means only when the container does not normally exit will restart (exit code is not 0). Default is `none` which means does not restart.
  - `Max Retry Count`: When `Restart Policy` is `on-failure`, you need to specify the maximum number of retries
- `DNS`: The container needs to use the DNS server
- `CPU Shares`: Specifies the CPU usage weight for the current container relative to other containers
- `Memory Limit`: Maximum available memory
- `Ports Binging`: When `Network Mode` is not `Host`, you can set port mapping
- `Volumes Binding`: Container volume mapping
- `Environment Variables`: Set container environment variables, such as "VAR=value"
- `Links`:  A link between containers, such as container_name:alias
- `Log Config`: Log configuration, there are two properties Driver and Opts. The default value is Driver:json-file,Opts:{max-size = 10m,max-file = 3}
