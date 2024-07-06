target "docker-metadata-action" {}

target "meal-planner" {
  inherits = ["docker-metadata-action"]
  context    = "./"
  dockerfile = "Dockerfile"
  platforms = [
    "linux/amd64",
    "linux/arm64",
  ]
  labels = {
    "org.opencontainers.image.version" = "latest"
  }
}
