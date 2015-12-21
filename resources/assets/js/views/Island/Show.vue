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
            seed {{ Seed }}
        </div>

        <div>
            <img class="img-responsive" :src="islandMapSrc">
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

            this.$resource('/island/new').save({name: this.Name, seed: parseInt(this.Seed, 10)}, function (data, status, request) {
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
