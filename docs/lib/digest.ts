import { getAllPosts, type BlogPost } from '@/lib/blog'

/**
 * The newest monthly digest, for the pages that should point at it.
 *
 * The changelog is thorough and that is its job: it says why each change was
 * made, which is what you want at two in the morning. It is 757 entries, and
 * nobody reads it to find out what happened this month (grit#266).
 *
 * Found by tag rather than listed in a config. A digest is a blog post like
 * any other, written by hand from a draft scripts/changelog-digest.py
 * produces, and a second place naming it is a second place to forget.
 */
export const DIGEST_TAG = 'release-notes'

export function latestDigest(): BlogPost | undefined {
  // getAllPosts is already sorted newest first.
  return getAllPosts().find((post) => post.tags.includes(DIGEST_TAG))
}

export function allDigests(): BlogPost[] {
  return getAllPosts().filter((post) => post.tags.includes(DIGEST_TAG))
}
