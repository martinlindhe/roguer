<template>

    <div class="container">
        <div class="row">
            <div class="col-sm-3 col-xs-12">
                <shop-nav-bar></shop-nav-bar>
            </div>

            <div class="col-sm-9 col-xs-12">

                <div class="panel panel-default">
                    <div class="panel-body">

                        <breadcrumb :list="[ {title: $root.shop.name, url: '/shop/show/' + $root.shop.id} ]"></breadcrumb>

                        <span v-if="$root.shop.open">
                            <p>
                                {{ $t('shop.is_open') }}
                            </p>
                        </span>
                        <span v-else>
                            <p>
                                {{ $t('shop.is_closed') }}
                            </p>

                            <p>
                                {{ $t('nav.you_can_open_the_shop_from_the') }} <a href="/shop/setting/overview">{{ $t('nav.settings') }}</a>.
                            </p>
                        </span>

                        <p>
                            {{ $t('shop.primary_host') }}:
                            <a href="{{ $root.shop.primary_host }}" target="_blank">{{ $root.shop.primary_host }}</a>
                        </p>

                        {{ $t('nav.created') }}
                        <when-for-humans :time="$root.shop.created_at"></when-for-humans>

                    </div>
                </div>
            </div>

        </div>
    </div>

</template>

<script>
// TODO:
//   - a RELOAD button, that does a json request to get current shops and update $root.shops, $root.shop

import Breadcrumb from './../../components/Breadcrumb.vue';
import WhenForHumans from './../../components/WhenForHumans.vue';

import ShopNavBar from './../../components/ShopNavBar.vue';

export default {
    components: {
        Breadcrumb,
        WhenForHumans,
        ShopNavBar,
    },
    ready() {

        // FACT: visiting a shop overview makes it the currently selected one
        for (var i = 0; i < this.$root.shops.length; i++) {
            if (this.$root.shops[i].id == this.$route.params.id) {
                this.$root.shop = this.$root.shops[i];
                window.localStorage.setItem('_shop', JSON.stringify(this.$root.shop));
            }
        }
    }
}
</script>
