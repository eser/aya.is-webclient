"use client";

import * as React from "react";
import type { Profile } from "@/shared/components/userland/profile-card/profile-card.tsx";
import { FilterBar } from "./filter-bar.tsx";
import { ProductListDisplay } from "./product-list-display.tsx";

export type ProductsContentProps = {
  initialProfiles: Profile[];
};

export function ProductsContent(props: ProductsContentProps) {
  const [searchText, setSearchText] = React.useState("");

  return (
    <>
      <FilterBar
        searchText={searchText}
        onSearchTextChange={setSearchText}
      />
      <ProductListDisplay
        allProducts={props.initialProfiles}
        searchText={searchText}
      />
    </>
  );
}
