<template>

    <footer class="footer">

        <pre>
            route params: {{$route.params | json}}
            {{ $root.$data | json }}
        </pre>

        <div class="row">
            <div class="col-sm-6 col-sm-offset-2">
                <a href="/forum">{{ $t('nav.forum') }}</a> |
                <a v-link="{ path: '/contact' }">{{ $t('nav.contact') }}</a>
                | <a href="/help">{{ $t('nav.help') }}</a>
                | <a href="/support">{{ $t('nav.support') }}</a>
                | <a href="/about">{{ $t('nav.about') }}</a><br/>
                <div class="footer_corp">&copy; {{ new Date().getFullYear() }} Example, Inc.</div>
            </div>
            <div class="col-sm-4 pull-right">
                <language-chooser></language-chooser>
            </div>
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
