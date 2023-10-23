FROM golang:1.21 AS build

# Populate the module cache based on the go.{mod,sum} files.
COPY ./go.mod ./go.sum src/service-app/
WORKDIR src/service-app

COPY ./.netrc /root/.netrc
RUN chmod 600 /root/.netrc

# add workaround to add own repositoriy packages
#RUN go list -m github.com/VerzCar/vyf-lib-logger && \
#go list -m github.com/VerzCar/vyf-lib-awsx

# Download all dependencies.
RUN go mod download

# copy all sources to container
COPY . .

# Build the Go app
RUN CGO_ENABLED=0 GOOS=linux go build -o /go/bin/main ./cmd/server/main.go
#RUN go get -d -v ./...
#RUN go install -v ./...


# abstract build layers forms the final image
#FROM scratch # use from scratch if no shell is needed
FROM alpine
# copy certs files into scratch image
COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

COPY --from=build /go/bin/main .

# copy additional data
COPY ./app/config/config.service.yml go/src/service-app/app/config/config.service.yml
#COPY ./app/config/secret.service.yml go/src/service-app/app/config/secret.service.yml
COPY ./app/config/config.logger.yml go/src/service-app/app/config/config.logger.yml
COPY ./repository/migrations go/src/service-app/repository/migrations

# start application
CMD ["./main"]
