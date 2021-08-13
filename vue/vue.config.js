module.exports = {
  publicPath: "/",
  configureWebpack: {
    devServer: {
      port: 3000,
      watchOptions: {
        poll: true,
      },
    },
  },
};
