# Create Service

## Group

> Group is Required Before Creating a Service​​

A group is a mechanism for dividing cluster resources into multiple virtual clusters, similar to `Kubernetes namespaces.`

After creating a group, you can add `nodes` to it.

You can also configure group permissions by adding `users` or `teams` with management privileges.

Once the group is set up, you can proceed to create a service.

##  Service

> Service: The Core Object for Managing Containers​​

A Service is the central entity used to manage containers. It consists of the following components:

- ​​Container Configuration​​: Defines the container's options and specifications.
- ReplicaSet: The number of container instances to create.
- Deployment Strategy​​: Specifies how containers are deployed across the cluster.


![](/_media/service-list.png)

## Service Status Definition

The service has the following three states:

- `Assigning`​​: The service's current state hasn't reached the desired state, and Humpback is currently scheduling resources.
- `Running`​​: The service has successfully achieved its desired state and is running normally.
- `Warning`​​: The service is experiencing anomalies and cannot reach its desired state, requiring operating.

## Setup Container Options

In the service application page, you can setup container options like below:

![](/_media/service-application.png)