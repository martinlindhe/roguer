<style>
.fade-transition {
    transition: opacity 0.3s ease;
}
.fade-enter, .fade-leave {
    height: 0;
    opacity: 0;
}
.alert.top {
    position: fixed;
    top: 30px;
    margin: 0 auto;
    left: 0;
    right: 0;
    z-index: 2;
}
.alert.top-right {
    position: fixed;
    top: 30px;
    right: 50px;
    z-index: 2;
}
</style>

<template>
    <div>
        <div
                v-if="show"
                v-bind:class="{
          'alert':		true,
          'alert-success':(type == 'success'),
          'alert-warning':(type == 'warning'),
          'alert-info':	(type == 'info'),
          'alert-danger':	(type == 'danger'),
          'top': 			(placement === 'top'),
          'top-right': 	(placement === 'top-right')
        }"
                transition="fade"
                v-bind:style="{width:width}"
                role="alert">
            <button v-show="dismissable" type="button" class="close" @click="show = false">
                <span><i class="fa fa-times"></i></span>
            </button>
            <slot></slot>
        </div>
    </div>
</template>

<script>
export default {
    props: {
        type: {
            type: String
        },
        dismissable: {
            type: Boolean,
            default: false,
        },
        show: {
            type: Boolean,
            default: true,
            twoWay: true
        },
        duration: {
            type: Number,
            default: 0
        },
        width: {
            type: String
        },
        placement: {
            type: String
        },
    },
    watch: {
        show(val) {
            //console.log("show = " + val + ", duration " + this.duration);

            // TODO can we register our own transition ont he component to have a function when it is closed?
    /*
            if (this._timeout) {
                clearTimeout(this._timeout);
            }
            if (Boolean(this.duration)) {
                console.log("creating timeout");
                this._timeout = setTimeout(()=> this.show = false, this.duration);
            }
            */
        }
    }
}
</script>

