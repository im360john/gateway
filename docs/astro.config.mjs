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
        { label: 'Helm Installation', link: '/helm/gateway' },
        { label: 'Kubernetes Example', link: '/example/k8s' },
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
            { label: 'Terms of Service', slug: 'content/terms' },
            { label: 'Privacy Policy', slug: 'content/privacy' },
            { label: 'Cookie Policy', slug: 'content/cookie' },
          ],
        },
      ],
    }),
    react(),
  ],
});
