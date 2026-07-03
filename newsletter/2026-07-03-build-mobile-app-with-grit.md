---
title: "Build a mobile store with Grit: generate resource → shoppable app"
subtitle: "Point grit generate resource at your Expo app and it scaffolds typed hooks, list + detail screens, and the plumbing for a real store. We wire a Category → Product shopping flow with similar products in one sitting."
series: "The Daily Grit"
edition: 3
date: 2026-07-03
readingTime: "9 min"
author: "Muke JohnBaptist"
tags: [grit, mobile, expo, react-native, code-generation, ecommerce, tutorial]
canonical: "https://gritframework.dev/blog/build-mobile-app-with-grit"
---

A while back we built a full-stack **web** store with Grit — describe your data, and
`grit generate resource` wrote the Go API, the React Query hooks, and an admin panel.

Today we do it on **mobile**. Same command, same models, but now `grit generate
resource` also scaffolds your **Expo** app: a typed hook, a paginated list screen,
and a detail screen for every resource. We'll take that and build a real shopping
flow — **browse categories → tap a category → see its products → open a product →
scroll its similar products** — and you'll write almost none of the plumbing.

Let's build.

## 1. Install or update Grit

Mobile code generation landed in Grit **v3.31.56**, so grab the latest release
before you scaffold.

```bash
# Install (macOS / Linux)
curl -fsSL https://gritframework.dev/install.sh | sh

# Install (Windows PowerShell)
iwr -useb https://gritframework.dev/install.ps1 | iex

# Already have Grit? Update in place:
grit update
```

Prefer Go? `go install github.com/MUKE-coder/grit/v3/cmd/grit@latest`. Confirm you're
on v3.31.56 or newer:

```bash
grit version
```

## 2. Start with a mobile project

If you're starting fresh:

```bash
grit new my-store --mobile   # Go API + an Expo (React Native) app
cd my-store
```

`--mobile` scaffolds an `apps/api` Go backend and an `apps/expo` React Native app
that already has auth (login/register), a themed light/dark UI, and a tab layout
wired to the API.

## 3. Generate the models

A product **belongs to** a category, so we create the parent first.

**Category:**

```bash
grit generate resource Category --fields "name:string,slug:slug,image:file:image"
```

**Product** — linked to Category:

```bash
grit generate resource Product --fields "name:string,slug:slug,price:int,description:text,thumbnail:file:image,category:belongs_to:Category"
```

> **Money tip:** `price:int` stores the amount in the smallest unit (cents) so you
> never fight float rounding. We divide by 100 for display.

### What you just generated for mobile

Alongside the Go model/service/handler, each command now writes into `apps/expo`:

```
✓ apps/expo/hooks/use-<plural>.ts          # typed React Query hook
✓ apps/expo/app/<plural>/index.tsx          # paginated list screen
✓ apps/expo/app/<plural>/[id].tsx           # detail screen
✓ apps/expo/components/ui/screen-header.tsx # shared safe-area header (once)
```

expo-router is file-based, so those files **are** the routes — `/products`,
`/products/:id`, `/categories`, `/categories/:id` all exist now, no registration
step. Open `apps/expo/hooks/use-products.ts` and you'll find exactly what you need:

```ts
export function useProducts(search = "", filters: Record<string, string> = {}) { … } // infinite list
export function useProduct(id: string) { … }                                          // one product
export function useCreateProduct() { … }  // + update / delete mutations
```

The list hook uses `useInfiniteQuery` — pagination for free — and takes **equality
filters**. That second argument is the whole reason the store works, and it leans on
one nice generator detail.

## 4. Get it running

Before we touch a single screen, let's get the whole thing live. From the project
root:

```bash
pnpm install            # install JS deps across the monorepo
docker compose up -d    # Postgres, Redis, MinIO, Mailhog
grit migrate            # create the categories + products tables
grit seed               # seed the admin user: admin@example.com / admin123
grit start              # boot the Go API with hot reload
```

Then, in a second terminal, start the Expo app:

```bash
cd apps/expo
pnpm start              # press a for Android, i for iOS, or scan the QR in Expo Go
```

Log in with **admin@example.com / admin123**. You now have a running mobile app
talking to a running API — time to turn it into a store.

## 5. The one thing that turns a list into a store

Because `Product` **belongs to** `Category`, the generated Go handler makes products
filterable by their foreign key:

```
GET /products?category_id=<id>   →  only that category's products
```

You don't write that — `grit generate resource` wires it for every `belongs_to`
field. On the client, that's just:

```ts
useProducts("", { category_id: id })
```

That's the spine of the entire shopping flow. Everything below is UI on top of the
generated hooks.

## 6. The Shop tab — a grid of categories

Create `apps/expo/app/(tabs)/shop.tsx`:

```tsx
import { View, Text, FlatList, Pressable, ActivityIndicator } from "react-native";
import { Image } from "expo-image";
import { useRouter } from "expo-router";
import { Ionicons } from "@expo/vector-icons";
import { ScreenHeader } from "@/components/ui/screen-header";
import { useTheme } from "@/lib/theme";
import { useCategories, type Category } from "@/hooks/use-categories";

export default function ShopScreen() {
  const router = useRouter();
  const { palette } = useTheme();
  const query = useCategories();
  const categories = query.data?.pages.flatMap((p) => p.data) ?? [];

  const renderItem = ({ item }: { item: Category }) => (
    <Pressable
      onPress={() =>
        router.push({ pathname: "/shop/category/[id]", params: { id: item.id, name: item.name } })
      }
      className="flex-1 m-2 bg-white dark:bg-[#111118] border border-[#E5E7EB] dark:border-[#1f1f2b] rounded-2xl overflow-hidden"
    >
      {item.image?.url ? (
        <Image source={{ uri: item.image.url }} style={{ width: "100%", height: 110 }} contentFit="cover" />
      ) : (
        <View style={{ height: 110 }} className="bg-[#6c5ce7]/10 items-center justify-center">
          <Ionicons name="pricetags-outline" size={28} color="#6c5ce7" />
        </View>
      )}
      <Text className="text-[15px] font-semibold text-[#0F1018] dark:text-white px-3 py-3" numberOfLines={1}>
        {item.name}
      </Text>
    </Pressable>
  );

  return (
    <View className="flex-1 bg-[#F4F4F6] dark:bg-[#0a0a0f]">
      <ScreenHeader title="Shop" subtitle="Browse by category" />
      <FlatList
        data={categories}
        keyExtractor={(item) => item.id}
        renderItem={renderItem}
        numColumns={2}
        contentContainerStyle={{ paddingHorizontal: 16, paddingBottom: 120 }}
        refreshing={query.isRefetching}
        onRefresh={query.refetch}
        ListEmptyComponent={
          query.isLoading ? (
            <ActivityIndicator color={palette.refresh} style={{ marginTop: 40 }} />
          ) : (
            <Text className="text-center text-[#6B7280] dark:text-[#9090a8] mt-16">No categories yet</Text>
          )
        }
      />
    </View>
  );
}
```

Then add it to the tab bar in `apps/expo/app/(tabs)/_layout.tsx`:

```tsx
<Tabs.Screen
  name="shop"
  options={{
    title: "Shop",
    tabBarIcon: ({ color, size }) => <Ionicons name="bag-outline" size={size} color={color} />,
  }}
/>
```

`ScreenHeader` (shipped by the scaffold) gives you a safe-area title bar — and a back
button on any screen that passes `showBack`.

## 7. A category's products

Create `apps/expo/app/shop/category/[id].tsx`. This is where the `category_id` filter
earns its keep:

```tsx
import { View, Text, FlatList, Pressable, ActivityIndicator } from "react-native";
import { Image } from "expo-image";
import { useRouter, useLocalSearchParams } from "expo-router";
import { Ionicons } from "@expo/vector-icons";
import { ScreenHeader } from "@/components/ui/screen-header";
import { useTheme } from "@/lib/theme";
import { useProducts, type Product } from "@/hooks/use-products";

export default function CategoryProductsScreen() {
  const router = useRouter();
  const { palette } = useTheme();
  const { id, name } = useLocalSearchParams<{ id: string; name?: string }>();
  const query = useProducts("", { category_id: id }); // ← the whole trick
  const products = query.data?.pages.flatMap((p) => p.data) ?? [];

  const renderItem = ({ item }: { item: Product }) => (
    <Pressable
      onPress={() => router.push({ pathname: "/shop/product/[id]", params: { id: item.id } })}
      className="flex-1 m-2 bg-white dark:bg-[#111118] border border-[#E5E7EB] dark:border-[#1f1f2b] rounded-2xl overflow-hidden"
    >
      {item.thumbnail?.url ? (
        <Image source={{ uri: item.thumbnail.url }} style={{ width: "100%", height: 130 }} contentFit="cover" />
      ) : (
        <View style={{ height: 130 }} className="bg-[#6c5ce7]/10 items-center justify-center">
          <Ionicons name="cube-outline" size={28} color="#6c5ce7" />
        </View>
      )}
      <View className="p-3">
        <Text className="text-[14px] font-semibold text-[#0F1018] dark:text-white" numberOfLines={1}>{item.name}</Text>
        <Text className="text-[14px] font-bold text-[#6c5ce7] mt-1">${(item.price / 100).toFixed(2)}</Text>
      </View>
    </Pressable>
  );

  return (
    <View className="flex-1 bg-[#F4F4F6] dark:bg-[#0a0a0f]">
      <ScreenHeader title={name || "Category"} subtitle="Products" showBack />
      <FlatList
        data={products}
        keyExtractor={(item) => item.id}
        renderItem={renderItem}
        numColumns={2}
        contentContainerStyle={{ paddingHorizontal: 16, paddingBottom: 40 }}
        onEndReached={() => { if (query.hasNextPage && !query.isFetchingNextPage) query.fetchNextPage(); }}
        onEndReachedThreshold={0.4}
        ListEmptyComponent={
          query.isLoading ? (
            <ActivityIndicator color={palette.refresh} style={{ marginTop: 40 }} />
          ) : (
            <Text className="text-center text-[#6B7280] dark:text-[#9090a8] mt-16">No products in this category</Text>
          )
        }
      />
    </View>
  );
}
```

Infinite scroll, pull-free pagination, and a category-scoped list — all from
`useProducts("", { category_id: id })`.

## 8. Product detail + similar products

Create `apps/expo/app/shop/product/[id].tsx`. "Similar" is just the same filter again
— other products in this product's category, minus itself:

```tsx
import { View, Text, ScrollView, Pressable, ActivityIndicator, FlatList } from "react-native";
import { Image } from "expo-image";
import { useRouter, useLocalSearchParams } from "expo-router";
import { ScreenHeader } from "@/components/ui/screen-header";
import { useTheme } from "@/lib/theme";
import { useProduct, useProducts, type Product } from "@/hooks/use-products";

export default function ProductScreen() {
  const router = useRouter();
  const { palette } = useTheme();
  const { id } = useLocalSearchParams<{ id: string }>();
  const { data: product, isLoading } = useProduct(id);

  // Similar = other products in the same category.
  const similarQuery = useProducts("", product ? { category_id: product.category_id } : {});
  const similar = (similarQuery.data?.pages.flatMap((p) => p.data) ?? [])
    .filter((p) => p.id !== id)
    .slice(0, 10);

  if (isLoading || !product) {
    return (
      <View className="flex-1 bg-[#F4F4F6] dark:bg-[#0a0a0f]">
        <ScreenHeader title="Product" showBack />
        <ActivityIndicator color={palette.refresh} style={{ marginTop: 40 }} />
      </View>
    );
  }

  return (
    <View className="flex-1 bg-[#F4F4F6] dark:bg-[#0a0a0f]">
      <ScreenHeader title="Product" showBack />
      <ScrollView contentContainerStyle={{ paddingBottom: 48 }} showsVerticalScrollIndicator={false}>
        {product.thumbnail?.url ? (
          <Image source={{ uri: product.thumbnail.url }} style={{ width: "100%", height: 300 }} contentFit="cover" />
        ) : null}

        <View className="p-6">
          <Text className="text-[24px] font-bold text-[#0F1018] dark:text-white">{product.name}</Text>
          <Text className="text-[22px] font-bold text-[#6c5ce7] mt-2">${(product.price / 100).toFixed(2)}</Text>
          {product.category?.name ? (
            <View className="self-start bg-[#6c5ce7]/12 px-3 py-1 rounded-full mt-3">
              <Text className="text-[12px] font-medium text-[#6c5ce7]">{product.category.name}</Text>
            </View>
          ) : null}
          <Text className="text-[15px] text-[#6B7280] dark:text-[#9090a8] leading-6 mt-4">{product.description}</Text>

          <Pressable className="bg-[#6c5ce7] rounded-full py-4 items-center mt-6">
            <Text className="text-white font-semibold text-[15px]">Add to Cart</Text>
          </Pressable>
        </View>

        {similar.length > 0 ? (
          <View className="mt-2">
            <Text className="text-[18px] font-bold text-[#0F1018] dark:text-white px-6 mb-3">Similar products</Text>
            <FlatList
              data={similar}
              horizontal
              keyExtractor={(item) => item.id}
              showsHorizontalScrollIndicator={false}
              contentContainerStyle={{ paddingHorizontal: 24 }}
              renderItem={({ item }: { item: Product }) => (
                <Pressable
                  onPress={() => router.push({ pathname: "/shop/product/[id]", params: { id: item.id } })}
                  className="w-40 mr-3 bg-white dark:bg-[#111118] border border-[#E5E7EB] dark:border-[#1f1f2b] rounded-2xl overflow-hidden"
                >
                  {item.thumbnail?.url ? (
                    <Image source={{ uri: item.thumbnail.url }} style={{ width: "100%", height: 110 }} contentFit="cover" />
                  ) : (
                    <View style={{ height: 110 }} className="bg-[#6c5ce7]/10" />
                  )}
                  <View className="p-3">
                    <Text className="text-[13px] font-semibold text-[#0F1018] dark:text-white" numberOfLines={1}>{item.name}</Text>
                    <Text className="text-[13px] font-bold text-[#6c5ce7] mt-1">${(item.price / 100).toFixed(2)}</Text>
                  </View>
                </Pressable>
              )}
            />
          </View>
        ) : null}
      </ScrollView>
    </View>
  );
}
```

Because `Product` preloads its `Category` on the API, `product.category?.name`
renders the category chip with zero extra requests.

## 9. Add a few products and shop

The API and app are already running from step 4. Now add a couple of categories and
products (through the API, or the admin panel if you scaffolded a `--triple` project),
each product pointing at a category. Then open the **Shop** tab and tap through:
category → products → product → similar.

That whole loop — pagination, relationship filtering, image handling, navigation — is
running on generated hooks plus the three screens above.

## The takeaway

`grit generate resource` isn't a backend-only tool anymore. Describe a model once and
you get a Go API, an admin panel, web hooks — **and** typed mobile screens that
already know how to paginate, filter by relationship, and render images. The store you
just built is mostly your generated code; the shop flow is the thin, fun layer on top.

*Go + React. Built with Grit.*

---

<!-- THUMBNAIL PROMPT (Neon style) — remove before publishing to the blog -->
**Gemini thumbnail prompt:** A dark, premium neon-style hero image for a developer
blog. Deep near-black background (#0a0a0f) with a subtle purple technical grid. Center:
a glossy smartphone mock rendering a clean mobile shopping app — a 2-column grid of
product cards with small price tags in electric purple (#6c5ce7). To the left, a
glowing terminal line reading `grit generate resource Product` with a purple glow.
Floating connector lines link the terminal to the phone, suggesting "command →
screens". Top-left corner: the Grit logo (rounded dark app-badge with a bright blue
"G"). Bold, crisp typography overlay: "Build a mobile store with Grit". Neon accents,
soft bloom, high contrast, 16:9.
