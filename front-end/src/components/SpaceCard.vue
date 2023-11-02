
<template>
    <div class="space-card-wapper">
        <div class="card clearfix">
            <div class="desc">
                <div class="logo-box">
                    <img class="space-logo" :src="space.avatar">
                </div>
                <h3>名称：</h3>
                <h3 class="h3-name">{{space.name}}</h3>
                <h3>环境：</h3>
                <h3 class="h3-environment">{{ space.environment }}</h3>
                <h3>创建时间：</h3>
                <h3 class="h3-time">{{ space.create_time | dateFormat }}</h3>
            </div>
            <div class="operations">
                <span class="space-spec">{{ spaceSpecDesc }}</span>
                <el-tooltip class="item" :class="{'green-border':space.running_status}" effect="dark" content="进入工作空间" placement="top">
                    <i class="iconfont icon-jinru" v-show="space.running_status" @click="enterWorkspace"></i>
                </el-tooltip>
                <el-tooltip class="item" effect="dark" content="停止" placement="top">
                    <i class="iconfont icon-icon_tingzhi" v-show="space.running_status" @click="stopWorkspace"></i>
                </el-tooltip>
                <el-tooltip class="item" effect="dark" content="启动" placement="top">
                    <i class="iconfont icon-qidong" v-show="!space.running_status" @click="startWorkspace"></i>
                </el-tooltip>
                <el-tooltip class="item" effect="dark" content="编辑名称" placement="top">
                    <i class="iconfont icon-bianji" v-show="!space.running_status" @click="openEditDialog"></i>    
                </el-tooltip>
                <el-tooltip class="item" effect="dark" content="删除" placement="top">
                    <i class="iconfont icon-shanchu" @click="deleteWorkspace"></i>    
                </el-tooltip>
            </div>
            
        </div>
    </div>
    

</template>



<script>

import "../assets/icon/iconfont.css"

export default {
    props: ["space", "index"],
    data() {
        return {
        }
    },
    computed: {
        spaceSpecDesc() {
            const spec = this.space.spec
            const strs = spec.mem_spec.split('i')
            const desc = spec.name + " " + spec.cpu_spec + "C" + strs[0]
            return desc
        }
    },
    methods: {
        enterWorkspace() {
            if (this.space.running_status) {
                const url = this.$axios.defaults.workspaceUrl + this.space.sid + "/"
                window.open(url, "_blank")
            }
        },
        async startWorkspace() {
            if (this.space.running_status) {
                return
            }
            // 记载中动画
            const loading = this.$loading({
                lock: true,
                text: 'Loading',
                spinner: 'el-icon-loading',
                background: 'rgba(0, 0, 0, 0.7)'
            });

            const {data:res} = await this.$axios.put("/api/workspace/start", {id: this.space.id})
            if (res.status) {
                this.$message.error(res.message)
                loading.close()
                return
            }

            // 2s钟后在打开
            setTimeout(() => {
                loading.close()
                this.$message.success(res.message)
                const url = this.$axios.defaults.workspaceUrl + res.data.sid + "/"
                window.open(url, "_blank")
                // 通知父组件改变space的running_status字段
                this.$emit("onStartSpace", this.index, true)
            }, 2000);           
        },
        async stopWorkspace() {
            if (!this.space.running_status) {
                return
            }
            
            console.log("sid:", this.space.sid)
            const {data:res} = await this.$axios.put("/api/workspace/stop", {id: this.space.id})
            if (res.status) {
                this.$message.error(res.message)
                return
            }

            this.$message.success(res.message)
            this.$emit("onStopSpace", this.index, false)

        },
        deleteWorkspace(){
            this.$messageBox({
                title: "操作",
                message: "确定删除此空间(删除后无法恢复)?",
                type: "warning",
                showCancelButton: true,
                confirmButtonText: '确定',
                cancelButtonText: '取消',
                customClass: "delete-confirm"
            }).then(async () => {
                
                if (this.space.running_status) {
                    this.$message.warning("工作空间正在运行,请先停止!")
                    return
                }
                
                console.log("id:", this.space.id)
                // 坑：使用delete时，数据必须要在data下即{data: your-data}
                const {data:res} = await this.$axios.delete("/api/workspace", {data: {id: this.space.id}})
                if (res.status) {
                    this.$message.error(res.message)
                    return
                }
                
                this.$message.success(res.message)
                this.$emit("onDeleteSpace", this.index)
            }, () => {
                this.$message({type: 'info', message: '已取消删除'});
            });
        },
        openEditDialog() {
            this.$messageBox.prompt('请输入名称', '修改工作空间名称', {
                    confirmButtonText: '确定',
                    cancelButtonText: '取消',
                    customClass: "delete-confirm"
                }).then(({ value }) => {
                    if (!value) {
                        this.$message.error("名称不能为空")
                        return
                    }

                    //TODO 
                    // 中文字符最多16个 英文32个
                    const chineseMatch = value.match(/[\u4e00-\u9fa5]/g)
                    const englishMatch = value.match(/[a-zA-Z]/g)
                    let chineseCount = 0
                    let englishCount = 0
                    if (chineseMatch) {
                        chineseCount = chineseMatch.length
                    }
                    if (englishMatch) {
                        englishCount = englishMatch.length
                    }

                    if (chineseCount * 2 + englishCount > 32) {
                        this.$message.warning("名称的长度过长,中文字符最多16个,英文字符最多32个")
                        return
                    }

                    this.editWorkspace(value)
                }).catch(() => {

                })
        },
        editWorkspace(newName) {
            if (this.space.name == newName) {
                return
            }

            // 本地先检查是否名称重复 
            this.$emit("onSpaceNameCheck", newName, this.index, async (ret) => {
                if (ret) {
                    this.$message.error("不能和已有的工作空间名称重复")
                    return
                }

                // 发送请求修改名称
                const {data:res} = await this.$axios.put("/api/workspace/name", {name: newName, id: this.space.id})
                if (res.status) {
                    this.$message.error(res.message)
                    return
                }

                this.$message.success(res.message)

                // 通知父组件修改名称
                this.$emit("onSpaceNameModified",newName, this.index)
            })
        }
    }
}
</script>



<style lang="less" scoped>

@import "../assets/style/confirm.css";

.space-card-wapper {
    background-color: #353941;
    height: 60px;
    border: 1px solid #484d55;
    border-radius: 8px;
    padding: 15px 0px;
    margin-bottom: 16px;

    .card{
        text-align: center;
        margin: 0 auto;

        div {
            display: inline-block;
        }

       .desc {
            float: left;
            padding-left: 16px;
            width: calc(80% - 80px);
            text-align: left;
            height: 100%;
            display: flex;
            
            .logo-box {
                padding-right: 16px;
                display: flex;
                align-items: center;
                justify-content: flex-end;
                user-select: none;
            }
            .space-logo {
                width: 40px;
                height: 40px;
            }
        }

       .operations {
            float: right;
            margin-right: 50px;
            line-height: 60px;

            .space-spec {
                color: #fff;
                font-size: 14px;
                margin-right: 24px;
            }
        }

        i {
            color: #dbdada;
            font-size: 25px;
            cursor: pointer;
            padding: 12px;
        }

        i:hover {
            background-color: #050F23;
            border-radius: 50%;
        }

    }

    h3 {
        font-size: 16px;
        font-weight: normal;
        line-height: 36px;
        font-family:"微软雅黑","黑体","宋体";
    }

    h3:nth-child(2n + 1) {
        color: #f5f4f4;
        text-align: left;
    }

    h3:nth-child(2n) {
        color: #a19f9f;
    }

    .h3-name {
        width: 22%;
    }

    .h3-environment {
        width: 38%;
    }

    .h3-time {
        width: 14%;
    }
}

.green-border {
    border: 2px solid #0a7e6f;
    padding: 9px !important;
    border-radius: 50%;
}

.clearfix:after {
    content: "";
    display: table;
    clear: both;
}


.delete-confirm {
    background-color: #363636;
}

</style>
