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
      social: {
        github: 'https://github.com/centralmind/gateway',
      },
      sidebar: [
        { label: 'Introduction', slug: '' },
        {
          label: 'Getting Started',
          items: [
            { label: 'Quickstart', slug: 'docs/content/getting-started/quickstart' },
            { label: 'Installation', slug: 'docs/content/getting-started/installation' },
            { label: 'Generating an API', slug: 'docs/content/getting-started/generating-api' },
            { label: 'Launching an API', slug: 'docs/content/getting-started/launching-api' },
          ],
        },
        {
          label: 'Deployment',
          items: [
            { label: 'Docker Compose', link: '/example/simple' },
            { label: 'Plugin Integrations', link: '/example/complex' },
            { label: 'Kubernetes Example', link: '/example/k8s' },
            { label: 'Helm Installation', link: '/helm/gateway' },
          ],
        },
        {
          label: 'Integration',
          items: [
            { label: 'ChatGPT', slug: 'docs/content/integration/chatgpt' },
            { label: 'LangChain', slug: 'docs/content/integration/langchain' },
            { label: 'Claude Desktop', slug: 'docs/content/integration/claude-desktop' },
            { label: 'Cursor', slug: 'docs/content/integration/cursor' },
          ],
        },
        { label: 'CLI (Command Line Interface)', link: '/cli' },
        {
          label: 'Database Connectors',
          autogenerate: {
            directory: 'connectors',
          },
        },
        {
          label: 'AI Providers',
          items: [
            { label: 'Overview', slug: 'docs/content/ai-providers/overview' },
            { label: 'OpenAI', slug: 'docs/content/ai-providers/openai' },
            { label: 'Anthropic', slug: 'docs/content/ai-providers/anthropic' },
            { label: 'Amazon Bedrock', slug: 'docs/content/ai-providers/bedrock' },
            { label: 'Google VertexAI', slug: 'docs/content/ai-providers/anthropic-vertexai' },
            { label: 'Local Models', slug: 'docs/content/ai-providers/local-models' },
          ],
        },
        {
          label: 'Plugins',
          autogenerate: {
            directory: 'plugins',
          },
        },
        {
          label: 'Terms of Service',
          items: [
            { label: 'Terms of Service', slug: 'docs/content/legal/terms' },
            { label: 'Privacy Policy', slug: 'docs/content/legal/privacy' },
            { label: 'Cookie Policy', slug: 'docs/content/legal/cookie' },
          ],
        },
      ],
    }),
    react(),
  ],
});
