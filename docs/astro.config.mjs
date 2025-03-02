import { defineConfig } from 'astro/config';
import starlight from '@astrojs/starlight';
import react from '@astrojs/react';

// https://astro.build/config
export default defineConfig({
  outDir: './dist',
  site: 'https://docs.centralmind.ai',
  build: {
    assets: 'app_assets',
  },
  integrations: [
    starlight({
      title: 'CentralMind',
      logo: { dark: './src/assets/logo-dark.svg', light: './src/assets/logo-light.svg' },

      customCss: ['./src/styles/custom.css'],
      sidebar: [
        { label: 'Introduction', slug: '' },
        {
          label: 'Getting Started',
          items: [
            { label: 'Installation', slug: 'content/getting-started/installation' },
            { label: 'Generate API', slug: 'content/getting-started/generate-api' },
            { label: 'Launch API', slug: 'content/getting-started/launch-api' },
          ],
        },
        {
          label: 'Deployment',
          items: [
            { label: 'Helm Installation', link: '/helm/gateway' },
            { label: 'Kubernetes Example', link: '/example/k8s' },
          ],
        },
        {
          label: 'Integration',
          items: [
            { label: 'ChatGPT', slug: 'content/integration/chatgpt' },
            { label: 'LangChain', slug: 'content/integration/langchain' },
            { label: 'Claude Desktop', slug: 'content/integration/claude-desktop' },
            { label: 'Local Running Models', slug: 'content/integration/local-running-models' },
          ],
        },
        { label: 'CLI (Command Line Interface)', link: '/cli' },
        {
          label: 'Database Connectors',
          autogenerate: {
            directory: 'connectors'
          }
        },
        {
          label: 'Plugins',
          autogenerate: {
            directory: 'plugins'
          }
        },
        {
          label: 'Terms of Service',
          items: [
            { label: 'Terms of Service', slug: 'content/legal/terms' },
            { label: 'Privacy Policy', slug: 'content/legal/privacy' },
            { label: 'Cookie Policy', slug: 'content/legal/cookie' },
          ],
        },
      ],
    }),
    react(),
  ],
});
