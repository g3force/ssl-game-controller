<template>
    <EditableLabelNumber
            label="Goalie: "
            title="The goalie number"
            :value="goalie"
            :callback="updateGoalie"
            :min="0"
            :max="15"/>
</template>

<script>
    import EditableLabelNumber from "../common/EditableLabelNumber";

    export default {
        name: "TeamGoalie",
        components: {EditableLabelNumber},
        props: {
            teamColor: String,
            goalie: Number
        },
        methods: {
            updateGoalie: function (v) {
                this.$socket.sendObj({
                    'modify': {
                        'forTeam': this.teamColor,
                        'goalie': Number(v)
                    }
                })
            },
        }
    }
</script>

<style scoped>

</style>