# Configuration Management

Humpback provides configuration management. You can store non-sensitive information in configs, allowing your images to remain as generic as possible.  

Humpback supports two types of configs:  
- Static: Key-value pairs. You can use these configs in environment variables.  
- Volume: Configuration content will be saved as files and mounted into containers via volumes.  

The screenshot below demonstrates the two different config types.  

![](/_media/manage-config.png)

You can then use them in your containers.  

![](/_media/use-config-volume.png)