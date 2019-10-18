FROM golang
WORKDIR /wenkuProject
COPY / .
EXPOSE 8991
ENTRYPOINT ["./main"]