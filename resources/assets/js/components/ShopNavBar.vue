<template>
    <nav class="navbar navbar-default">

        <bootstrap-dropdown-collapse>
            <button type="button" class="navbar-toggle collapsed" data-toggle="collapse" data-target="#navbar-collapse-shop">
                <span class="sr-only">Toggle navigation</span>
                {{ $root.shop.name }}
            </button>
        </bootstrap-dropdown-collapse>

        <div class="collapse navbar-collapse" id="navbar-collapse-shop">
            <ul class="nav nav-pills nav-stacked">
                <li role="presentation">
                    <!-- XXX FIXME: it is the li tag who need class="active" if path match, not the a...-->
                    <a v-link="{ path: '/shop/show/' + $root.shop.id, activeClass: 'active' }">{{ $t('nav.overview') }}</a>
                </li>

                <li role="presentation">
                    <a v-link="{ path: '/shop/products', activeClass: 'active' }">{{ $t('nav.products') }}
                        <span class="badge">{{ $root.shop.products }}</span>
                    </a>
                </li>

                <li role="presentation">
                    <a v-link="{ path: '/shop/product/model/overview', activeClass: 'active' }"> &gt; {{ $t('nav.models') }}</a>
                </li>

                <li role="presentation">
                    <a v-link="{ path: '/shop/order/overview', activeClass: 'active' }">{{ $t('nav.orders') }}
                        <span class="badge">{ count(shop->orders) }</span>
                    </a>
                </li>
            </ul>

            <p v-if="$root.shop.logos">
                <!-- FIXME not showing logo -->
                <br/>
                <img src="/image/resize/width/{ $shop->logos[0]->image->id }/32"/>
            </p>
        </div>

    </nav>
</template>

<script>
import Vue from 'vue';

//import BootstrapDropdownToggle from './BootstrapDropdownToggle.vue';
import BootstrapDropdownCollapse from './BootstrapDropdownCollapse.vue';

export default {
    components: {
        //BootstrapDropdownToggle,
        BootstrapDropdownCollapse,
    },
    methods: {
    }
}
</script>
