# Service Scheduling Strategy

**Humpback** offers highly flexible deployment strategies to meet various business scenarios.  
As shown in the figure below, you can configure deployment strategies on the Service's Deployment page.  

![](/_media/humpback-schduler.png)

The deployment strategy consists of three configurations:  
- `Dispatch Mode`: Controls the number of service replicas in a group.  
   - Selecting **Global** will create one container on each eligible node.  
   - Selecting **Replicated** and set replica count. Humpback will create the specified number of containers based on resource usage across nodes in the group.  
- `Placement Constraints`: Controls node distribution for containers by setting constraints via IP or labels.  
- `Schedules Info`: Configures containers to run in scheduled task mode.  
   • Humpback supports **Cron expressions** for defining scheduled tasks.  
   • It also allows **manual startup**.  
   • **Timeouts**: Humpback will automatically stop containers when the configured timeout is reached.  
