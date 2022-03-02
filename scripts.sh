# Compile and run coordinator
go build -o coordinator ./cmd/coordinator && ./coordinator
# Compile and run attacker
go build -o attacker ./cmd/attacker && ./attacker

mkdir -p docker-images/coordinator
mkdir -p docker-images/attacker

### ATTACKER
# Compile for alpine
docker run -v $(pwd)/docker-images/attacker:/outfile -v $(pwd):/app --workdir="/app" --env CGO_ENABLED=0 golang:latest go build -o /outfile ./cmd/attacker

## In folder with binary
# Build
docker build --no-cache -t "tbsitg/artillery-attacker:v1.0.0" .
# Publish
docker push tbsitg/artillery-attacker:v1.0.0

### COORDINATOR
# Compile for alpine
docker run -v $(pwd)/docker-images/coordinator:/outfile -v $(pwd):/app --workdir="/app" --env CGO_ENABLED=0 golang:latest go build -o /outfile ./cmd/coordinator

## In folder with binary
# Build
docker build --no-cache -t "tbsitg/artillery-coordinator:v1.0.0" .
# Publish
docker push tbsitg/artillery-coordinator:v1.0.0
