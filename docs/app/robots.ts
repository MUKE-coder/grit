import type { MetadataRoute } from 'next'

export default function robots(): MetadataRoute.Robots {
  return {
    rules: [
      {
        userAgent: '*',
        allow: ['/'],
        disallow: ['/api/'],
      },
    ],
    sitemap: 'https://gritframework.dev/sitemap.xml',
    host: 'https://gritframework.dev',
  }
}
