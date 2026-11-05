package scaffold

import "strings"

// The migrate service, in one place, because two have to agree about it: the
// production compose file a new project gets, and the repair that puts it into
// a project that already exists.
//
// A first deploy on Postgres used to come up with no tables. cmd/server
// migrates only when DB_PROVIDER=memory, where the schema cannot outlive the
// process; the API image carried only the server binary; and the production
// compose file had no migrate step. So the API started, answered "relation
// does not exist" to every request, and there was nothing in the image to run
// by hand (grit#177).

// composeMigrateService is the service itself, written between "services:" and
// the api service.
//
// The environment block is the api's own, minus Redis and MinIO, which a
// migration does not touch.
const composeMigrateService = `  # The schema, before anything serves a request.
  #
  # cmd/server migrates only for DB_PROVIDER=memory, where the schema cannot
  # outlive the process. On a real database this is the step that creates it,
  # and without it a first deploy came up with no tables and answered
  # "relation does not exist" to every request (grit#177).
  #
  # One-shot: it runs, it exits, and the api waits for that exit rather than
  # for a health check. A failed migration therefore stops the api from
  # starting, which is what you want. restart: "no" because a migration that
  # failed should be read, not retried in a loop.
  migrate:
    build:
      context: ./apps/api
      dockerfile: Dockerfile
    command: ["./migrate"]
    restart: "no"
    env_file:
      - .env
    environment:
`

// composeMigrateDependency is what the api's depends_on gains.
//
// service_completed_successfully, not service_healthy: this container exits,
// and waiting for it to become healthy would wait for ever.
const composeMigrateDependency = "      migrate:\n        condition: service_completed_successfully\n"

// composeAPIRedisDepends is where the dependency is inserted, so the order
// matches the file a new project gets.
const composeAPIRedisDepends = "      redis:\n        condition: service_healthy\n"

// repairComposeMigrateSource adds the migrate service to a production compose
// file that predates it.
//
// Both halves or neither: a service nothing waits for would run beside the API
// rather than before it, which is the race this exists to remove, and a
// depends_on naming a service that is not there stops the whole stack from
// starting.
func repairComposeMigrateSource(src string) (string, []string, []string) {
	// Already there, or not a file this understands.
	if strings.Contains(src, "\n  migrate:\n") {
		return src, nil, nil
	}
	start, end, ok := serviceBlock(src, "api")
	if !ok {
		return src, nil, nil
	}
	// The api must be the Grit one, waiting on redis in the shape the template
	// writes. Anything else is somebody's own compose file, and inserting into
	// it would be a guess.
	block := src[start:end]
	if !strings.Contains(block, composeAPIRedisDepends) {
		return src, nil, []string{
			"the api service in docker-compose.prod.yml is not the one Grit wrote, so the " +
				"migrate service was not added: without it a first deploy on Postgres starts " +
				"with no tables (grit#177)",
		}
	}

	// The environment the migrate service needs, taken from the api's own so
	// the two cannot disagree about which database they point at.
	env := composeDBProvider + composeAPIPostgresNew
	service := composeMigrateService + env + composeAPIDependsNew + "\n"

	anchor := "services:\n"
	if !strings.Contains(src, anchor) {
		return src, nil, nil
	}
	out := strings.Replace(src, anchor, anchor+service, 1)

	// And the api waits for it. Recomputed because the insert above moved
	// everything after it.
	start, end, ok = serviceBlock(out, "api")
	if !ok {
		return src, nil, nil
	}
	block = out[start:end]
	block = strings.Replace(block, composeAPIRedisDepends,
		composeMigrateDependency+composeAPIRedisDepends, 1)
	out = out[:start] + block + out[end:]

	return out, []string{
		"a migrate service runs before the api, so a first deploy on Postgres has tables",
	}, nil
}
