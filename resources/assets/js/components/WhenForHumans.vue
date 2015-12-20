<template>
    <bootstrap-tooltip
        trigger="hover"
        effect="scale"
        placement="bottom"
        content="Lorem ipsum dolor sit amet consectetur adipisicing elit, sed do eiusmod"
    >
        {{ val() }}
    </bootstrap-tooltip>
</template>

<script>

import moment from 'moment';

// XXX dont work + weird errors in console:
//import BootstrapTooltip from './BootstrapTooltip.vue';

export default {
    props: ['time'],
    components: [
        //BootstrapTooltip,
    ],
    methods: {
        val() {
            return moment.unix(this.time).fromNow();
        }
    }
}
</script>
