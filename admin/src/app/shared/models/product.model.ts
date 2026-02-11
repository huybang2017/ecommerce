export interface Product {
  id: number;
  shop_id: number;
  name: string;
  description: string;
  base_price: number;
  category_id?: number | null;
  category?: { id: number; name: string; slug: string } | null;
  status: string;
  images: string[] | null;
  is_active: boolean;
  sold_count: number;
  created_at: string;
  updated_at: string;
}

export type ProductCreate = Partial<Omit<Product, 'id' | 'created_at' | 'updated_at'>>;
export type ProductUpdate = Partial<Omit<Product, 'id' | 'created_at' | 'updated_at'>>;
