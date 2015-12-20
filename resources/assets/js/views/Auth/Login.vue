<template>

    <form @submit.prevent="postLogin" method="post" action="login" class="form-horizontal">

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

            {{ $t('auth.or_login_with_existing_account') }}

            <h3>{{ $t('nav.sign_in') }}</h3>
        </div>

        <div class="form-group">
            <div class="col-md-5 col-md-offset-4">
                <input v-model="email" id="email" class="form-control" placeholder="{{ $t('nav.email') }}">
            </div>
        </div>

        <div class="form-group">
            <div class="col-md-5 col-md-offset-4">
                <input type="password" v-model="password" class="form-control" placeholder="{{ $t('nav.password') }}">
            </div>
        </div>

        <div class="form-group">
            <div class="col-md-5 col-md-offset-4">
                <button v-show="email && password" class="btn btn-primary form-control">{{ $t('nav.sign_in') }}</button>
            </div>
        </div>

        <div class="form-group">
            <div class="col-md-5 col-md-offset-4">
                <a href="/password/email">{{ $t('nav.has_forgot_password') }}</a>
            </div>
        </div>

    </form>

</template>

<script>
import Vue from 'vue';

export default {
    data: function() {
        return {
            email: '',
            password: ''
        }
    },
    ready() {
        document.getElementById('email').focus();
    },
    methods: {
        postLogin: function(el) {

            var resource = this.$resource('/api/auth/login');
            resource.save({email: this.email, password: this.password}, function (data, status, request) {
                console.log("login ok");

                // jwt-auth
                Vue.http.headers.common['Authorization'] = 'Bearer ' + data.token;

                this.$root.username = data.username;
                this.$root.token = data.token;
                this.$root.shops = data.shops;
                this.$root.shop = data.shops[0];
                this.$root.email = this.email;

                window.sessionStorage.setItem('_token', this.$root.token);
                window.localStorage.setItem('_username', this.$root.username);
                window.localStorage.setItem('_email', this.$root.email);
                window.localStorage.setItem('_shops', JSON.stringify(this.$root.shops));
                window.localStorage.setItem('_shop', JSON.stringify(this.$root.shop));

                this.$root.registerRefreshTokenTimer();

                this.$route.router.go('/');

            }).error(function (data, status, request) {
                // handle error
                console.log("login error: " + data.status);
                console.log(request);
            });
        }
    }
}
</script>
