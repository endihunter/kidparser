FROM gobuffalo/buffalo:v0.9.5

RUN mkdir -p $GOPATH/src/recipes
WORKDIR $GOPATH/src/recipes

ADD . .
RUN go get $(go list ./... | grep -v /vendor/)
RUN buffalo build --static -o /bin/app

EXPOSE 3000

# Comment out to run the migrations before running the binary:
# CMD /bin/app migrate; /bin/app
CMD /bin/app
