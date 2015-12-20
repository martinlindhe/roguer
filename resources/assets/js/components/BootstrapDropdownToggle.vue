<template>
    <li class="dropdown">
        <slot></slot>
        <slot name="dropdown-menu"></slot>
    </li>
</template>

<script>
import EventListener from './../EventListener';

export default {
    methods: {
        toggleDropdown(e) {
            e.preventDefault();
            this.$el.classList.toggle('open');
        }
    },
    ready() {
        const toggle = this.$el.querySelector('[data-toggle="dropdown"]');
        if (toggle) {
            toggle.style.borderRadius = '4px';
            toggle.addEventListener('click', this.toggleDropdown);
        }
        this._closeEvent = EventListener.listen(window, 'click', (e)=> {
            if (!this.$el.contains(e.target)) {
                this.$el.classList.remove('open');
            }
        })
    },
    beforeDestroy() {
        if (this._closeEvent) {
            this._closeEvent.remove();
        }
    }
}
</script>
