import { Metadata } from "next";
import { notFound } from "next/navigation";
import ProductDetailClient from "./ProductDetailClient";
import { Product } from "@/types/product";
import { getProduct } from "@/services/product.service";

interface ProductDetailPageProps {
  params: Promise<{
    id: string;
  }>;
}

// Generate dynamic metadata for SEO
export async function generateMetadata({
  params,
}: ProductDetailPageProps): Promise<Metadata> {
  const resolvedParams = await params;
  const id = parseInt(resolvedParams.id, 10);

  if (isNaN(id)) {
    return {
      title: "Invalid Product",
    };
  }

  try {
    const product = await getProduct(id);

    const title = `${product.name} - Giá ₫${product.base_price.toLocaleString(
      "vi-VN"
    )} | Shopee`;
    const description =
      product.description?.substring(0, 160) ||
      `Mua ${product.name} giá tốt, chất lượng cao tại Shopee Việt Nam`;

    return {
      title,
      description,
      openGraph: {
        title,
        description,
        images:
          product.images && product.images.length > 0
            ? product.images.map((img) => ({
                url: img,
                alt: product.name,
              }))
            : [],
        type: "website",
      },
      twitter: {
        card: "summary_large_image",
        title,
        description,
        images: product.images || [],
      },
    };
  } catch {
    return {
      title: "Product Not Found",
      description: "The product you are looking for could not be found",
    };
  }
}

// Server Component - SSR for SEO
export default async function ProductDetailPage({
  params,
}: ProductDetailPageProps) {
  const resolvedParams = await params;
  const id = parseInt(resolvedParams.id, 10);

  if (isNaN(id)) {
    notFound();
  }

  // Fetch product on server using service layer
  const product = await getProduct(id);

  // JSON-LD Schema for rich snippets
  const jsonLd = {
    "@context": "https://schema.org",
    "@type": "Product",
    name: product.name,
    image: product.images || [],
    description: product.description || `${product.name} tại Shopee`,
    offers: {
      "@type": "Offer",
      price: product.base_price,
      priceCurrency: "VND",
      availability: product.is_active
        ? "https://schema.org/InStock"
        : "https://schema.org/OutOfStock",
      url: `${
        process.env.NEXT_PUBLIC_SITE_URL || "http://localhost:3000"
      }/products/${product.id}`,
    },
    aggregateRating: {
      "@type": "AggregateRating",
      ratingValue: "4.9",
      reviewCount: "10200",
    },
  };

  return (
    <>
      {/* JSON-LD for SEO */}
      <script
        type="application/ld+json"
        dangerouslySetInnerHTML={{ __html: JSON.stringify(jsonLd) }}
      />

      {/* Client Component with hydrated data */}
      <ProductDetailClient initialProduct={product} />
    </>
  );
}
