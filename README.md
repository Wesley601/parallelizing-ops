# Parallelizing comparison, NodeJs child process vs Golang Coroutines

This is a comparison based on video [how to Migrate 1M items from MongoDB to Postgres in just a few minutes](https://youtu.be/EnK8-x8L9TY).

## Running

You'll need to install Docker and Docker compose to be able to spin up the DBs instances, after that run:
- docker-compose up -d

### for node
- cd node-version
- npm ci
- npm run seed
- npm start

### for Golang
- cd go-version
- go build -o main
- ./main
