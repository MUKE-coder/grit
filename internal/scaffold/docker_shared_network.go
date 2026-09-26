package scaffold

import "fmt"

// dockerComposeSharedNetwork emits docker-compose.shared-network.yml.
//
// The second half of the address-pool fix. Dropping the named network from the
// production stack halves what it costs the daemon, but on a host whose pool is
// already empty even Compose's implicit <project>_default fails to allocate and
// the stack still will not start. Joining a network that already exists asks for
// no subnet at all, so there is nothing left to fail.
//
// An overlay rather than a second copy of the stack: the services, images and
// environment all come from docker-compose.prod.yml and only the networking is
// replaced, so the two cannot drift.
//
// Every service also answers to a project-prefixed alias, and the API is told to
// use those. A shared network is shared: two projects that both call a service
// "postgres" resolve each other's, and connecting to somebody else's database
// without noticing is a worse day than a deploy that fails loudly.
func dockerComposeSharedNetwork(opts Options) string {
	name := opts.ProjectName

	alias := func(svc string) string {
		return fmt.Sprintf("  %s:\n    networks:\n      default:\n        aliases: [%s-%s]\n", svc, name, svc)
	}

	services := ""
	for _, svc := range []string{"postgres", "pgbouncer", "redis", "minio"} {
		services += alias(svc)
	}
	if opts.ShouldIncludeWeb() {
		services += alias("web")
	}
	if opts.ShouldIncludeAdmin() {
		services += alias("admin")
	}
	if opts.ShouldIncludeDocs() {
		services += alias("docs")
	}

	// The API needs more than an alias: it has to be told to use the prefixed
	// names for everything it talks to.
	services += fmt.Sprintf(`  api:
    networks:
      default:
        aliases: [%s-api]
    environment:
      # The prefixed aliases, not the bare service names. On a shared network
      # "postgres" may very well be somebody else's.
      POSTGRES_HOST: %s-pgbouncer
      REDIS_URL: redis://:`+"${REDIS_PASSWORD:?set REDIS_PASSWORD in .env}"+`@%s-redis:6379
      MINIO_ENDPOINT: http://%s-minio:9000
`, name, name, name, name)

	return `# Join an existing network instead of creating one.
#
# Use this when a deploy fails with:
#
#   Network <project>_default  Error
#   all predefined address pools have been fully subnetted
#
# That is the Docker daemon, not this project. Docker slices bridge subnets out
# of 172.17.0.0/12 in /16 blocks, which is about sixteen networks for the whole
# daemon, and it does not release them when a deploy fails or a project is
# deleted. Once they are gone no new stack starts, however idle the host is.
#
# This is an overlay, not a replacement. Everything else comes from
# docker-compose.prod.yml, so the two cannot drift:
#
#   docker compose -f docker-compose.prod.yml \
#                  -f docker-compose.shared-network.yml up -d --build
#
# On Dokploy, set the Compose Path to both files, in that order.
#
# THE TRADE-OFF, AND IT IS A REAL ONE.
#
# A shared network is shared. Containers on it resolve each other by service
# name, so if another project on the host also calls a service "postgres",
# whichever answers first is what you get, and a silent connection to somebody
# else's database is worse than a failed deploy. Every service below therefore
# also answers to a project-prefixed alias, and the API is pointed at those
# rather than at the bare names.
#
# PREFER FIXING THE DAEMON. On the host:
#
#   docker network prune -f       # returns the leftovers of failed deploys
#
# and then, so it cannot recur, in /etc/docker/daemon.json:
#
#   { "default-address-pools": [ { "base": "172.17.0.0/12", "size": 24 } ] }
#
# followed by ` + "`systemctl restart docker`" + `, which gives about 4096 networks
# instead of 16. That restart bounces every container on the host, so pick a
# quiet moment. Once it is done, go back to docker-compose.prod.yml on its own
# and keep one private network per stack.

services:
` + services + `
networks:
  # Already exists, so Docker is never asked to allocate a subnet. Change the
  # name if your proxy uses a different one; ` + "`docker network ls`" + ` will say.
  default:
    external: true
    name: ` + "${SHARED_NETWORK:-dokploy-network}" + `
`
}
