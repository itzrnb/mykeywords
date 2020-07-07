#FROM golang:1.14-alpine AS build

#WORKDIR /src/
#COPY main.go go.* /src/
#RUN CGO_ENABLED=0 go build -o /bin/MyKeywords

#FROM scratch
#COPY --from=build /bin/MyKeywords /bin/MyKeywords
#ENTRYPOINT ["/bin/MyKeywords"]

FROM golang

COPY main.go go.* /src/
WORKDIR /src/
RUN CGO_ENABLED=0 go build -o /src/MyKeywords /src/main.go
CMD ["/src/MyKeywords"]