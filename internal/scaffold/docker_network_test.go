package scaffold

import (
	"strings"
	"testing"
)

// A stack must not ask the daemon for a subnet it does not need.
//
// Docker hands out bridge subnets from 172.17.0.0/12 sliced into /16 blocks,
// which is about sixteen networks for the entire daemon, and it does not
// release them when a deploy fails or a project is deleted. A host running a
// handful of Compose stacks runs out, and every further deploy dies at the last
// step with "all predefined address pools have been fully subnetted" after
// every image has already built. The message says nothing about disk, memory or
// CPU because none of them are the problem.
//
// The production stack used to declare a named network and attach every service
// to it. Compose's implicit <project>_default gives the same service-name DNS
// and the same isolation, for one fewer thing asking the daemon for a subnet.
func TestProductionStackDeclaresNoNetwork(t *testing.T) {
	for _, arch := range []Architecture{ArchTriple, ArchDouble, ArchAPI} {
		opts := Options{ProjectName: "shop", Architecture: arch, Frontend: FrontendNext}
		yaml := dockerComposeProd(opts)

		for _, line := range strings.Split(yaml, "\n") {
			if strings.HasPrefix(line, "networks:") {
				t.Errorf("%s: the production stack declares a network, which costs a second subnet", arch)
			}
			if strings.HasPrefix(strings.TrimSpace(line), "container_name:") {
				t.Errorf("%s: %q pins a container name; names are global to the daemon, so a second copy of the stack cannot start",
					arch, strings.TrimSpace(line))
			}
		}
		// And no service is attached to one either, which is what would pull
		// the declaration back in.
		if strings.Contains(yaml, "    networks:\n      - shop") {
			t.Errorf("%s: a service is still attached to a named network", arch)
		}
		// The comment has to explain it, or the next person adds it back.
		if !strings.Contains(yaml, "predefined address pools") &&
			!strings.Contains(yaml, "sixteen networks") {
			t.Errorf("%s: nothing says why there is no network here", arch)
		}
	}
}

// And a way out for a host whose pool is already empty.
//
// Dropping the named network halves what a stack costs, but on a host with no
// subnets left even <project>_default fails to allocate and the stack still
// will not start. Joining a network that already exists asks for nothing.
func TestSharedNetworkOverlayCreatesNothing(t *testing.T) {
	opts := Options{ProjectName: "shop", Architecture: ArchTriple, Frontend: FrontendNext}
	yaml := dockerComposeSharedNetwork(opts)

	if !strings.Contains(yaml, "external: true") {
		t.Fatal("the overlay does not join an existing network, so it allocates a subnet like any other")
	}
	if !strings.Contains(yaml, "${SHARED_NETWORK:-dokploy-network}") {
		t.Error("the network name cannot be overridden, so this only works on Dokploy")
	}

	// A shared network is shared: without prefixed aliases, two projects that
	// both call a service "postgres" resolve each other's, and connecting to
	// somebody else's database unnoticed is worse than a deploy that fails.
	for _, svc := range []string{"postgres", "pgbouncer", "redis", "minio", "api", "web", "admin"} {
		if !strings.Contains(yaml, "aliases: [shop-"+svc+"]") {
			t.Errorf("%s has no project-prefixed alias on the shared network", svc)
		}
	}
	// And the API is pointed at those aliases rather than the bare names.
	for _, want := range []string{
		"POSTGRES_HOST: shop-pgbouncer",
		"@shop-redis:6379",
		"MINIO_ENDPOINT: http://shop-minio:9000",
	} {
		if !strings.Contains(yaml, want) {
			t.Errorf("the API still resolves a bare service name: missing %q", want)
		}
	}

	// The fix that actually belongs on the host is named, because the overlay
	// is a stopgap and saying so is the difference between a workaround and a
	// new default nobody revisits.
	for _, want := range []string{"docker network prune -f", "default-address-pools"} {
		if !strings.Contains(yaml, want) {
			t.Errorf("the overlay does not mention the host fix: %q", want)
		}
	}
}

// An architecture with no web or admin app must not alias services it does not
// have: Compose rejects an override for a service that is not in the base file.
func TestSharedNetworkOverlayMatchesTheArchitecture(t *testing.T) {
	yaml := dockerComposeSharedNetwork(Options{
		ProjectName: "shop", Architecture: ArchAPI, Frontend: FrontendNext,
	})
	for _, absent := range []string{"  web:", "  admin:", "  docs:"} {
		if strings.Contains(yaml, absent) {
			t.Errorf("an API-only project's overlay names %q, which the base file does not define", strings.TrimSpace(absent))
		}
	}
	if !strings.Contains(yaml, "  api:") {
		t.Error("the overlay does not cover the API")
	}
}
