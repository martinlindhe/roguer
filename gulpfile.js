var elixir = require('laravel-elixir');

elixir.config.sourcemaps = false;

elixir(function(mix) {

    // shell.task(['php artisan vue-i18n:generate']);

    mix
        .sass('app.scss')

        .browserify(
            './resources/assets/js/bootstrap.js',
            'public/js/app.js'
        )
/*
        .version([
            'css/app.css',
            'js/app.js'
        ])
*/
/*
        .browserSync({
            proxy: 'example.app'
        })
*/
        .copy(
            'node_modules/bootstrap-sass/assets/fonts/bootstrap',
            'public/fonts/bootstrap'
        )

        .copy(
            'node_modules/font-awesome/fonts',
            'public/fonts/font-awesome'
        )

        .copy(
            'node_modules/flag-icon-css/flags/1x1',
            'public/flags/1x1'
        )

        .copy(
            'node_modules/flag-icon-css/flags/4x3',
            'public/flags/4x3'
        )

    ;

});
