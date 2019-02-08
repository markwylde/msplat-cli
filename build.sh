go build -o dist/mtk -ldflags="-s -w" -i ./src/*.go
cd dist && PROG=mtk source ../autocomplete/bash_autocomplete
cd ../
