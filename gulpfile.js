var elixir = require('laravel-elixir');

elixir.config.sourcemaps = false;

elixir(function(mix) {

    mix
        .sass('app.scss')

        .browserify(
            [
                './resources/assets/js/game.js',
            ],
            'public/js/app.js'
        )

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
