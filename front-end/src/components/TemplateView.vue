
<template>
    <div>
      <div v-if="dataLoaded" class="tmpl-area" v-for="item in spaceTemplates" :key="item.id">
        <div class="title">{{ item.name }}</div>
        <div class="card-area clear-fix">
          <ul class="clear-fix">
            <li v-for="tmpl in item.tmpls" :key="tmpl.id" :info="tmpl">
              <TemplateCard :info="tmpl" @click.native="tmplSelected(tmpl.id)"></TemplateCard>
            </li>
          </ul>
        </div>
      </div>


      <el-dialog custom-class="space-create-dialog" title="基本信息" :visible.sync="dialogFormVisible" width="40%" :close-on-click-modal="false">
        <el-form :model="spaceForm">
          <el-form-item label="空间名称:" label-width="180px">
            <el-input v-model="spaceForm.name" autocomplete="off" placeholder="请输入空间名称"></el-input>
          </el-form-item>
          <el-form-item label="空间规格:" label-width="180px">
            <el-select v-model="spaceForm.space_spec_id" placeholder="请选择空间规格">
              <el-option v-for="item in spaceSpecs" :key="item.id" :label="item.desc" :value="item.id"></el-option>
            </el-select>
          </el-form-item>
          <el-form-item label="Git仓库:" label-width="180px">
            <el-input v-model="spaceForm.git_repository" autocomplete="off" placeholder="请输要克隆的Git仓库或者忽略"></el-input>
          </el-form-item>
        </el-form>
        <div slot="footer" class="dialog-footer">
          <el-button type="primary" @click="createSpaceAndStart">创建并启动</el-button>
          <el-button type="primary" @click="createSpace">创建</el-button>
          <el-button type="info" @click="dialogFormVisible = false">取 消</el-button>
        </div>
      </el-dialog>
    </div>
</template>



<script>

import TemplateCard from "./TemplateCard.vue"

export default {
    components: {
        TemplateCard
    },
    data() {
        return {
          dataLoaded: false,
          spaceTemplates: [
            {
              id : 1,
              name: "编程语言",
              tmpls: [
                {id: 1, avatar: "", name: "Go", desc: "Go语言环境, 包含go sdk、make、git工具", tags:["Go", "Git"]},
                {id: 2, avatar: "", name: "C++", desc: "C++语言环境, 包含gcc、g++、make、git工具", tags:["Cpp", "Git"]}
              ]
            }
          ],
          spaceSpecs: [
            {id: 1,name: "测试专用",desc: "测试型 2CPU 2GB内存 / 4GB存储"},
            {id: 2,name: "测试专用",desc: "测试型 2CPU 2GB内存 / 4GB存储"}
          ],
          dialogFormVisible: false,
          spaceForm: {
            name: "",
            space_spec_id: "",
            tmpl_id: 0,
            user_id: 0,
            git_repository: "",
          },
        }
    },
    methods: {
        tmplSelected(id) {
          this.dialogFormVisible = true
          this.spaceForm.tmpl_id = id
          this.spaceForm.name = ""
          this.spaceForm.space_spec_id = ""
        },
        
        joinPath(p1, p2) {
            return p1.replace(/\/$/, '') + "/" + p2.replace(/^\//, '');
        },
        
        async getTemplates() {
            const {data: res} = await this.$axios.get("/api/template/list")
            if (res.status) {
                this.$message.error(res.message)
                return
            }
            const kinds = res.data.kinds
            const tmpls = res.data.tmpls.sort((a, b) => {
                return a.id - b.id
            })
            
            kinds.forEach((ele, index) => {
            this.spaceTemplates[index].id = ele.id
            this.spaceTemplates[index].name = ele.name
            this.spaceTemplates[index].tmpls = []
            for (let i = 0; i < tmpls.length; i++) {
                if (ele.id === tmpls[i].kind_id) {
                    var t = tmpls[i]
                    const tags = t.tags.split(',')
                    this.spaceTemplates[index].tmpls.push({...t, tags})
                    for (let j = 0; j < this.spaceTemplates[index].tmpls.length; j++) {
                        const avatar = this.spaceTemplates[index].tmpls[j].avatar
                        if (!avatar.startsWith("http") && !avatar.startsWith("https")) {
                            this.spaceTemplates[index].tmpls[j].avatar = this.joinPath(this.$axios.defaults.baseURL, avatar)
                        }
                    }
                }
            }
            })
      },
      
      validateCreateInfo() {
        if (!(this.spaceForm.name.trim())) {
          this.$message.warning("请输入要创建的工作空间的名称")
          return false
        }
        const value = this.spaceForm.name
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
            return false
        }

        if (!this.spaceForm.space_spec_id) {
          this.$message.warning("请选择要创建的工作空间的规格")
          return false
        }
        
        console.log("验证git仓库：", this.spaceForm.git_repository)
        const regex = /^https:\/\/\S+\.git$/
        this.spaceForm.git_repository = this.spaceForm.git_repository.trim()
        if (this.spaceForm.git_repository.length === 0) {
          return true
        }
        console.log("验证git仓库：", this.spaceForm.git_repository)
        if (!regex.test(this.spaceForm.git_repository.trim())) {
          console.log("git地址无效")
          this.$message.warning("请输入有效的Git仓库地址")
          this.spaceForm.git_repository = ""
          return false
        }
        
        console.log("git正则通过")
        return true
      },
      async getSpaceSpecs() {
        const {data:res} = await this.$axios.get("/api/spec/list")
        if (res.status) {
          this.$message.error(res.message)
          return
        }
        this.spaceSpecs = res.data
      },
      async createSpaceAndStart() {
        this.dialogFormVisible = false

        if (!this.validateCreateInfo()) {
          return          
        }

        const loading = this.$loading({
            lock: true,
            text: 'Loading',
            spinner: 'el-icon-loading',
            background: 'rgba(0, 0, 0, 0.7)'
        });

        const {data:res} = await this.$axios.post("/api/workspace/cas", this.spaceForm)
        if (res.status) {
          this.$message.error(res.message)
          loading.close()
          return
        }

        setTimeout(() => {
          loading.close()
          const spaceUrl =  this.$axios.defaults.workspaceUrl + res.data.sid + "/"
          window.open(spaceUrl, '_blank')
        }, 2000);
        
      },
      async createSpace() {
        this.dialogFormVisible = false

        if (!this.validateCreateInfo()) {
          return          
        }

        const {data:res} = await this.$axios.post("/api/workspace", this.spaceForm)
        if (res.status) {
          this.$message.error(res.message)
        } else {
          this.$message.success(res.message)
        }
      }
      
    },
    
    async mounted() {
      this.spaceForm.user_id = parseInt(window.sessionStorage.getItem("userId"))
      await this.getTemplates()
      this.dataLoaded = true
      await this.getSpaceSpecs()
    }
}
</script>



<style lang="less">

ul, li {
  list-style: none;
}

.header {
    height: 30px;
    color: white;
    font-size: 10px;
    padding: 16px 0 16px 24px;
    text-align: left;
}

.title {
    color: white;
    font-size: 18px;
    line-height: 30px;
    font-weight: 550;
    padding-top: 20px;
}

.space-create-dialog {
  background-color: #323640 !important;
  .el-form-item__label {
    color: #FFF;
  }

  .el-input__inner  {
    background-color: #3C414C;
    border-color: #494D57;
    color: #cfcdcd;
  }
  
}

.el-scrollbar__view, .el-select-dropdown__item {
  
  background-color: #3C414D !important;
  border-color: #494D57 ;
  color: #dfdede !important;
}

.el-select-dropdown {
  border: none !important;
}

.popper__arrow::after {
  border-bottom-color: #3C414D !important;
}

.el-select-dropdown__item:hover {
  background: #6e6180 !important;
}

.el-dialog {
  
  .el-form {
    width: 76%;
  }

  .el-input {
    width: 100%;
  }

  .el-select {
    width: 100%
  }
}

.clear-fix:after {
  content: '';
  display: block;
  height: 0;
  clear: both;
  visibility: hidden;
}

ul {
  padding: 20px 4%;
}

.card-area {
  display: block;
  width: 100%;

  ul {

    li {
      float: left;
      margin: 10px;
      width: 23%;
    }
  }

}

</style>