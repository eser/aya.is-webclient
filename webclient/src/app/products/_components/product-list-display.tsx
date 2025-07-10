"use client";

import * as React from "react";
import { useTranslations } from "@/shared/modules/i18n/use-translations.tsx";
import { type Profile, ProfileCard } from "@/shared/components/userland/profile-card/profile-card.tsx";

export type ProductListDisplayProps = {
  allProducts: Profile[]; // Assuming products also use the Profile type structure
  searchText: string;
};

export function ProductListDisplay(props: ProductListDisplayProps) {
  const { t } = useTranslations();
  const { allProducts, searchText } = props;

  const filteredProducts = React.useMemo(() => {
    let products = allProducts;

    if (searchText.trim() !== "") {
      const lowerSearchText = searchText.toLowerCase();
      products = products.filter(
        (product) =>
          product.title?.toLowerCase().includes(lowerSearchText) ||
          product.description?.toLowerCase().includes(lowerSearchText),
      );
    }
    return products;
  }, [allProducts, searchText]);

  if (filteredProducts.length === 0) {
    return (
      <div className="py-10 text-center">
        <p className="text-xl text-muted-foreground">{t("Products", "NoProductsFound")}</p>
      </div>
    );
  }

  return (
    <div className="grid grid-cols-1 gap-8 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4">
      {filteredProducts.map((product) => <ProfileCard key={product.slug} profile={product} />)}
    </div>
  );
}
