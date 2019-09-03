module.exports = {
    pathPrefix: '/source-chat-relay',
    plugins: [
        {
            resolve: 'gatsby-plugin-manifest',
            options: {
              name: 'Source Chat Relay',
              short_name: 'SCR',
              start_url: '/',
              background_color: '#1d2330',
              theme_color: '#1fb6de',
              display: 'standalone',
            },
        },
        'gatsby-plugin-offline',
        'gatsby-theme-docz',
        'gatsby-plugin-sharp', {
            resolve: 'gatsby-source-filesystem',
            options: {
                path: '${__dirname}/src/pages',
            },
        },
    ]
}