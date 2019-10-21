export default {
    htmlContext: {
        favicon: 'favicon.ico'
    },
    menu: [
        'Home',
        {
            name: 'Getting Started',
            menu: [
                { name: 'Setup', route: 'setup' },
                { name: 'Bot Commands', route: 'bot-commands'},
                { name: 'Recommended Hosts', route: 'recommended-hosts'}
            ]
        },
        {
            name: 'Extended',
            menu: [
                { name: 'Filters', route: 'filters' },
                { name: 'Tips and Tricks', route: 'tips-and-tricks' },
                { name: 'Service', route: 'service' },
                { name: 'Protocol', route: 'protocol' }
            ]
        },
        {
            name: 'Support',
            menu: [
                { name: 'Troubleshooting', route: 'troubleshooting' },
            ]
        }
    ],
    themeConfig: {
        title: 'Source Chat Relay',
        description: 'Communicate between Discord & In-Game, monitor server without being in-game, control the flow of messages and user base engagement!',
        mode: 'dark',
        mdPlugins: [
            {
                resolve: `gatsby-remark-images`,
                options: {
                  maxWidth: 1200,
                },
            },
        ]
    },
}
