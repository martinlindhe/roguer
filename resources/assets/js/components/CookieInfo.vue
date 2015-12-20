<template>
    <div>
        <bootstrap-alert v-if="!seen" type="warning" @click="click()">
            {{ $t('auth.cookie_banner') }}
            <a href="/cookie/policy/{{ this.$root.locale }}" target="_blank">{{ $t('nav.learn_more') }}</a>
        </bootstrap-alert>
    </div>
</template>

<script>
import Vue from 'vue';

import BootstrapAlert from './BootstrapAlert.vue';
/*
Vue.transition('dismissed', {
    css: false,
    leave: function (el, done) {
        console.log("dismissed yea!");
    }
});
*/

export default {
    data() {
        return {
            seen: 0,
        }
    },
    components: {
        BootstrapAlert,
    },
    methods: {
        click() {
            // XXX UGLY HACK: 25 nov, the vue-strap Alert component has a dismissable prop, but cant
            // override function trigger, so i made whole alert clickable. fix it properly, maybe with a transition?

            // XXXX also css the div is only invisible, not gone
            console.log("clikck");

            this.seen = 1;
            window.localStorage.setItem('_cookie_info', 1);
        }
    }, ready() {
        this.seen = window.localStorage.getItem('_cookie_info');
    }
}
</script>

