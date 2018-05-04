NAME=fuelsales-rpt
VERSION=0.2
PORT_MAP=3011:3011

AWS_ACCOUNT=407205661819
AWS_REGION=ca-central-1
TAG=galesd/$(NAME):$(VERSION)
REPO=$(AWS_ACCOUNT).dkr.ecr.$(AWS_REGION).amazonaws.com/$(TAG)


default: buildp

# production build
buildp: buildgop build push

# build go for production targeting a linux OS
buildgop:
	CGO_ENABLED=0 GOOS=linux go build -ldflags "-s" -a -installsuffix cgo -o $(NAME) .

# build go for development
buildgod:
	go build -o $(NAME) .

# standard docker build
build:
	docker build --rm --tag=$(REPO) .

push:
	docker push $(REPO)

# run in development
rundd:
	docker run --name $(NAME) -p $(PORT_MAP) --env Stage=test -d $(REPO)