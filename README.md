# Meal Planning

This app is a meal planner to help you track the meals for the week and manage what you eat.

The app is a self-contained Docker image and can be easily be pulled and deployed.

```
docker run -p 8080:8080 ghcr.io/untanky/meal-planning:<version>
```

Afterwards the app can be opened at http://localhost:8080. Please consider mounting a volume at like this `volume:/level.db` to store the state of the app outside of the container itself (highly recommended).

## Architecture

The is written in Go and utilizes only a single dependency in LevelsDB. LevelsDB is the persistence layer of the app. The web app is handled using HTMX, AlpineJS and Tailwind.