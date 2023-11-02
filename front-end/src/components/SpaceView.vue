
<template>
    <div class="space-view-wapper">
        <div class="title">全部空间</div>
        <SpaceCard v-for="(item, i) in spaces"  v-if="dataLoaded"
            :key="item.id" :space="item" :index="i" @onDeleteSpace="deleteElement"
            @onStopSpace="setSpaceStatus" @onStartSpace="setSpaceStatus"
            @onSpaceNameCheck="checkSpaceName" @onSpaceNameModified="modifySpaceName"
            @click="startSpace(item.id)">
        </SpaceCard>
    </div>
    

</template>



<script>
import SpaceCard from './SpaceCard.vue'

export default {
    components: {
        SpaceCard,
    },
    data() {
        return {
            dataLoaded: false,
            spaces: [
                {
                    id: 1,
                    sid: "",
                    name: "test",
                    create_time: "2023-02-22 02:14:33",
                    environment: "go workspace with go 1.19.3, make"
                },
                {
                    id: 2,
                    sid: "",
                    name: "test",
                    create_time: "2023-02-22 02:14:33",
                    environment: "go workspace with go 1.19.3, make"
                }
            ],
        }
    },
    methods: {
        async getAllSpaces() {
            const {data:res} = await this.$axios.get("/api/workspace/list")
            if (res.status) {
                this.$message.error(res.message)
                return
            }
            this.spaces = res.data
        },

        async startSpace(id) {

        },
        deleteElement(index) {
            this.spaces.splice(index, 1)
        },
        setSpaceStatus(index, status) {
            this.spaces[index].running_status = status
        },
        checkSpaceName(name, index, callback) {
            // 检查是否有工作空间的名称和name相同
            for (let i = 0; i < this.spaces.length; i++) {
                const ele = this.spaces[i]
                if (i == index) {
                    continue
                }
                
                if (ele.name === name) {
                    callback(true)
                    return
                }
            }

            callback(false)
        },
        modifySpaceName(newName, index) {
            this.spaces[index].name = newName
        }
    },
    async mounted() {
        await this.getAllSpaces()
        this.dataLoaded = true
    }
}
</script>



<style scoped>

.title {
    color: #FFF;
    font-size: 18px;
    font-weight: 520;
    padding: 8px 0 36px 0px;
    text-align: left;
    line-height: 40px;
}
.space-view-wapper {
    margin: 25px;
    height: 100%;
}    

</style>