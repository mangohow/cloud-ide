
<template>
    <el-container class="root">
      <el-header class="header">
        <h3 class="site-name">Cloud Code</h3>
        <div class="user-area clearfix">
          <img class="user-avatar" :src="user.avatar" alt="">
          <h4 class="user-nickname">{{user.nickname}}</h4>
        </div>
      </el-header>

      <!-- "#1e1e22" -->
      <el-container class="bottom-main">
        <el-aside>
          <!-- <SpaceTemplateCard :kinds="tmplKinds" @triggered="changeTmpls"></SpaceTemplateCard> -->
          <el-menu
            :default-active="activePath"
            class="el-menu-vertical-demo"
            :router="true"
            background-color="#303336"  
            text-color="#fff"
            active-text-color="#ffd04b">
            <el-menu-item index="/dash/templates">
              <i class="el-icon-menu"></i>
              <span slot="title">空间模板</span>
            </el-menu-item>
            <el-menu-item index="/dash/workspaces">
              <i class="el-icon-document"></i>
              <span slot="title">工作空间</span>
            </el-menu-item>
          </el-menu>
        </el-aside>

        <el-main>
          <router-view></router-view>
        </el-main>
      </el-container>
    </el-container>

</template>

<script>

import {Base64} from "js-base64"

export default {
  name: 'DashBoard',
  data() {
    return {
      user: {},
      activePath: ""
    }
  },
  methods: {
  },
  mounted() {
    this.activePath = this.$route.path

    const data = window.sessionStorage.getItem("userData")
    const jdata = Base64.decode(data)
    this.user = JSON.parse(jdata)
  }
}
</script>

<style lang="less" scoped>

.root {
  height: 100%;
}

.el-header {
  background-color: #373b42;
  color: #333;
  line-height: 60px;
  padding: 0 80px;

  .site-name {
    color: rgb(75, 196, 165);
    text-align: left;
    margin: 3px;
    font-size: 24px;
    font-weight: 550;
    font-family: Georgia;
    float: left;
  }

  .user-area {
    float: right;
    cursor: pointer;
    height: 34px;
    margin: 13px 0;
  }

  .user-avatar {
    float: left;
    width: 30px;
    height: 30px;
    border-radius: 50%;
    border: 2px solid #fff;
  }

  .user-nickname {
    float: left;
    color: #fff;
    margin: 0 0 0 8px;
    font-weight: normal;
    line-height: 2em;
    font-size: 17px;
  }
}


.bottom-main {
  width: 100%;
  height: calc(100% - 60px);
}

.el-aside {
  background-color: #303336;
  height: 100%;
  width: 150px !important;
  float: left;

  .el-menu {
    border: 0 !important;
  }
}

.el-main {
  height: 100%;
  width: calc(100% - 260px);
  float: left;
  padding: 0;
  background-color: rgb(33, 35, 41);
}



</style>