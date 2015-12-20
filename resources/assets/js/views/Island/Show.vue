<template>

    <div v-if="!Created">
        <form @submit.prevent="postCreate" class="form-horizontal">

            <div class="form-group">
                <div class="col-md-5 col-md-offset-4">
                    <input v-model="Name" class="form-control" placeholder="{{ $t('island.name') }}">
                </div>
            </div>

            <div class="form-group">
                <div class="col-md-5 col-md-offset-4">
                    <input v-model="Seed" class="form-control" placeholder="{{ $t('island.seed') }}">
                </div>
            </div>

            <div class="form-group">
                <div class="col-md-5 col-md-offset-4">
                    <button v-show="Name && Seed" class="btn btn-primary form-control">{{ $t('nav.create') }}</button>
                </div>
            </div>
        </form>
    </div>

    <div v-if="Created">

        <div>
            <img :src="islandMapSrc">
        </div>

    </div>

</template>

<script>
import Vue from 'vue';

export default {
    data: function() {

        return {
            Seed: 666,
            Name: "Untitled Island",
            HeightMap: {},
            Width: 0,
            Height: 0,
            Created: false,
            islandMapSrc: "",
        }
    },
    methods: {
        postCreate: function() {

            this.$http.post('/island/new',{name: this.Name, seed: this.Seed}, {emulateJSON:true}).then(function (response) {
                var data = response.data
                console.log("island create ok");
                console.log(data);

                window.localStorage.setItem('_island', JSON.stringify(data));

                this.$set('Seed', data.Seed);
                this.$set('Name', data.Name);
                this.$set('Width', data.Width);
                this.$set('Height', data.Height);
                this.$set('HeightMap', data.HeightMap);
                this.$set('islandMapSrc', '/img/islands/' + data.Seed + '.png')
                this.$set('Created', true);
            }, function (response) {
                // handle error
                console.log("island create error")
                console.log(response)
            });

        }
    }
}
</script>
