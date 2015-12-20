<template>

    <ol class="breadcrumb">
        <div v-for="b in list">
            <div v-if="!b.url">
                <li class="active">
                    {{ b.title }}
                </li>
            </div>
            <div v-else>
                <li>
                    <a v-link="{ path: b.url }">{{ b.title }}</a>
                </li>
            </div>

        </div>

    </ol>
</template>

<script>
// FIXME: never make the last link in list clickable

export default {
    props: ['list']
}

</script>
