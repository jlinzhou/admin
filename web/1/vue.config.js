module.exports = {
    devServer: {
        // host: "localhost", //要设置当前访问的ip 否则失效
        port: 4000,//当前web服务端口
        // open: false, //浏览器自动打开页面
        proxy: {
            '/api': {
                target: 'http://localhost:2001/api/',//目标地址
                ws: true,//是否代理websocket
                changeOrigin: true,//是否跨域
                pathRewrite: {
                    '^/api': ''//url重写
                }
            }
        }
    }
}