package e2e

import "flag"

var NoDocker = flag.Bool("no-docker", false, "Wether the tests are running inside an environment with docker available")
