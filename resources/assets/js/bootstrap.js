// XXX::
// './resources/assets/js/filereader-0.99.js',


import Vue from 'vue';
import VueAsyncData from 'vue-async-data';
import VueInternationalization from 'vue-i18n';
import VueResource from 'vue-resource';
import VueRouter from 'vue-router';
import VueValidator from 'vue-validator';

Vue.config.debug = true;


// XXX move to function file
Array.prototype.contains = function(obj) {
    var i = this.length;
    while (i--) {
        if (this[i] === obj) {
            return true;
        }
    }
    return false;
}


// generated from Laravel translations with "php artisan vue-i18n:generate"
import Locales from './vue-i18n-locales.generated.js';

var locale = window.localStorage.getItem('_locale');
if (!locale) {
    locale = 'sv_SE';
}

Vue.use(VueInternationalization, {
    lang: locale,
    locales: Locales
});

Vue.use(VueAsyncData);
Vue.use(VueResource);
Vue.use(VueRouter);
Vue.use(VueValidator);

import './moment-locales.js';

import moment from 'moment';
moment.locale(locale);



// **********
// ** ROUTES
// **********
const router = new VueRouter({
    hashbang: false
});

import App from './App.vue';
import MainView from './views/Main.vue';
import AuthLoginView from './views/Auth/Login.vue';
import AuthRegisterView from './views/Auth/Register.vue';
import ContactView from './views/Corporate/Contact.vue';
import ShopShowView from './views/Shop/Show.vue';
import ShopProductsView from './views/Shop/Products.vue';

router.map({
    '/': { component: MainView },
    '/auth/login': { component: AuthLoginView },
    '/auth/register': { component: AuthRegisterView },
    '/contact': { component: ContactView },
    '/shop/show/:id': { component: ShopShowView },
    '/shop/products': { component: ShopProductsView },
});

// Redirect to the home route if any routes are unmatched
router.redirect({
    '*': '/'
});

router.start(App, '#app');

