FROM golang:1.17 AS build

# Populate the module cache based on the go.{mod,sum} files.
COPY ./go.mod ./go.sum src/vecapp/
WORKDIR src/vecapp

# add workaround to add own repositoriy packages
RUN git config --global \
url."https://vecLibsToken:oac9pW1xsTMYbxK4DeYK@gitlab.vecomentman.com/".insteadOf "https://gitlab.vecomentman.com/" && \
go list -m gitlab.vecomentman.com/libs/logger && \
go list -m gitlab.vecomentman.com/libs/sso && \
go list -m gitlab.vecomentman.com/libs/email

# Download all dependencies.
RUN go mod download

# copy all sources to container
COPY . .

# Build the Go app
RUN CGO_ENABLED=0 GOOS=linux go build -o /go/bin/main ./cmd/server/main.go
#RUN go get -d -v ./...
#RUN go install -v ./...


# abstract build layers forms the final image
FROM scratch
# copy certs files into scratch image
COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

COPY --from=build /go/bin/main .

# copy additional data
COPY ./app/config/config.service.yml go/src/vecapp/app/config/config.service.yml
COPY ./app/config/secret.service.yml go/src/vecapp/app/config/secret.service.yml
COPY ./app/config/config.logger.yml go/src/vecapp/app/config/config.logger.yml
COPY ./email-templates/dist go/src/vecapp/email-templates/dist
COPY ./repository/migrations go/src/vecapp/repository/migrations

# start application
ENTRYPOINT ["/main"]