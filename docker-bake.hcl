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
    "org.opencontainers.image.source" = "https://github.com/Untanky/meal-planning"
    "org.opencontainers.image.version" = "latest"
    "org.opencontainers.image.title" = "Meal Planner"
    "org.opencontainers.image.description" = "Self-contained, web-app to manage meal planning"
    "org.opencontainers.image.authors" = "Lukas Grimm"
  }
}
