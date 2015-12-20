<template>

    <form @submit.prevent="postRegister" class="form-horizontal">

        <div class="col-md-5 col-md-offset-4">
            <h3>{{ $t('auth.sign_in_with_your_social_network') }}</h3>

            <nav class="navbar navbar-default">
                <ul class="nav navbar-nav">
                    <li>
                        <a href="/oauth/facebook"><i class="fa fa-facebook fa-fw"></i> Facebook</a>
                    </li>
                    <li>
                        <a href="/oauth/instagram"><i class="fa fa-instagram fa-fw"></i> Instagram</a>
                    </li>
                    <li>
                        <a href="/oauth/live"> Microsoft Live</a>
                    </li>
                </ul>
            </nav>

            {{ $t('auth.or_register_new_user') }}

            <h3>{{ $t('user.new') }}</h3>
        </div>

        <div class="form-group">
            <div class="col-md-5 col-md-offset-4">
                <input v-model="username" class="form-control" placeholder="{{ $t('nav.username') }}">
            </div>
        </div>

        <div class="form-group">
            <div class="col-md-5 col-md-offset-4">
                <input v-model="email" class="form-control" placeholder="{{ $t('nav.email') }}">
            </div>
        </div>

        <div class="form-group">
            <div class="col-md-5 col-md-offset-4">
                <input type="password" v-model="password" class="form-control" placeholder="{{ $t('nav.password') }}">

            </div>
        </div>

        <div class="form-group">
            <div class="col-md-5 col-md-offset-4">
                <input type="password" v-model="password_confirmation" class="form-control" placeholder="{{ $t('nav.confirm_password') }}">
            </div>
        </div>

        <div class="form-group">
            <div class="col-md-5 col-md-offset-4">
                <input type="checkbox" v-model="terms">
                &nbsp; <a href="/legal/terms/{{ $root.locale }}" target="_blank">{{ $t('nav.agree_to_terms') }}</a>
            </div>
        </div>

        <div class="form-group">
            <div class="col-md-5 col-md-offset-4">
                <button v-show="username && email && password && password_confirmation && terms" class="btn btn-primary form-control">{{ $t('nav.register') }}</button>
            </div>
        </div>

    </form>

</template>

<script>
import Vue from 'vue';

export default {
    data: function() {
        return {
            username: '',
            email: '',
            password: '',
            password_confirmation: '',
            terms: false,
        }
    },
    methods: {
        postRegister: function() {
            /*
            var resource = this.$resource('/api/auth/login');
            resource.save({email: this.email, password: this.password}, function (data, status, request) {
                console.log("login ok");

                // jwt-auth
                Vue.http.headers.common['Authorization'] = 'Bearer ' + data.token;

                this.$root.username = data.username;
                this.$root.token = data.token;

                window.sessionStorage.setItem('_token', data.token);
                window.localStorage.setItem('_username', data.username);

                this.$route.router.go('/');

            }).error(function (data, status, request) {
                // handle error
                console.log("login error: " + data.status);
                console.log(request);
            });
            */

        }
    }
}
</script>
