<template>

    <footer class="footer">

        <div class="row">
            <div class="col-sm-6 col-sm-offset-2">
                <a v-link="{ path: '/contact' }">{{ $t('nav.contact') }}</a>
                <div class="footer_corp">&copy; {{ new Date().getFullYear() }} Example, Inc.</div>
            </div>
            <div class="col-sm-4 pull-right">
                <language-chooser></language-chooser>
            </div>
        </div>

        <div class="row">
            <pre>
                {{ $root.$data | json }}
            </pre>
        </div>

    </footer>
</template>

<script>
import Vue from 'vue';

import LanguageChooser from './../LanguageChooser.vue';

export default {
    components: {
        LanguageChooser,
    }
}
</script>
