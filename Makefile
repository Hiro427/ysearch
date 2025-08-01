gorun: 
	go run main.go 

watch: 
	air

tailwind: 
	bun run tailwindcss -i views/css/styles.css -o public/styles.css --minify --watch

templ: 
	templ generate --watch

