<template>
    <nav class="navbar navbar-custom navbar-fixed-top">
        <div class="container-fluid">

            <bootstrap-dropdown-collapse>
                <button type="button" class="navbar-toggle collapsed" data-toggle="collapse" data-target="#navbar-collapse-top">
                    <span class="sr-only">Toggle navigation</span>
                    <span class="icon-bar"></span>
                    <span class="icon-bar"></span>
                    <span class="icon-bar"></span>
                </button>
                <a class="navbar-brand" v-link="{ path: '/' }">brand name</a>
            </bootstrap-dropdown-collapse>

            <div class="collapse navbar-collapse" id="navbar-collapse-top">
                <ul class="nav navbar-nav navbar-right">

                    <bootstrap-dropdown-toggle v-if="$root.token && $root.shop">
                        <a href="#" class="dropdown-toggle" data-toggle="dropdown" role="button" aria-expanded="false">
                            <i class="fa fa-home"></i> {{ $root.shop.name }}
                            <span class="caret"></span>
                        </a>

                        <ul class="dropdown-menu" role="menu">
                            <li><a v-link="{ path: '/shop/show/' + $root.shop.id }">{{ $t('nav.overview') }}</a></li>
                            <li class="divider"></li>
                            <li><h4>{{ $t('nav.settings') }}</h4></li>

                            <li role="presentation"><a href="/shop/setting/overview">{{ $t('nav.overview') }}</a></li>
                            <li role="presentation"><a href="/shop/setting/name/overview">{{ $t('nav.name') }}</a></li>
                            <li role="presentation"><a href="/shop/setting/order/overview">{{ $t('nav.orders') }}</a></li>
                            <li role="presentation"><a href="/shop/setting/company/overview">{{ $t('nav.company') }}</a></li>
                            <li role="presentation"><a href="/shop/setting/bank/overview">{{ $t('nav.bank_accounts') }}</a></li>
                            <li role="presentation"><a href="/shop/setting/contact/overview">{{ $t('nav.contact') }}</a></li>
                            <li role="presentation"><a href="/shop/setting/theme/overview">{{ $t('nav.theme') }}</a></li>
                            <li role="presentation"><a href="/shop/setting/text/overview">{{ $t('nav.text') }}</a></li>
                            <li role="presentation"><a href="/shop/setting/language/overview">{{ $t('nav.language') }}</a></li>
                            <li role="presentation"><a href="/shop/setting/product/overview">{{ $t('nav.product') }}</a></li>

                            <li role="presentation"><a href="/shop/setting/plugin/overview">{{ $t('nav.plugins') }}</a></li>

                            <li role="presentation"><a href="/shop/setting/logo/overview">{{ $t('nav.logo') }}</a></li>
                            <li role="presentation"><a href="/shop/setting/background/overview">{{ $t('nav.backgrounds') }}</a></li>

                            <li role="presentation"><a href="/shop/setting/campaign/overview">{{ $t('nav.campaigns') }}</a></li>

                            <li role="presentation"><a href="/shop/setting/extra/overview">{{ $t('nav.extras') }}</a></li>

                            <!--if (shop.hasExtra('message') -->
                            <li role="presentation"><a href="/shop/setting/extra/message/overview">{{ $t('nav.messages') }}</a></li>

                            <!-- if (shop.hasExtra('faq')) -->
                            <li role="presentation"><a href="/shop/setting/extra/faq/overview">{{ $t('nav.faq') }}</a></li>

                            <li role="presentation"><a href="/shop/setting/user/overview">{{ $t('nav.users') }} <span class="badge">X</span></a></li>
                            <li role="presentation"><a href="/shop/setting/host/overview">{{ $t('nav.hosts') }}</a></li>
                        </ul>
                    </bootstrap-dropdown-toggle>


                    <bootstrap-dropdown-toggle v-show="$root.token">
                        <a href="#" class="dropdown-toggle" data-toggle="dropdown" role="button" aria-expanded="false">
                            <i class="fa fa-shopping-cart"></i> {{ $t('nav.shops') }} <span class="caret"></span>
                        </a>

                        <ul class="dropdown-menu" role="menu">

                            <li v-for="shop in $root.shops">
                                <a v-link="{ path: '/shop/show/' + shop.id }">
                                    <i class="fa fa-arrow-right fa-fw"></i>&nbsp; {{ shop.name }}
                                </a>
                            </li>

                            <li class="divider"></li>
                            <li><a href="/shop/create"><i class="fa fa-plus fa-fw"></i>&nbsp; {{ $t('nav.new_shop') }}</a></li>
                        </ul>
                    </bootstrap-dropdown-toggle>

                    <bootstrap-dropdown-toggle v-show="$root.token">
                        <a href="#" class="dropdown-toggle" data-toggle="dropdown" role="button" aria-expanded="false">
                            {{ $root.username }} <span class="caret"></span>
                        </a>

                        <ul class="dropdown-menu" role="menu">
                            <li><a href="/user/setting/overview"><i class="fa fa-cogs fa-fw"></i>&nbsp; {{ $t('nav.settings') }}</a></li>
                            <li><a href="/user/account/overview"><i class="fa fa-user fa-fw"></i>&nbsp; {{ $t('nav.account') }}</a></li>
                            <li><a href="/user/subscription/overview"><i class="fa fa-money fa-fw"></i>&nbsp; {{ $t('nav.subscription') }}</a></li>
                            <li class="divider"></li>
                            <li><a href="/server/setting/overview"><i class="fa fa-server fa-fw"></i>&nbsp; {{ $t('nav.server') }}</a></li>
                            <li class="divider"></li>
                            <li><a href="#" @click="$root.logout"><i class="fa fa-sign-out fa-fw"></i>&nbsp; {{ $t('nav.sign_out') }}</a></li>
                        </ul>
                    </bootstrap-dropdown-toggle>

                    <li v-show="!$root.token">
                        <a v-link="{ path: '/auth/register' }"><i class="fa fa-plus fa-fw"></i>{{ $t('user.new') }}</a>
                    </li>

                    <li v-show="!$root.token">
                        <a v-link="{ path: '/auth/login' }"><i class="fa fa-sign-in fa-fw"></i>{{ $t('nav.sign_in') }}</a>
                    </li>

                </ul>
            </div>
        </div>
    </nav>
</template>

<script>
import Vue from 'vue';

import BootstrapDropdownToggle from './BootstrapDropdownToggle.vue';
import BootstrapDropdownCollapse from './BootstrapDropdownCollapse.vue';

export default {
    calculated: {
        shop: function() {
            return this.shops[0];
        }
    },
    components: {
        BootstrapDropdownToggle,
        BootstrapDropdownCollapse,
    },
}
</script>
