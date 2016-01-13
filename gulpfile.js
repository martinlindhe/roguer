var elixir = require('laravel-elixir');

elixir.config.production = true; // minify
elixir.config.sourcemaps = true;

elixir(function(mix) {

    mix
        .sass('app.scss')

        .browserify('app.js')

        .copy(
            'node_modules/bootstrap-sass/assets/fonts/bootstrap',
            'public/fonts/bootstrap'
        )

        .copy(
            'node_modules/font-awesome/fonts',
            'public/fonts/font-awesome'
        )

    ;

});
