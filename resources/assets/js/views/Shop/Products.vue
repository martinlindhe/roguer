<template>

    <div class="container">
        <div class="row">
            <div class="col-sm-3 col-xs-12">
                <shop-nav-bar></shop-nav-bar>
            </div>

            <div class="col-sm-9 col-xs-12">

                <div class="panel panel-default">
                    <div class="panel-body">

                        <breadcrumb :list="[
                            {title: $root.shop.name, url: '/shop/show/' + $root.shop.id},
                            {title: $t('nav.products'), url: '/shop/products'}
                            ]"></breadcrumb>

                        @foreach ($shop->products as $product)
                            @include('_partials.product-admin-list', $product)
                        @endforeach

                        <a href="/shop/product/add">{{ $t('nav.add_product') }}</a>

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
    }
}
</script>
