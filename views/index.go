package views

import (
	"bytes"
)

func Index() string {
	var _buffer bytes.Buffer
	_buffer.WriteString("<!DOCTYPE html>\n<html>\n<head>\n    <meta charset=\"utf-8\" />\n    <meta http-equiv=\"X-UA-Compatible\" content=\"IE=edge\">\n    <meta name=\"viewport\" content=\"width=device-width, initial-scale=1\">\n    <title>rogue</title>\n    <link href=\"/css/app.css\" rel=\"stylesheet\" type=\"text/css\">\n</head>\n<body>\n    <noscript>\n        <div class=\"alert alert-danger\" role=\"alert\">\n            <h4>{{ trans('help.please_enable_javascript') }}</h4>\n            {{ trans('help.enable_javascript_info') }}\n        </div>\n    </noscript>\n\n    <!--[if lt IE 10]>\n    <div>\n        {!! trans('help.old_ie', ['link' => 'http://ie.microsoft.com/']) !!}\n    </div>\n    <![endif]-->\n\n\n    <script src=\"/js/phaser.js\"></script>\n    <script src=\"/js/app.js\"></script>\n</body>\n</html>")

	return _buffer.String()
}
