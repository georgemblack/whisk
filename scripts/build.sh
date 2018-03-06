sass ./resources/main.scss ./public/styles/main.css --style compressed
sass ./resources/themes/empty.scss ./public/styles/empty.css --style compressed
rm -rf ./public/styles/*.css.map ./.sass-cache/
go run cmd/whisk/main.go