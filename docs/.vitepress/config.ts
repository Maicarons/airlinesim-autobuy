import { defineConfig } from 'vitepress'

export default defineConfig({
  title: 'AirlineSim Autobuy',
  description: 'Automated aircraft market monitoring and purchasing',
  ignoreDeadLinks: true,

  locales: {
    '/': { label: 'English', lang: 'en' },
    '/zh-CN/': { label: '简体中文', lang: 'zh-CN' },
    '/ko/': { label: '한국어', lang: 'ko' },
  },

  themeConfig: {
    logo: '✈️',
    socialLinks: [{ icon: 'github', link: 'https://github.com/Maicarons/airlinesim-autobuy' }],

    locales: {
      '/': {
        nav: [
          { text: 'Guide', link: '/guide/quickstart' },
          { text: 'Developer', link: '/reference/architecture' },
          { text: 'GitHub', link: 'https://github.com/Maicarons/airlinesim-autobuy' },
        ],
        sidebar: {
          '/guide/': [
            {
              text: 'User Guide',
              items: [
                { text: 'Quick Start', link: '/guide/quickstart' },
                { text: 'Configuration', link: '/guide/configuration' },
                { text: 'Purchase Rules', link: '/guide/rules' },
                { text: 'Web UI', link: '/guide/webui' },
                { text: 'Deployment', link: '/guide/deployment' },
                { text: 'FAQ', link: '/guide/faq' },
              ],
            },
          ],
          '/reference/': [
            {
              text: 'Developer Guide',
              items: [
                { text: 'Architecture', link: '/reference/architecture' },
                { text: 'API Reference', link: '/reference/api' },
                { text: 'Configuration Reference', link: '/reference/config' },
                { text: 'Contributing', link: '/reference/contributing' },
              ],
            },
          ],
        },
      },

      '/zh-CN/': {
        nav: [
          { text: '用户指南', link: '/zh-CN/guide/quickstart' },
          { text: '开发者文档', link: '/zh-CN/reference/architecture' },
          { text: 'GitHub', link: 'https://github.com/Maicarons/airlinesim-autobuy' },
        ],
        sidebar: {
          '/zh-CN/guide/': [
            {
              text: '用户指南',
              items: [
                { text: '快速开始', link: '/zh-CN/guide/quickstart' },
                { text: '配置说明', link: '/zh-CN/guide/configuration' },
                { text: '购买规则', link: '/zh-CN/guide/rules' },
                { text: 'Web 界面', link: '/zh-CN/guide/webui' },
                { text: '部署指南', link: '/zh-CN/guide/deployment' },
                { text: '常见问题', link: '/zh-CN/guide/faq' },
              ],
            },
          ],
          '/zh-CN/reference/': [
            {
              text: '开发者文档',
              items: [
                { text: '架构说明', link: '/zh-CN/reference/architecture' },
                { text: 'API 参考', link: '/zh-CN/reference/api' },
                { text: '配置参考', link: '/zh-CN/reference/config' },
                { text: '参与贡献', link: '/zh-CN/reference/contributing' },
              ],
            },
          ],
        },
      },

      '/ko/': {
        nav: [
          { text: '사용자 가이드', link: '/ko/guide/quickstart' },
          { text: '개발자 문서', link: '/ko/reference/architecture' },
          { text: 'GitHub', link: 'https://github.com/Maicarons/airlinesim-autobuy' },
        ],
        sidebar: {
          '/ko/guide/': [
            {
              text: '사용자 가이드',
              items: [
                { text: '빠른 시작', link: '/ko/guide/quickstart' },
                { text: '설정', link: '/ko/guide/configuration' },
                { text: '구매 규칙', link: '/ko/guide/rules' },
                { text: '웹 UI', link: '/ko/guide/webui' },
                { text: '배포', link: '/ko/guide/deployment' },
                { text: 'FAQ', link: '/ko/guide/faq' },
              ],
            },
          ],
          '/ko/reference/': [
            {
              text: '개발자 문서',
              items: [
                { text: '아키텍처', link: '/ko/reference/architecture' },
                { text: 'API 참조', link: '/ko/reference/api' },
                { text: '설정 참조', link: '/ko/reference/config' },
                { text: '기여하기', link: '/ko/reference/contributing' },
              ],
            },
          ],
        },
      },
    },

    footer: {
      message: 'Released under AGPL-3.0.',
      copyright: 'Copyright © 2024',
    },
  },
})