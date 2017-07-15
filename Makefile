default: clean setup server client

clean:
	rm -rf dist

setup:
	mkdir dist; mkdir dist/client; mkdir dist/server

server: FORCE
	GOOS=windows go build -o ./dist/server/WorldBackup.exe ./server; cp server/config.json dist/server

client: FORCE
	cd client; yarn run build; cd ..; cp -R client/dist/* dist/client

FORCE: