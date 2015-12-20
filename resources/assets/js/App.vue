<template>
    <div>
        <nav-bar></nav-bar>

        <cookie-info></cookie-info>

        <router-view></router-view>

        <corporate-footer></corporate-footer>
    </div>
</template>

<script>
import moment from 'moment';
import Vue from 'vue';

import NavBar from './components/NavBar.vue';
import CookieInfo from './components/CookieInfo.vue';
import CorporateFooter from './components/Footer/Corporate.vue';

export default {
    data() {
        return {
            username: '',
            email: '',
            locale: '', // "sv_SE"
            token: '', // jwt-auth
            _refreshTokenTimer: null,
            // all accessible shops
            shops: [],

            // currently selected shop
            shop: {},

            locales: [
                {
                    code: 'sv_SE',
                    name: 'Svenska'
                },
                {
                    code: 'en_US',
                    name: 'American English'
                }
            ]
        }
    },
    components: {
        NavBar,
        CookieInfo,
        CorporateFooter,
    },
    ready() {

        this.username = window.localStorage.getItem('_username');
        this.email = window.localStorage.getItem('_email');
        this.locale = window.localStorage.getItem('_locale');
        this.token = window.sessionStorage.getItem('_token');
        this.shops = JSON.parse(window.localStorage.getItem('_shops'));
        this.shop = JSON.parse(window.localStorage.getItem('_shop'));

        if (!Boolean(this.locale)) {
            this.locale = this.defaultLocale();
            console.log('no locale pref found, defaulting to ' + this.locale);
        }

        console.log("[booted] vue " + Vue.version + ", locale " + this.locale);
        if (Boolean(this.token)) {
            console.log("re-using token: " + this.token);

            // jwt-auth
            Vue.http.headers.common['Authorization'] = 'Bearer ' + this.token;

            this.refreshToken();
        }
    },
    methods: {
        refreshToken() {
            console.log("refreshing api token ... " + moment());

            this.registerRefreshTokenTimer();

            this.$http.get('/api/auth/refresh-token', function (data, status, request) {
                console.log("success refreshing token");

                Vue.http.headers.common['Authorization'] = 'Bearer ' + data.token;
                this.$root.token = data.token;
                window.sessionStorage.setItem('_token', data.token);

            }).error(function (data, status, request) {
                console.error("error refreshing token, ending session");
                console.log(data);

                this.clearRefreshTokenTimer();
                this.endSession();
            });
        },
        registerRefreshTokenTimer() { // private

            //this.clearRefreshTokenTimer();

            if (!Boolean(this.token)) {
                console.error("registerRefreshTokenTimer called without a token");
                return;
            }

            var intervalMillisec = 300 * 1000; // 300s = 5m

            this._refreshTokenTimer = setTimeout(this.refreshToken, intervalMillisec);
        },
        clearRefreshTokenTimer() { // private
            if (Boolean(this._refreshTokenTimer)) {
                clearTimeout(this._refreshTokenTimer);
            }

            console.log("XXXX clearing Authorization header");
            Vue.http.headers.common['Authorization'] = '';
        },
        defaultLocale() {

            var loc = navigator.languages != undefined ? navigator.languages[0] : navigator.language;

            // "en-US" => "en_US"
            loc = loc.replace('-', '_');
            //console.log("browser locale is " + loc);

            var ok = ['sv_SE', 'en_US'];
            if (!ok.contains(loc)) {
                console.log("unsupported browser locale " + loc + ", defaulting to " + Vue.config.lang);
                return Vue.config.lang;
            }
            return loc;
        },
        logout() {
            console.log("logging out");
            this.clearRefreshTokenTimer();

            if (!Boolean(this.token)) {
                return;
            }

            var resource = this.$http.get('/api/auth/logout', function (data, status, request) {
                this.endSession();

                this.$route.router.go('/');

                // XXX show alert popup with "You are now logged out."
            }).error(function (data, status, request) {
                console.error("error logging out, clearing session");

                this.clearRefreshTokenTimer();
                this.endSession();
                this.$route.router.go('/');
            });
        }, endSession() {
            this.$root.token = '';

            window.sessionStorage.setItem('_token', '');
            Vue.http.headers.common['Authorization'] = '';
        }
    }
}
</script>
