go build -o dist/mtk -i ./src/*.go
cd dist && PROG=mtk source ../autocomplete/bash_autocomplete
cd ../
