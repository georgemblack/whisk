sass ./resources/main.scss ./public/styles/main.css --style compressed
sass ./resources/themes/minimal.scss ./public/styles/minimal.css --style compressed
rm -rf ./public/styles/*.css.map ./.sass-cache/