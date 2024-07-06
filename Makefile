build: buildJsAndCss buildGo copyViews

clean: dist
	rm -r dist

buildJsAndCss: esbuild.mjs
	node esbuild.mjs

buildGo: cmd/main.go
	go build -o dist/meal-planner ./cmd

copyViews: views
	@mkdir dist/views
	cp views/* dist/views

run: dist/meal-planner
	cd ./dist; ./meal-planner
