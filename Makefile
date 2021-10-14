.PHONY: build clean package deploy log create-bucket remove-bucket describe remove outputs urls

STACK=ederoyd

build:
	env GOOS=linux go build -ldflags="-s -w" -o bin/hello src/handlers/hello/hello.go

offline: build
	sam local start-api

offline-skip: build
	sam local start-api --skip-pull-image

clean:
	rm -rf ./bin

package: build
	sam package \
		--template-file template.yml \
		--output-template-file packaged.yml \
		--s3-bucket $(STACK)

deploy:
	sam deploy \
		--template-file packaged.yml \
		--stack-name $(STACK) \
		--capabilities CAPABILITY_IAM

log: 
	sam logs --stack-name $(STACK) --name HelloFunction

describe: 
	aws cloudformation describe-stacks --stack-name $(STACK)

remove: 
	aws cloudformation delete-stack --stack-name $(STACK)

outputs: 
	aws cloudformation describe-stacks \
    	--stack-name $(STACK) --query 'Stacks[].Outputs'

# create-bucket: 
# 	aws s3 mb s3://$(STACK)

# remove-bucket: 
# 	aws s3 rb s3://$(STACK) --force
