.PHONY: mockgen build test docker docker-build docker-push

GIT_COMMIT   ?= $(shell git rev-parse HEAD)
GO_VERSION   ?= $(shell go version | awk {'print $$3'})
DOCKER_IMAGE ?= figmentnetworks/oasishub-indexer
DOCKER_TAG   ?= latest

# Generate mocks
mockgen:
	@echo "[mockgen] generating mocks"
	@mockgen -destination mock/store/mocks.go github.com/figment-networks/oasishub-indexer/store DatabaseStore,SyncablesStore,ReportsStore,SystemEventsStore,BlockSeqStore,DebondingDelegationSeqStore,DelegationSeqStore,StakingSeqStore,TransactionSeqStore,ValidatorSeqStore,BlockSummaryStore,ValidatorSummaryStore,AccountAggStore,ValidatorAggStore
	@mockgen -destination mock/indexer/mocks.go github.com/figment-networks/oasishub-indexer/indexer AccountAggCreatorTaskStore,BackfillSourceStore,BalanceEventPersistorTaskStore,BlockSeqCreatorTaskStore,BlockSeqPersistorTaskStore,ConfigParser,DebondingDelegationSeqCreatorTaskStore,DelegationSeqCreatorTaskStore,SourceIndexStore,StakingSeqCreatorTaskStore,SyncerPersistorTaskStore,SyncerTaskStore,SystemEventCreatorStore,TransactionSeqCreatorTaskStore,ValidatorAggCreatorTaskStore,ValidatorAggPersistorTaskStore,ValidatorSeqCreatorTaskStore,ValidatorSeqPersistorTaskStore
	@mockgen -destination mock/client/mocks.go github.com/figment-networks/oasishub-indexer/client AccountClient,BlockClient,ChainClient,EventClient,StateClient,TransactionClient,ValidatorClient

# Build the binary
build:
	go build \
		-ldflags "\
			-X github.com/figment-networks/oasishub-indexer/cli.gitCommit=${GIT_COMMIT} \
			-X github.com/figment-networks/oasishub-indexer/cli.goVersion=${GO_VERSION}"

# Run tests
test:
	go test -race -cover ./...

# Build a local docker image for testing
docker:
	docker build -t oasishub-indexer -f Dockerfile .

# Build a public docker image
docker-build:
	docker build \
		-t ${DOCKER_IMAGE}:${DOCKER_TAG} \
		-f Dockerfile \
		.

# Push docker images
docker-push:
	docker tag ${DOCKER_IMAGE}:${DOCKER_TAG} ${DOCKER_IMAGE}:latest
	docker push ${DOCKER_IMAGE}:${DOCKER_TAG}
	docker push ${DOCKER_IMAGE}:latest