variable "TAG" {
  default = "latest"
}

group "default" {
  targets = ["meal-planning"]
}

target "meal-planning" {
  dockerfile = "Dockerfile"
  platforms = ["linux/amd64", "linux/arm64"]
  tags = [
    "ghcr.io/untanky/meal-planning:${TAG}"
  ]
  labels = {
    "org.opencontainers.image.version" = "latest"
  }
}
