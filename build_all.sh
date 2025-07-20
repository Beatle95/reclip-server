version=$(git describe --tags)
go build -o ${PWD}/out/server -ldflags "-X 'internal.applicationVersionString=$version'"
go build -C client_integration_tests -o ${PWD}/out/client_integration_tests -ldflags "-X 'internal.applicationVersionString=$version'"
