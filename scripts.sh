# Compile and run coordinator
go build -o coordinator ./cmd/coordinator && ./coordinator
# Compile and run attacker
go build -o attacker ./cmd/attacker && ./attacker

mkdir -p docker-images/coordinator
mkdir -p docker-images/attacker

### ATTACKER
# Compile for alpine
docker run -v $(pwd)/docker-images/attacker:/outfile -v $(pwd):/app --workdir="/app" --env CGO_ENABLED=0 golang:latest go build -o /outfile ./cmd/attacker
# Build
docker build --no-cache -t "tbsitg/artillery-attacker:latest" ./docker-images/attacker
# Publish
docker push tbsitg/artillery-attacker:latest
# Run
docker run -it tbsitg/artillery-attacker:latest

### COORDINATOR
# Compile for alpine
docker run -v $(pwd)/docker-images/coordinator:/outfile -v $(pwd):/app --workdir="/app" --env CGO_ENABLED=0 golang:latest go build -o /outfile ./cmd/coordinator
# Build
docker build --no-cache -t "tbsitg/artillery-coordinator:latest" ./docker-images/coordinator
# Publish
# docker push tbsitg/artillery-coordinator:latest
# Run
docker run -it -v $(pwd)/targets.yml:/app/targets.yml -p 9000:9000 tbsitg/artillery-coordinator:latest
