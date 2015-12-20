<template>
    <div>
        <div v-for="lang in $root.locales">
            <language :name="lang.name" :code="lang.code"></language>
        </div>
    </div>
</template>

<script>
import Vue from 'vue';

export default {

    components: {
        language: {
            props: ['code', 'name'],
            template: '<p @click="Select" style="cursor: pointer">{{ name }}</p>',
            methods: {
                Select() {
                    console.log("changing language to " + this.code);
                    Vue.config.lang = this.code;

                    window.localStorage.setItem('_locale', this.code);
                    this.$root.locale = this.code;

                    // XXX dont seem to work to change on the fly, see https://github.com/kazupon/vue-i18n/issues/2
                    // HACK so instead we do a full reload here
                    location.reload();
                }
            }
        }
    }
}
</script>

