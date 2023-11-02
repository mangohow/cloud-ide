/*
 * @Author: mangohow mghgyf@qq.com
 * @Date: 2022-12-17 15:38:36
 * @LastEditors: mangohow mghgyf@qq.com
 * @LastEditTime: 2022-12-17 15:38:37
 * @FilePath: \front-end\vue.config.js
 * @Description: 这是默认设置,请设置`customMade`, 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
 */

module.exports = {
    devServer: {              //设置本地域名
        host: "localhost",
        port: 8080,
        proxy: {          //设置代理解决跨域问题
            "/": {
                target: "https://192.168.44.100:30443",    //要跨域的域名
                changeOrigin: true                  //是否开启跨域
            }
        }
    }
}