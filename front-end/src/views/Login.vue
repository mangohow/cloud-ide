<template>

  <div class="login_container">
      <div class="login_box">
          <div class="login_font">Cloud Code</div>
          <!-- 登陆区域 -->
          <el-form ref="loginFormRef" :model="loginForm" :rules="loginFormRules" label-width="0px" class="login_form">
              <!-- 用户名 -->
              <el-form-item prop="username">
                  <el-input v-model="loginForm.username" prefix-icon="el-icon-user-solid" placeholder="请输入用户名"></el-input>
              </el-form-item>
              <!-- 密码 -->
              <el-form-item prop="password">
                  <el-input v-model="loginForm.password" @keyup.enter.native="login" prefix-icon="el-icon-lock" type="password" placeholder="请输入密码"></el-input>
              </el-form-item>
              <!-- 按钮区域 -->
              <el-form-item class="btns">
                  <el-button type="primary" @click="login">登录</el-button>
                  <el-button type="primary" @click="resetLoginForm">重置</el-button>
                  <el-button type="primary" @click="showRegisterDialog">注册</el-button>
              </el-form-item>
          </el-form>
      </div>
      
      <!-- 注册dialog -->
      <el-dialog title="用户注册" :visible.sync="dialogFormVisible">

        <el-form ref="registerFormRef" :model="registerForm" :rules="registerFormRules" label-width="0px" label-position="right" class="register_form">
          <el-form-item prop="nickname" label="昵称" label-width="100px">
            <el-input v-model="registerForm.nickname"></el-input>
          </el-form-item>
          <el-form-item prop="username" label="用户名" label-width="100px">
            <el-input v-model="registerForm.username"></el-input>
          </el-form-item>
          <el-form-item prop="password" label="密码" label-width="100px">
            <el-input v-model="registerForm.password" type="password"></el-input>
          </el-form-item>
          <el-form-item prop="email" label="邮箱" label-width="100px" class="form-email">
            <el-input v-model="registerForm.email"></el-input>
            <el-button @click="getEmailCode" :disabled="getEmailCodeButtonEnable" class="vcb">{{ getEmailCodeButtonName }}</el-button>
          </el-form-item>
          <el-form-item prop="emailCode" label="验证码" label-width="100px">
            <el-input v-model="registerForm.emailCode"></el-input>
          </el-form-item>
        </el-form>
        <div slot="footer" class="dialog-footer">
          <el-button @click="dialogFormVisible = false">取 消</el-button>
          <el-button type="primary" @click="register">提 交</el-button>
        </div>

      </el-dialog>

  </div>

</template>



<script>
import { json } from 'body-parser'
import md5 from 'js-md5'
import {Base64} from "js-base64"

export default {
  data() {
    const validateUsername = async (rule, value, callback) => {
      const reg = /^\w+$/
      if (!reg.test(value)) {
        callback(new Error("用户名只能包含数字,英文字母和下划线"))
        return
      }
      const url = "/auth/username/check" + "?username=" + value
      const {data: res} = await this.$axios.get(url)
      console.log(res)
      if (res.status) {
        callback(new Error(res.message))
        return
      }
      callback()
    };
    
    const validateEmail = (rule, value, callback) => {
      const regEmail = /^([a-zA-Z0-9_-])+@([a-zA-Z0-9_-])+(\.[a-zA-Z0-9_-])+/
      if (!regEmail.test(value)) {
        callback(new Error('请输入合法的邮箱'))
      } else {
        callback()
      }
    };
    
    return {
          // 登录表单的数据绑定对象
          loginForm: {
              username: "",
              password: ""
          },
          //表单的验证规则对象
          loginFormRules: {
              // 验证用户名
              username: [
                  { required: true, message: "请输入用户名", trigger: "blur"},
                  { min: 3, max: 11, message: "长度在 3 到 10 个字符之间", trigger: "blur"}
              ],
              // 验证密码
              password: [
                  { required: true, message: "请输入密码", trigger: "blur"},
                  { min: 4, max: 15, message: "长度在 4 到 15 个字符之间", trigger: "blur"}
              ]
          },
          registerFormRules: {
            nickname: [
              {required: true, message: "请输入昵称", trigger: "blur"},
              {min: 5, max: 32, message: "昵称太短或太长", trigger: "blur"}
            ],
            // 验证用户名
            username: [
                { required: true, message: "请输入用户名", trigger: "blur"},
                { min: 8, max: 16, message: "长度在 8 到 16 个字符之间", trigger: "blur"},
                { validator: validateUsername, trigger: "blur"}
            ],
            // 验证密码
            password: [
                { required: true, message: "请输入密码", trigger: "blur"},
                { min: 8, max: 24, message: "长度在 8 到 24 个字符之间", trigger: "blur"}
            ],
            email: [
                { required: true, message: "请输入邮箱", trigger: "blur"},
                { validator: validateEmail, trigger: "blur"}
            ],
            emailCode: [
                { required: true, message: "请输入验证码", trigger: "blur"},
                { min: 6, max: 6, message: "长度为6", trigger: "blur"}
            ]

          },
          dialogFormVisible: false,
          registerForm: {
            nickname: "",
            username: "",
            password: "",
            email: "",
            emailCode: ""
          },
          getEmailCodeButtonName: "发送验证码",
          getEmailCodeButtonEnable: false,
          countDown: 60,
      }
  },
  methods: {
      resetLoginForm() {
          this.$refs.loginFormRef.resetFields();
      },
      login() {   //表单预校验
          this.$refs.loginFormRef.validate(async (valid) =>{
              if(!valid) return;     //预验证没有通过
              const encodedPasswd = md5(this.loginForm.password)
              const forms = {
                  username: this.loginForm.username,
                  password: encodedPasswd
              }
              const {data: res} = await this.$axios.post("/auth/login", forms);
              if(res.status) {   //登录失败
                  return this.$message.error(res.message);
              }
              if (!res.data) {
                  return this.$message.error(res.message);
              }
              const jsonData = JSON.stringify(res.data)
              const encodedData = Base64.encode(jsonData)
              
              // 登录成功之后:
                  //1 将登陆成功的Token保存到客户端的sessionStorage中，token只应在当前网站
                  //    打开期间生效，所以将Token保存到客户端的sessionStorage中
                  // sessionStorage 是会话期间的存储 localStorage是持久化的存储
                  //2 通过编程式导航跳转到后台主页.路由地址为home
              window.sessionStorage.setItem("userData", encodedData)
              window.sessionStorage.setItem("token", res.data.token)
              window.sessionStorage.setItem("userId", res.data.id)
              await this.$router.push("/dash")
          });
      },
      showRegisterDialog() {
        if (this.dialogFormVisible == false) {
          this.dialogFormVisible = true
        }
      },
      // 用户注册
      register() {
        this.$refs.registerFormRef.validate(async (valid) =>{
              if(!valid) {
                this.$message.warning("输入信息不合法")
                return;     //预验证没有通过
              }
              const encodedPassword = md5(this.registerForm.password)
              const registerForm = {...this.registerForm, password: encodedPassword}
              
              const {data: res} = await this.$axios.post("/auth/register", registerForm)
              if (res.status) {
                this.$message.error(res.message)
                return
              }
              
              this.$message.success(res.message)
              this.dialogFormVisible = false
              this.loginForm.username = this.registerForm.username
              this.registerForm = {nickname: "", username: "", password: "", email: "", emailCode: ""}
          });
      },
      // 获取邮箱验证码
      async getEmailCode() {
        const regEmail = /^([a-zA-Z0-9_-])+@([a-zA-Z0-9_-])+(\.[a-zA-Z0-9_-])+/
        if (!regEmail.test(this.registerForm.email)) {
          this.$message.error('请输入合法的邮箱')
          return
        }

        const url = "/auth/emailCode?email=" + this.registerForm.email
        const {data: res} = await this.$axios.get(url)
        if (res.status) {
          this.$message.error(res.message)
          return
        } else {
          this.$message.success(res.message)
        }

        this.getEmailCodeButtonEnable = true

        var timer = setInterval(() => {
          this.countDown -= 1
          this.getEmailCodeButtonName = this.countDown + "秒后重新发送"
          if (this.countDown == 0) {
            this.getEmailCodeButtonEnable = false
            this.getEmailCodeButtonName = "发送验证码"
            this.countDown = 60
              clearInterval(timer)
            }
        }, 1000)

      }
  }
}
</script>



<style lang="less" scoped>
.login_container {
  background: url("~@/assets/images/back7.jpg");
  width:100%;
  height: 100%;
  background-size:100% 100%;
}
.login_box {
  width: 450px;
  height: 300px;
  background-color: rgba(255, 255, 255, .3);
  border-radius: 3px;
  position: absolute;
  left: 50%;
  top: 50%;
  transform: translate(-50%, -60%);
}

.login_font {
  font-size: 30px;
  font-weight: bold;
  color: #00B5AD;
  width: 100%;
  padding-top: 30px;
}
.login_form {
  position: absolute;
  bottom: 0;
  width: 100%;
  padding: 0 20px;
  box-sizing: border-box;
}
.el-input {
  opacity: 0.5;
}
.btns {
  display: flex;
  justify-content: flex-end;
}

.register_form {
  width: 70%;
  margin-left: 10%;
}

.form-email {
  .el-button {
    position: absolute;
    margin-left: 10px;
  }
}

.vcb {
  width: 120px;
  font-size: 10px;
}
</style>

