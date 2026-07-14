package supervisor

import (
	"fmt"

	"github.com/pafthang/pocketagent/internal/ctrl/catalog"
)

func topoSort(services []catalog.Service) ([]string, error) {
	names := make(map[string]int)
	inDegree := make(map[string]int)
	graph := make(map[string][]string)

	for _, svc := range services {
		names[svc.Name]++
		if names[svc.Name] > 1 {
			return nil, fmt.Errorf("duplicate service %q", svc.Name)
		}
		inDegree[svc.Name] = 0
	}

	for _, svc := range services {
		for _, dep := range svc.DependsOn {
			if _, ok := names[dep]; !ok {
				return nil, fmt.Errorf("service %q depends on unknown %q", svc.Name, dep)
			}
			graph[dep] = append(graph[dep], svc.Name)
			inDegree[svc.Name]++
		}
	}

	queue := make([]string, 0)
	for name := range inDegree {
		if inDegree[name] == 0 {
			queue = append(queue, name)
		}
	}

	order := make([]string, 0, len(services))
	for len(queue) > 0 {
		n := queue[0]
		queue = queue[1:]
		order = append(order, n)
		for _, next := range graph[n] {
			inDegree[next]--
			if inDegree[next] == 0 {
				queue = append(queue, next)
			}
		}
	}

	if len(order) != len(services) {
		return nil, fmt.Errorf("circular service dependency detected")
	}

	return order, nil
}