/**
 * Amazing Outfits - Beads & Art Single Source of Truth
 * Collections, subcategories, brand copy, option presets, and cards.
 */

export const BRAND = {
  name: "Amazing Outfits",
  tagline: "Beads & Art",
  currency: "KSh",
  city: "Nairobi",
  country: "Kenya",
  origin: "Nairobi, Kenya",
  description: "Handmade beaded jewellery, bags and art, delivered across Kenya.",
  contactEmail: "hello@amazingoutfits.co.ke",
  contactPhone: "+254 700 000 000",
  socialInstagram: "@amazingoutfits.ke"
};

export const COLLECTIONS = [
  {
    slug: "bracelets",
    name: "Bracelets",
    shortName: "Bracelets",
    eyebrow: "Handcrafted Adornments",
    title: "Beaded Bracelets & Cuffs",
    description: "Vibrant Maasai cuffs, delicate seed-bead strands, and cowrie accents handcrafted by Nairobi artisans.",
    subcategories: ["Beaded Bracelets", "Cowrie Shell Bracelets"],
    sizePresets: ["S", "M", "L"],
    sizeLabel: "Size",
    cardImage: "/assets/images/bracelet1.png",
    cardDesc: "Cuffs & beaded strands",
    menuExtra: [
      { name: "Waist Beads", url: "/shop.html?category=Waist%20Beads" }
    ]
  },
  {
    slug: "necklaces",
    name: "Necklaces",
    shortName: "Necklaces",
    eyebrow: "Statement Neckwear",
    title: "Beaded & Cowrie Necklaces",
    description: "Intricate beaded collars, cowrie choker strands, and pendant necklaces celebrating African craftsmanship.",
    subcategories: ["Beaded Necklaces", "Cowrie Shell Necklaces", "Pendant Necklaces"],
    sizePresets: ["16\"", "18\"", "20\"", "24\""],
    sizeLabel: "Length",
    cardImage: "/assets/images/placeholders/necklaces.svg",
    cardDesc: "Collars, chokers & pendants"
  },
  {
    slug: "rings",
    name: "Rings",
    shortName: "Rings",
    eyebrow: "Everyday Details",
    title: "Handmade Beaded Rings",
    description: "Woven seed-bead stack rings and custom initial bands crafted with premium glass seed beads.",
    subcategories: ["Beaded Rings"],
    sizePresets: ["US 5", "US 6", "US 7", "US 8", "US 9", "US 10"],
    sizeLabel: "Ring Size",
    cardImage: "/assets/images/ring1.png",
    cardDesc: "Stack rings & initial bands"
  },
  {
    slug: "bags",
    name: "Bags",
    shortName: "Bags",
    eyebrow: "Wearable Art",
    title: "Beaded Bags & Purses",
    description: "Substantial, fully beaded structured handbags, clutches, and crossbody bags designed to turn heads.",
    subcategories: ["Beaded Bags"],
    sizePresets: ["One size"],
    sizeLabel: "Option",
    singleOption: true,
    cardImage: "/assets/images/bag1.png",
    cardDesc: "Structured bead handbags & clutches"
  },
  {
    slug: "fashion-accessories",
    name: "Fashion Accessories",
    shortName: "Accessories",
    eyebrow: "Heritage Flair",
    title: "Fashion & Wardrobe Accessories",
    description: "Beaded ties, statement belts, hair charms, waist beads, and beaded Ankara accessories for standout style.",
    subcategories: ["Beaded Ties", "Beaded Belts", "Hair Charms", "Waist Beads", "Beaded Ankaras"],
    sizePresets: ["One size"],
    sizeLabel: "Option",
    singleOption: true,
    cardImage: "/assets/images/placeholders/fashion-accessories.svg",
    cardDesc: "Ties, belts, charms & waist beads"
  },
  {
    slug: "art-decor",
    name: "Art & Home Decor",
    shortName: "Art & Decor",
    eyebrow: "Cultural Expressions",
    title: "Wall Art & Decorative Beadwork",
    description: "Framed beaded wall portraits, cultural tapestries, and sculptural decorative bead pieces for your home.",
    subcategories: ["Beaded Wall Portraits", "Beaded Wall Art", "Decorative Beadwork"],
    sizePresets: ["One size", "30×40 cm", "40×50 cm"],
    sizeLabel: "Dimensions",
    cardImage: "/assets/images/placeholders/art-decor.svg",
    cardDesc: "Wall portraits, tapestries & decor"
  }
];

export const VALID_COLLECTION_SLUGS = COLLECTIONS.map(c => c.slug);

export function getCollection(slug) {
  if (!slug) return null;
  return COLLECTIONS.find(c => c.slug.toLowerCase() === slug.toLowerCase()) || null;
}

export function getCollectionBySubcategory(subcategory) {
  if (!subcategory) return null;
  const subLower = subcategory.toLowerCase();
  return COLLECTIONS.find(c => 
    c.subcategories.some(s => s.toLowerCase() === subLower)
  ) || null;
}

export function getSizeOptions(collectionSlug, category) {
  if (category && category.toLowerCase() === "waist beads") {
    return ["28\"", "32\"", "36\"", "40\"", "44\"", "Adjustable"];
  }
  const col = getCollection(collectionSlug);
  if (!col) return ["One size"];
  return col.sizePresets || ["One size"];
}

export function isSingleOptionProduct(collectionSlug, category) {
  if (category && category.toLowerCase() === "waist beads") {
    return false;
  }
  const col = getCollection(collectionSlug);
  if (!col) return true;
  return !!col.singleOption;
}
