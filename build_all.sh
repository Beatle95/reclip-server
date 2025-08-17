SCRIPT_DIR=$(cd -- "$( dirname -- "${BASH_SOURCE[0]}" )" &> /dev/null && pwd)
cd ${SCRIPT_DIR}
version=$(git describe --abbrev=0)
go build -o out/server -ldflags "-X 'internal.applicationVersionString=$version'"
go build -C client_integration_tests -o ${SCRIPT_DIR}/out/client_integration_tests -ldflags "-X 'internal.applicationVersionString=$version'"
