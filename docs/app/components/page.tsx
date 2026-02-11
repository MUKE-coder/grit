'use client'

import Link from 'next/link'
import { useState, useEffect } from 'react'
import { useSearchParams, useRouter } from 'next/navigation'
import { Code2, Search } from 'lucide-react'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { Badge } from '@/components/ui/badge'
import { Button } from '@/components/ui/button'
import { SiteHeader } from '@/components/site-header'
import { DocsSidebar } from '@/components/docs-sidebar'
import { components, categories, type Platform } from '@/lib/components-data'

export default function ComponentsPage() {
  const searchParams = useSearchParams()
  const router = useRouter()
  const [selectedCategory, setSelectedCategory] = useState('All')
  const [selectedPlatform, setSelectedPlatform] = useState<Platform | 'All'>('All')

  useEffect(() => {
    const categoryParam = searchParams.get('category')
    const platformParam = searchParams.get('platform')

    if (categoryParam && categories.includes(categoryParam as typeof categories[number])) {
      setSelectedCategory(categoryParam)
    }

    if (platformParam) {
      const validPlatforms: Array<Platform | 'All'> = ['All', 'Next.js', 'Expo', 'Go']
      if (validPlatforms.includes(platformParam as Platform | 'All')) {
        setSelectedPlatform(platformParam as Platform | 'All')
      }
    }
  }, [searchParams])

  const updateFilters = (category: string, platform: Platform | 'All') => {
    const params = new URLSearchParams()
    if (category !== 'All') params.set('category', category)
    if (platform !== 'All') params.set('platform', platform)

    const queryString = params.toString()
    const url = queryString ? `/components?${queryString}` : '/components'
    router.push(url, { scroll: false })
  }

  const handleCategoryChange = (category: string) => {
    setSelectedCategory(category)
    updateFilters(category, selectedPlatform)
  }

  const handlePlatformChange = (platform: Platform | 'All') => {
    setSelectedPlatform(platform)
    updateFilters(selectedCategory, platform)
  }

  const filteredComponents = components.filter((component) => {
    const categoryMatch = selectedCategory === 'All' || component.category === selectedCategory
    const platformMatch = selectedPlatform === 'All' || component.platform.includes(selectedPlatform as Platform)
    return categoryMatch && platformMatch
  })

  return (
    <div className="min-h-screen bg-background">
      <SiteHeader />
      <DocsSidebar />

      <main className="lg:pl-64">
        <div className="container max-w-screen-2xl py-10 px-6">
          {/* Header */}
          <div className="mb-10">
            <span className="tag-mono text-primary/80 mb-3 block">Library</span>
            <h1 className="text-4xl font-bold tracking-tight mb-3">
              Components
            </h1>
            <p className="text-muted-foreground max-w-xl">
              Production-ready, full-stack components. Install with a single command.
            </p>
          </div>

          {/* Filters */}
          <div className="mb-10 space-y-5">
            {/* Category filter */}
            <div>
              <h3 className="tag-mono text-muted-foreground/60 mb-3">Category</h3>
              <div className="flex flex-wrap gap-1.5">
                {categories.map((category) => (
                  <button
                    key={category}
                    onClick={() => handleCategoryChange(category)}
                    className={`px-3 py-1.5 rounded-md text-[13px] font-medium transition-all ${
                      selectedCategory === category
                        ? 'bg-primary/15 text-primary border border-primary/20'
                        : 'text-muted-foreground hover:text-foreground hover:bg-accent/50 border border-transparent'
                    }`}
                  >
                    {category}
                  </button>
                ))}
              </div>
            </div>

            {/* Platform filter */}
            <div>
              <h3 className="tag-mono text-muted-foreground/60 mb-3">Platform</h3>
              <div className="flex flex-wrap gap-1.5">
                {(['All', 'Next.js', 'Expo', 'Go'] as const).map((platform) => (
                  <button
                    key={platform}
                    onClick={() => handlePlatformChange(platform)}
                    className={`px-3 py-1.5 rounded-md text-[13px] font-medium transition-all ${
                      selectedPlatform === platform
                        ? 'bg-primary/15 text-primary border border-primary/20'
                        : 'text-muted-foreground hover:text-foreground hover:bg-accent/50 border border-transparent'
                    }`}
                  >
                    {platform}
                  </button>
                ))}
              </div>
            </div>
          </div>

          {/* Divider */}
          <div className="w-full h-px bg-border/40 mb-8" />

          {/* Results count */}
          <div className="mb-6">
            <span className="text-xs text-muted-foreground/50">
              {filteredComponents.length} component{filteredComponents.length !== 1 ? 's' : ''}
            </span>
          </div>

          {/* Components Grid */}
          <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-3">
            {filteredComponents.map((component) => (
              <Link key={component.id} href={`/components/${component.id}`}>
                <Card className="group h-full border-border/40 bg-card/50 hover:bg-card/80 hover:border-primary/20 transition-all duration-300">
                  <CardHeader className="space-y-4 pb-3">
                    {/* Thumbnail */}
                    <div className="aspect-[16/10] rounded-lg bg-accent/50 border border-border/30 flex items-center justify-center relative overflow-hidden">
                      <div className="absolute inset-0 bg-gradient-to-br from-primary/[0.08] to-transparent" />
                      <Code2 className="h-10 w-10 text-primary/40 group-hover:text-primary/60 transition-colors relative z-10" />
                      {component.isPaid && (
                        <div className="absolute top-2.5 right-2.5">
                          <Badge className="bg-amber-500/90 text-black text-[10px] px-2 py-0.5 font-medium hover:bg-amber-500">
                            PRO
                          </Badge>
                        </div>
                      )}
                    </div>

                    {/* Title & Price */}
                    <div className="flex items-start justify-between gap-3">
                      <CardTitle className="text-base font-semibold group-hover:text-primary transition-colors">
                        {component.name}
                      </CardTitle>
                      {component.isPaid ? (
                        <span className="text-sm font-mono font-semibold text-amber-400/90">${component.price}</span>
                      ) : (
                        <span className="text-xs font-mono text-primary/60">FREE</span>
                      )}
                    </div>

                    {/* Description */}
                    <CardDescription className="line-clamp-2 text-sm leading-relaxed text-muted-foreground/70">
                      {component.description}
                    </CardDescription>
                  </CardHeader>

                  <CardContent className="pt-0 space-y-3">
                    {/* Platforms */}
                    <div className="flex flex-wrap gap-1.5">
                      {component.platform.map((platform) => (
                        <span key={platform} className="tag-mono text-[10px] text-muted-foreground/50 px-2 py-1 rounded-md bg-accent/50 border border-border/30">
                          {platform}
                        </span>
                      ))}
                    </div>

                    {/* Meta */}
                    <div className="flex items-center gap-2 text-[11px] text-muted-foreground/40">
                      <span>{component.category}</span>
                      <span className="text-border">|</span>
                      <span>{component.functionalities.length} features</span>
                    </div>
                  </CardContent>
                </Card>
              </Link>
            ))}
          </div>

          {/* No Results */}
          {filteredComponents.length === 0 && (
            <div className="text-center py-20">
              <div className="h-12 w-12 rounded-xl bg-accent/50 border border-border/30 flex items-center justify-center mx-auto mb-4">
                <Search className="h-5 w-5 text-muted-foreground/40" />
              </div>
              <h3 className="text-sm font-medium mb-2">No components found</h3>
              <p className="text-sm text-muted-foreground/60 mb-6">
                Try adjusting your filters
              </p>
              <Button
                variant="outline"
                size="sm"
                className="rounded-lg border-border/60"
                onClick={() => {
                  setSelectedCategory('All')
                  setSelectedPlatform('All')
                  router.push('/components', { scroll: false })
                }}
              >
                Clear Filters
              </Button>
            </div>
          )}
        </div>
      </main>
    </div>
  )
}
