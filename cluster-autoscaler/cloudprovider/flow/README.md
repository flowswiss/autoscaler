# Cluster Autoscaler for Flow

The cluster autoscaler for Flow Cloud scales worker nodes within any specified Flow Kubernetes cluster.

# Configuration
As there is no concept of a node group within Flow's Kubernetes offering, the configuration required is quite 
simple. You need to set:

- The ID of the cluster
- The Flow access token literally defined
- The Flow URL (optional; defaults to `https://api.flow.ch/`)
- The minimum and maximum number of **worker** nodes you want (the master is excluded)

See the [cluster-autoscaler-standard.yaml](examples/cluster-autoscaler-standard.yaml) example configuration, but to 
summarise you should set a `nodes` startup parameter for cluster autoscaler to specify a node group called `workers` 
e.g. `--nodes=3:10:workers`.

The remaining parameters can be set via environment variables (`FLOW_CLUSTER_ID`, `FLOW_API_TOKEN` and `FLOW_API_URL`) as in the
example YAML.

It is also possible to get these parameters through a YAML file mounted into the container
(for example via a Kubernetes Secret). The path configured with a startup parameter e.g.
`--cloud-config=/etc/kubernetes/cloud.config`. In this case the YAML keys are `api_url`, `api_token` and `cluster_id`.


# Development

Make sure you're inside the root path of the [autoscaler
repository](https://github.com/kubernetes/autoscaler)

1.) Build the `cluster-autoscaler` binary:


```
make build-in-docker
```

2.) Build the docker image:

```
docker build -t cloudbit/cluster-autoscaler:dev .
```

3.) Push the docker image to Docker hub:

```
docker push cloudbit/cluster-autoscaler:dev
```