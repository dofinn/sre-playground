promgraph := o11y-promgraph
run-promgraph:
	podman play kube promgraph/pod_promgraph.yaml --configmap promgraph/cm_grafana.yaml --configmap promgraph/cm_prometheus.yaml

stop-promgraph:
	podman pod stop ${promgraph}

clean-promgraph:
	podman pod rm ${promgraph}