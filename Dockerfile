FROM golang:1.22

# Install java 
RUN apt-get update && apt-get install -y openjdk-17-jre

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o server .

EXPOSE 8080 25565

CMD ["./server"]