<template>

    <form @submit.prevent="postCreate" class="form-horizontal">

        <div class="form-group">
            <div class="col-md-5 col-md-offset-4">
                <input v-model="name" class="form-control" placeholder="{{ $t('island.name') }}">
            </div>
        </div>

        <div class="form-group">
            <div class="col-md-5 col-md-offset-4">
                <input v-model="seed" class="form-control" placeholder="{{ $t('island.seed') }}">
            </div>
        </div>

        <div class="form-group">
            <div class="col-md-5 col-md-offset-4">
                <button v-show="name && seed" class="btn btn-primary form-control">{{ $t('nav.create') }}</button>
            </div>
        </div>

    </form>

</template>

<script>
import Vue from 'vue';

export default {
    data: function() {
        return {
            name: 'Untitled Island',
            seed: "12345",
        }
    },
    methods: {
        postCreate: function() {
            // XXX submit, show ajax. display island when done

            this.$http.post('/island/new', {name: this.name, seed: this.seed}).then(function (response) {
                console.log("island create ok");

            }, function (response) {
                // handle error
                console.log("error: " + data.status);
                console.log(request);
            });

        }
    }
}
</script>
